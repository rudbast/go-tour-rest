package auth

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rudbast/go-tour-rest/config"
	"github.com/rudbast/go-tour-rest/core/redis"
	"github.com/rudbast/go-tour-rest/models"
	"github.com/rudbast/go-tour-rest/util"
)

type JWTAuth struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

const EXPIRE_OFFSET = 3600

func (jwtAuth *JWTAuth) Authenticate(user *models.User) (int, bool) {
	// TODO: hash password.
	var userId int

	err := util.DBConn.QueryRow(`
		SELECT id FROM users
		WHERE username = $1
		AND password = $2`, user.Username, user.Password).Scan(&userId)

	switch {
	case err == sql.ErrNoRows:
		return -1, false
	case err != nil:
		log.Fatal(err)
		return -1, false
	default:
		return userId, true
	}
}

func (jwtAuth *JWTAuth) InitBackend() error {
	var (
		err        error
		publicKey  []byte
		privateKey []byte
	)

	// Read public key.
	if publicKey, err = ioutil.ReadFile(config.Data.JWT.PublicKeyPath); err != nil {
		return err
	} else {
		if jwtAuth.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey); err != nil {
			// return err
			panic(err)
		}
	}

	// Read private key.
	if privateKey, err = ioutil.ReadFile(config.Data.JWT.PrivateKeyPath); err != nil {
		return err
	} else {
		if jwtAuth.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey); err != nil {
			// return err
			panic(err)
		}
	}

	return nil
}

func (jwtAuth *JWTAuth) GenerateToken(userId int) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = userId

	if tokenString, err := token.SignedString(jwtAuth.privateKey); err != nil {
		log.Fatal(err)
		return "", err
	} else {
		return tokenString, nil
	}
}

func (jwtAuth *JWTAuth) isBlacklisted(token string) bool {
	redisConn := redis.Connect()

	if redisToken, _ := redisConn.GetValue(token); redisToken == nil {
		return false
	}

	return true
}

func (jwtAuth *JWTAuth) RequireTokenAuthentication(rw http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
	var tokenString string = rq.Header.Get("Authorization")

	// Validate token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtAuth.publicKey, nil
	})

	if err == nil {
		if token.Valid && !jwtAuth.isBlacklisted(tokenString) {
			next(rw, rq)
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(rw, "Invalid token.")
		}
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(rw, "Unauthorized access to resource "+err.Error())
	}
}

func (jwtAuth *JWTAuth) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())

		if remainder > 0 {
			return int(remainder.Seconds() + EXPIRE_OFFSET)
		}
	}

	return EXPIRE_OFFSET
}

func (jwtAuth *JWTAuth) InvalidateToken(tokenString string, tokenHelper interface{}) error {
	req := tokenHelper.(*http.Request)

	tokenRequest, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.publicKey, nil
	})

	if err != nil {
		return err
	}

	redisConn := redis.Connect()

	claims := tokenRequest.Claims.(jwt.MapClaims)

	return redisConn.SetValue(tokenString, tokenString, jwtAuth.getTokenRemainingValidity(claims["exp"]))
}
