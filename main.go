package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Go websocket")

	setUpAPI()

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}

func setUpAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.loginHandler)
}
