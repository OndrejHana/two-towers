package backend

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {

	mux.Handle("/", http.RedirectHandler("http://localhost:5173/", http.StatusTemporaryRedirect))

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got to /auth")

		if _, err := gothic.CompleteUserAuth(w, r); err == nil {
			fmt.Println("already has user")
			http.Redirect(w, r, "http://localhost:5173/", http.StatusFound)
		} else {
			fmt.Println("dont have user", err.Error())
			gothic.BeginAuthHandler(w, r)
		}
	})

	mux.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got callback")

		gothUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("got user on callback; ", gothUser)
		http.Redirect(w, r, "http://localhost:5173/", http.StatusFound)
	})

	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
