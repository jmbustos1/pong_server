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
	log.Println("NEW CLIENT ADDED", client.PlayerID)

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

// Método para enviar un mensaje a un cliente
func (c *Client) SendMessage(msg interface{}) {
	err := c.Conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Error al enviar mensaje al cliente %s: %v\n", c.PlayerID, err)
		c.Conn.Close()
		Clients.Lock()
		delete(Clients.m, c.PlayerID) // Elimina el cliente del mapa por PlayerID
		Clients.Unlock()
	}
}

func BroadcastMessage(msg interface{}) {
	Clients.Lock()
	defer Clients.Unlock()

	for _, client := range Clients.m {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error al enviar mensaje al cliente %s: %v\n", client.PlayerID, err)
			client.Conn.Close()
			delete(Clients.m, client.PlayerID)
		}
	}
}
