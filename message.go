package main

import (
	"log"
	"pong_server/game"
	"pong_server/lobby"
	"pong_server/ws"
)

func HandleMessages() {
	for {
		msg := <-ws.Broadcast
		log.Printf("Recibido mensaje: %+v\n", msg)

		switch msg.Event {
		case "start_game":
			lobby.Lobbies.Lock()
			if lobby, exists := lobby.Lobbies.M[msg.LobbyID]; exists {
				if len(lobby.Players) == 2 {
					go game.GameLoop(func(update interface{}) {
						lobby.BroadcastMessageToLobby(update)
					})
					lobby.Started = true
				} else {
					log.Println("No se puede iniciar el juego: no hay suficientes jugadores.")
				}
			}
			lobby.Lobbies.Unlock()
		default:
			ws.BroadcastMessage(msg)
		}
	}
}
