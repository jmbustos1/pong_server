package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

// Estructura para los mensajes
type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
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
