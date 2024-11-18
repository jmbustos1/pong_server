package main

import (
	"log"
	"net/http"
	"pong_server/game"
	"pong_server/ws"
	"time"
)

func main() {
	http.HandleFunc("/ws", ws.HandleConnections)

	go ws.HandleMessages()

	// Bucle principal del juego
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
		defer ticker.Stop()

		for range ticker.C {
			// Actualiza la lógica del juego
			game.UpdateBall()

			// Convierte el estado actual a una estructura serializable
			currentState := game.GameState

			// Transmite el estado del juego a los clientes
			ws.Broadcast <- ws.Message{
				Event:     "game_update",
				GameState: currentState,
			}
		}
	}()

	log.Println("Servidor WebSocket iniciado en ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
