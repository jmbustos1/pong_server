package ws

import (
	"log"
	"net/http"
	"pong_server/game"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}
	defer ws.Close()
	// Generar un PlayerID único
	playerID := generatePlayerID()
	client := &Client{Conn: ws, PlayerID: playerID}

	// Registrar el nuevo cliente
	Clients.Lock()
	Clients.m[playerID] = client
	Clients.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error al leer el mensaje:", err)
			Clients.Lock()
			Clients.m[playerID] = client
			Clients.Unlock()
			break
		}

		switch msg.Event {
		case "move_paddle":
			if msg.PlayerID == "1" && msg.Direction == "up" && game.GameState.Paddle1Y > 0 {
				game.GameState.Paddle1Y -= game.PaddleSpeed
			} else if msg.PlayerID == "1" && msg.Direction == "down" && game.GameState.Paddle1Y < game.ScreenHeight-game.PaddleHeight {
				game.GameState.Paddle1Y += game.PaddleSpeed
			} else if msg.PlayerID == "2" && msg.Direction == "up" && game.GameState.Paddle2Y > 0 {
				game.GameState.Paddle2Y -= game.PaddleSpeed
			} else if msg.PlayerID == "2" && msg.Direction == "down" && game.GameState.Paddle2Y < game.ScreenHeight-game.PaddleHeight {
				game.GameState.Paddle2Y += game.PaddleSpeed
			}
		}

		// Broadcast del estado del juego actualizado
		Broadcast <- Message{
			Event:     "sync_game_state",
			GameState: game.GameState,
		}
	}
}

func HandleMessages() {
	for {
		// Recibir un mensaje del canal broadcast
		msg := <-Broadcast
		log.Printf("Transmitiendo mensaje a todos los clientes: %+v\n", msg)

		// Enviar el mensaje a todos los clientes conectados
		Clients.Lock()
		for _, client := range Clients.m { // Cambiado para iterar sobre el nuevo mapa
			err := client.Conn.WriteJSON(msg) // Accede a la conexión WebSocket
			if err != nil {
				log.Println("Error al enviar mensaje:", err)
				client.Conn.Close()
				delete(Clients.m, client.PlayerID) // Elimina por PlayerID
			}
		}
		Clients.Unlock()
	}
}
