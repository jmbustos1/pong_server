package lobby

import (
	"log"
	"pong_server/ws"
	"sync"
)

// Lobby para manejar jugadores
type Lobby struct {
	ID      string
	Name    string
	Players []*ws.Client
	Host    *ws.Client
	Started bool
}

// Mapa de lobbies
var Lobbies = struct {
	sync.Mutex
	M map[string]*Lobby
}{M: make(map[string]*Lobby)}

func HandleCreateLobby(msg ws.Message, client *ws.Client) {
	if msg.LobbyName == "" {
		client.SendMessage(ws.Message{
			Event: "error",
			Data:  "Lobby name cannot be empty",
		})
		return
	}

	lobbyID := generateLobbyID() // Función para generar IDs únicos para lobbies
	newLobby := &Lobby{
		ID:      lobbyID,
		Name:    msg.LobbyName,
		Host:    client,
		Players: []*ws.Client{client},
		Started: false,
	}

	// Registrar el lobby
	Lobbies.Lock()
	Lobbies.M[lobbyID] = newLobby
	Lobbies.Unlock()

	// Asociar al cliente con el lobby
	client.LobbyID = lobbyID

	// Notificar al cliente que el lobby fue creado
	client.SendMessage(ws.Message{
		Event:     "lobby_created",
		LobbyID:   lobbyID,
		LobbyName: msg.LobbyName,
	})
	log.Printf("Lobby creado: %s con ID: %s", msg.LobbyName, lobbyID)
	log.Printf("Lobby creado: %s", Lobbies)
}

func handleJoinLobby(msg ws.Message, client *ws.Client) {
	lobbyID := msg.LobbyID

	// Validar que el lobby exista
	Lobbies.Lock()
	lobby, exists := Lobbies.M[lobbyID]
	Lobbies.Unlock()
	if !exists {
		client.SendMessage(ws.Message{
			Event: "error",
			Data:  "Lobby not found",
		})
		return
	}

	// Agregar al cliente al lobby
	Lobbies.Lock()
	lobby.Players = append(lobby.Players, client)
	Lobbies.Unlock()

	// Asociar al cliente con el lobby
	client.LobbyID = lobbyID

	// Notificar al cliente que se unió al lobby
	client.SendMessage(ws.Message{
		Event:     "joined_lobby",
		LobbyID:   lobby.ID,
		LobbyName: lobby.Name,
	})

	// Notificar a los jugadores en el lobby
	for _, player := range lobby.Players {
		player.SendMessage(ws.Message{
			Event:    "player_joined",
			PlayerID: client.PlayerID,
		})
	}
}

func handleStartGame(msg ws.Message, client *ws.Client) {
	lobbyID := client.LobbyID

	// Validar que el lobby exista
	Lobbies.Lock()
	lobby, exists := Lobbies.M[lobbyID]
	Lobbies.Unlock()
	if !exists {
		client.SendMessage(ws.Message{
			Event: "error",
			Data:  "Lobby not found",
		})
		return
	}

	// Verificar que el cliente sea el host del lobby
	if lobby.Host != client {
		client.SendMessage(ws.Message{
			Event: "error",
			Data:  "Only the host can start the game",
		})
		return
	}

	// Marcar el juego como iniciado
	Lobbies.Lock()
	lobby.Started = true
	Lobbies.Unlock()

	// Notificar a los jugadores que el juego ha comenzado
	for _, player := range lobby.Players {
		player.SendMessage(ws.Message{
			Event: "game_started",
		})
	}
}

func (l *Lobby) BroadcastToPlayers(msg ws.Message) {
	for _, player := range l.Players {
		player.SendMessage(msg)
	}
}

func (l *Lobby) BroadcastMessageToLobby(msg interface{}) {
	for _, client := range l.Players {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error al enviar mensaje al cliente %s: %v\n", client.PlayerID, err)
			client.Conn.Close()
			removeClientFromLobby(client.PlayerID, l.ID)
		}
	}
}

func removeClientFromLobby(playerID string, lobbyID string) {
	Lobbies.Lock()
	defer Lobbies.Unlock()

	if lobby, exists := Lobbies.M[lobbyID]; exists {
		// Filtra los jugadores para eliminar al cliente
		var updatedPlayers []*ws.Client
		for _, player := range lobby.Players {
			if player.PlayerID != playerID {
				updatedPlayers = append(updatedPlayers, player)
			}
		}
		lobby.Players = updatedPlayers

		// Si no quedan jugadores, puedes eliminar el lobby
		if len(lobby.Players) == 0 {
			delete(Lobbies.M, lobbyID)
			log.Printf("Lobby %s eliminado porque no quedan jugadores.\n", lobbyID)
		}
	}
}

func HandleGetLobbies(client *ws.Client) {
	lobbies := []string{}
	Lobbies.Lock()
	for _, lobby := range Lobbies.M {
		lobbies = append(lobbies, lobby.Name)
	}
	Lobbies.Unlock()

	client.SendMessage(ws.Message{
		Event:   "lobbies_list",
		Lobbies: lobbies,
	})
}

// Función para actualizar los jugadores en un lobby
func UpdateLobbyPlayers(lobbyID string) {
	Lobbies.Lock()
	defer Lobbies.Unlock()

	if lobby, exists := Lobbies.M[lobbyID]; exists {
		playerNames := []string{}
		for _, player := range lobby.Players {
			playerNames = append(playerNames, player.PlayerID) // Usa el ID del jugador
		}

		// Enviar la lista actualizada a los jugadores
		for _, client := range lobby.Players {
			client.SendMessage(ws.Message{
				Event:   "lobby_players",
				Lobbies: playerNames,
			})
		}
	}
}
