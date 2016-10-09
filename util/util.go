package util

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rudbast/go-tour-rest/config"
)

var DBConn *sql.DB

// Render json response helper function.
func JsonResponse(rw http.ResponseWriter, rq *http.Request, data interface{}, err error) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	if js, err := json.Marshal(data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	}
}

// Initialize database connection.
func InitDatabase() error {
	var err error

	DBConn, err = sql.Open(config.Data.Database.Driver, config.Data.Database.ConnectionString)

	return err
}
