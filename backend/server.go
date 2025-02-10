package backend

import (
	"net/http"
	"two-towers/backend/router"
)

func Serve() error {
	backend.RegisterRoutes()
	return http.ListenAndServe(":8080", nil)
}
