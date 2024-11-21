package lobby

import (
	"pong_server/ws"
	"sync"
)

// Lobby para manejar jugadores
type Lobby struct {
	ID      string
	Players []*ws.Client
	Started bool
}

// Mapa de lobbies
var Lobbies = struct {
	sync.Mutex
	m map[string]*Lobby
}{m: make(map[string]*Lobby)}

// Crear un lobby
func handleCreateLobby(client *ws.Client) {
	lobbyID := generateLobbyID()
	lobby := &Lobby{
		ID:      lobbyID,
		Players: []*ws.Client{client},
		Started: false,
	}

	// Registrar lobby
	Lobbies.Lock()
	Lobbies.m[lobbyID] = lobby
	Lobbies.Unlock()

	client.LobbyID = lobbyID

	// Responder al cliente
	client.Conn.WriteJSON(ws.Message{
		Event:   "lobby_created",
		LobbyID: lobbyID,
	})
}
