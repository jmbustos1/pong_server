package game

import (
	"time"
)

const (
	FrameDuration = time.Second / 60 // 60 FPS
)

// GameLoop actualiza el estado del juego y difunde cambios.
func GameLoop(broadcastFunc func(msg interface{})) {
	ticker := time.NewTicker(FrameDuration)
	defer ticker.Stop()

	for range ticker.C {
		// Actualizar posición de la bola
		UpdateBall()

		// Crear mensaje para difundir
		state := GameState
		msg := map[string]interface{}{
			"event":     "game_update",
			"ball_pos":  map[string]float64{"x": state.BallPos.X, "y": state.BallPos.Y},
			"paddle1_y": state.Paddle1Y,
			"paddle2_y": state.Paddle2Y,
		}

		// Difundir a los clientes conectados
		broadcastFunc(msg)
	}
}
