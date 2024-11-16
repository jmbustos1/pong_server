package ws

import (
	"log"
	"net/http"
	"pong_server/game"
)

// HandleConnections gestiona las conexiones WebSocket de los clientes
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registrar el nuevo cliente
	Clients.Lock()
	Clients.m[ws] = true
	Clients.Unlock()

	// Leer mensajes del cliente
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error al leer el mensaje:", err)
			Clients.Lock()
			delete(Clients.m, ws)
			Clients.Unlock()
			break
		}

		// Procesar mensajes según el evento
		switch msg.Event {
		case "move_paddle":
			game.UpdatePaddlePosition(msg.PlayerID, msg.Paddle1Y)
			// case "ball_update":
			// 	game.UpdateBallPosition(msg.BallPosition)
		}

		// Enviar mensaje al canal broadcast
		Broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Recibir un mensaje del canal broadcast
		msg := <-Broadcast

		// Enviar el mensaje a todos los clientes conectados
		Clients.Lock()
		for client := range Clients.m {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error al enviar mensaje:", err)
				client.Close()
				delete(Clients.m, client)
			}
		}
		Clients.Unlock()
	}
}
