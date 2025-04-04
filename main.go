package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/OndrejHana/two-towers/lib"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func main() {
	router := pat.New()
	router.HandleFunc("/game/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got it")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.CreateMock())
	})

	router.HandleFunc("/game/ws", func(w http.ResponseWriter, r *http.Request) {

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

	})

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
	http.ListenAndServe(":8000", router)
}
