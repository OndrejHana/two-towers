package server

import "net/http"

func Run() error {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	return http.ListenAndServe("", mux)
}
