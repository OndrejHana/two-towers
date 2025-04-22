package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OndrejHana/two-towers/lib"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	router := pat.New()
	router.Handle("/game/new", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got it")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.CreateMock())
	})))

	router.Handle("/game/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade:", err)
			return
		}
		defer c.Close()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}
			fmt.Printf("recv: %s ", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		}

		fmt.Println("ending")
	}))

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
	http.ListenAndServe(":8000", router)
}
