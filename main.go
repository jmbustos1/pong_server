package main

import (
	"log"
	"net/http"
	"pong_server/ws"
)

func main() {
	http.HandleFunc("/ws", ws.HandleConnections)

	// Manejo de mensajes de clientes
	go HandleMessages()

	log.Println("Servidor WebSocket iniciado en ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
