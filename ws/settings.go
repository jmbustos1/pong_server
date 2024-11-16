package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Configuración de WebSocket
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Mapa seguro para manejar clientes
var Clients = struct {
	sync.Mutex
	m map[*websocket.Conn]bool
}{m: make(map[*websocket.Conn]bool)}

// Canal para mensajes entre clientes
var Broadcast = make(chan Message)

// Estructura para mensajes
type Message struct {
	Event        string  `json:"event"`
	PlayerID     int     `json:"player_id,omitempty"`
	Direction    string  `json:"direction,omitempty"`
	BallPosition Vector  `json:"ball_position,omitempty"`
	Paddle1Y     float64 `json:"paddle1_y,omitempty"`
	Paddle2Y     float64 `json:"paddle2_y,omitempty"`
}

// Vector para representar coordenadas
type Vector struct {
	X, Y float64
}
