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

		case "join_lobby":
			if client, exists := ws.Clients.M[msg.PlayerID]; exists {
				lobby.HandleJoinLobby(msg, client)
			} else {
				log.Println("Cliente no encontrado para unirse al lobby:", msg.PlayerID)
			}
		case "create_lobby":
			if client, exists := ws.Clients.M[msg.PlayerID]; exists {
				log.Printf("Recibido mensaje: %+v\n", ws.Clients.M[msg.PlayerID])
				lobby.HandleCreateLobby(msg, client)

			} else {
				log.Println("Cliente no encontrado para crear lobby:", msg.PlayerID)
			}

		case "leave_lobby":
			if client, exists := ws.Clients.M[msg.PlayerID]; exists {
				lobbyID := client.LobbyID
				lobby.Lobbies.Lock()
				if lobbyInstance, exists := lobby.Lobbies.M[lobbyID]; exists {
					// Remover jugador del lobby
					for i, p := range lobbyInstance.Players {
						if p.PlayerID == client.PlayerID {
							lobbyInstance.Players = append(lobbyInstance.Players[:i], lobbyInstance.Players[i+1:]...)
							break
						}
					}
					client.LobbyID = "" // Desasociar cliente del lobby

					// Enviar la lista actualizada de jugadores a todos los jugadores restantes
					lobby.UpdateLobbyPlayers(lobbyInstance)

					// Eliminar el lobby si no quedan jugadores
					if len(lobbyInstance.Players) == 0 {
						// lobby.Lobbies.Unlock()
						delete(lobby.Lobbies.M, lobbyID)
						log.Printf("Lobby %s eliminado porque no quedan jugadores.\n", lobbyID)
						// Notificar a todos los clientes que la lista de lobbies ha cambiado
						lobby.NotifyAllClientsWithLobbies()
					}
				} else {
					log.Println("Lobby no encontrado:", lobbyID)
				}
				lobby.Lobbies.Unlock()
			}

		case "get_lobbies":
			// Manejar el evento `get_lobbies`
			if client, exists := ws.Clients.M[msg.PlayerID]; exists {
				lobby.HandleGetLobbies(client)
			}

		default:
			ws.BroadcastMessage(msg)
		}
	}
}
