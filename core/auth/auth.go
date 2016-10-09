package auth

var auth Auth

func GetAuthBackend() Auth {
	if auth == nil {
		auth = &JWTAuth{}

		if err := auth.InitBackend(); err != nil {
			panic(err)
		}
	}

	return auth
}
