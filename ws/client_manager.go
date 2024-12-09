package ws

import (
	"log"
	"net/http"
)

// // CAMBIAR A ESTRUCTURA DE MENSJAE
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
	Clients.M[playerID] = client
	Clients.Unlock()

	log.Println("NEW CLIENT ADDED", client.PlayerID)
	client.SendMessage(Message{
		Event: "test_message",
		Data:  "Conexión establecida correctamente HOLAAAAAA",
	})

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error al leer el mensaje:", err)
			Clients.Lock()
			delete(Clients.M, playerID)
			Clients.Unlock()
			break
		}

		// Añadir PlayerID al mensaje
		msg.PlayerID = playerID

		// Enviar mensaje al canal Broadcast
		Broadcast <- msg
	}
}

// Método para enviar un mensaje a un cliente
func (c *Client) SendMessage(msg interface{}) {
	err := c.Conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Error al enviar mensaje al cliente %s: %v\n", c.PlayerID, err)
		c.Conn.Close()
		Clients.Lock()
		delete(Clients.M, c.PlayerID) // Elimina el cliente del mapa por PlayerID
		Clients.Unlock()
	}
}

func BroadcastMessage(msg interface{}) {
	Clients.Lock()
	defer Clients.Unlock()

	for _, client := range Clients.M {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error al enviar mensaje al cliente %s: %v\n", client.PlayerID, err)
			client.Conn.Close()
			delete(Clients.M, client.PlayerID)
		}
	}
}
