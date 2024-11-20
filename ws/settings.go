package ws

import (
	"net/http"
	"pong_server/game"
	"sync"

	"github.com/gorilla/websocket"
)

// Configuración de WebSocket
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Cliente conectado
type Client struct {
	Conn     *websocket.Conn
	PlayerID string
	LobbyID  string
}

// Mapa seguro para manejar clientes
var Clients = struct {
	sync.Mutex
	m map[string]*Client
}{m: make(map[string]*Client)}

// Canal para mensajes entre clientes
var Broadcast = make(chan Message)

// Lobby para manejar jugadores
type Lobby struct {
	ID      string
	Players []*Client
	Started bool
}

// Mapa de lobbies
var Lobbies = struct {
	sync.Mutex
	m map[string]*Lobby
}{m: make(map[string]*Lobby)}

// Estructura para mensajes
type Message struct {
	Event        string               `json:"event"`
	PlayerID     string               `json:"player_id,omitempty"`
	Direction    string               `json:"direction,omitempty"`
	GameState    game.GameStateStruct `json:"game_state,omitempty"`
	BallPosition *game.Vector         `json:"ball_position,omitempty"`
	Paddle1Y     float64              `json:"paddle1_y,omitempty"`
	Paddle2Y     float64              `json:"paddle2_y,omitempty"`
}

// Vector para representar coordenadas
type Vector struct {
	X, Y float64
}
