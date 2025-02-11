package backend

import (
	"net/http"
)

func RegisterRoutes() {
	http.Handle("/", http.FileServer(http.Dir("dist")))
}
