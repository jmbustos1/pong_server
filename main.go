package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

// Configuración de la conexión WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Lista de clientes conectados
var clients = make(map[*websocket.Conn]bool)

// Canal para mensajes entre clientes
var broadcast = make(chan Message)

type Vector struct {
	X, Y float64
}

const (
	screenWidth  = 640
	screenHeight = 480
)

// Estructura para los mensajes
type Message struct {
	Event        string  `json:"event"`
	PlayerID     int     `json:"player_id,omitempty"`
	Direction    string  `json:"direction,omitempty"`
	BallPosition Vector  `json:"ball_position,omitempty"`
	Paddle1Y     float64 `json:"paddle1_y,omitempty"`
	Paddle2Y     float64 `json:"paddle2_y,omitempty"`
}

// Estructura para sincronizar el estado del juego
var gameState = struct {
	BallPos  Vector
	Paddle1Y float64
	Paddle2Y float64
}{
	BallPos:  Vector{X: screenWidth / 2, Y: screenHeight / 2},
	Paddle1Y: screenHeight / 2,
	Paddle2Y: screenHeight / 2,
}

func main() {
	// Ruta para manejar las conexiones WebSocket
	http.HandleFunc("/ws", handleConnections)

	// Goroutine para manejar los mensajes entrantes
	go handleMessages()

	// Iniciar el servidor en el puerto 8080
	log.Println("Servidor WebSocket iniciado en ws://localhost:8088/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

// Manejador para establecer nuevas conexiones WebSocket
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Actualizar la conexión a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registrar el nuevo cliente
	clients[ws] = true

	// Leer mensajes del cliente
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error al leer el mensaje:", err)
			delete(clients, ws)
			break
		}
		// Enviar mensaje al canal broadcast
		broadcast <- msg
	}
}

// Manejador de mensajes para transmitir a todos los clientes conectados
func handleMessages() {
	for {
		// Recibir un mensaje del canal broadcast
		msg := <-broadcast

		// Enviar el mensaje a todos los clientes conectados
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error al enviar mensaje:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func sendPlayerInput(conn *websocket.Conn) {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			conn.WriteJSON(Message{Event: "move", Data: "up"})
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			conn.WriteJSON(Message{Event: "move", Data: "down"})
		}
		time.Sleep(16 * time.Millisecond) // 60 envíos por segundo
	}
}
