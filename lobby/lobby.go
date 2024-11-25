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
	Host    []*ws.Client
	Started bool
}

// Mapa de lobbies
var Lobbies = struct {
	sync.Mutex
	m map[string]*Lobby
}{m: make(map[string]*Lobby)}

func handleCreateLobby(msg ws.Message, client *ws.Client) {
	lobbyID := generateLobbyID()
	newLobby := &Lobby{
		ID:      lobbyID,
		Name:    msg.LobbyName,
		Host:    client,
		Players: []*ws.Client{client},
	}
	Lobbies.Lock()
	Lobbies.m[lobbyID] = newLobby
	Lobbies.Unlock()

	client.SendMessage(Message{
		Event:     "lobby_created",
		LobbyID:   lobbyID,
		LobbyName: msg.LobbyName,
	})
}
