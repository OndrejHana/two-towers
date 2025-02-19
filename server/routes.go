package server

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

func getAuth(r *http.Request) (*clerk.SessionClaims, bool) {
	ctx := r.Context()
	return clerk.SessionClaimsFromContext(ctx)
}

func RegisterRoutes(mux *http.ServeMux) {
}
