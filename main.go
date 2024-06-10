package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Go websocket")

	setUpAPI()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setUpAPI() {

	manager := NewManager()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
}
