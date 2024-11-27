package lobby

import (
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
	m map[string]*Lobby
}{m: make(map[string]*Lobby)}

func handleCreateLobby(msg ws.Message, client *ws.Client) {
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
	Lobbies.m[lobbyID] = newLobby
	Lobbies.Unlock()

	// Asociar al cliente con el lobby
	client.LobbyID = lobbyID

	// Notificar al cliente que el lobby fue creado
	client.SendMessage(ws.Message{
		Event:     "lobby_created",
		LobbyID:   lobbyID,
		LobbyName: msg.LobbyName,
	})
}

func handleJoinLobby(msg ws.Message, client *ws.Client) {
	lobbyID := msg.LobbyID

	// Validar que el lobby exista
	Lobbies.Lock()
	lobby, exists := Lobbies.m[lobbyID]
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
	lobby, exists := Lobbies.m[lobbyID]
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
