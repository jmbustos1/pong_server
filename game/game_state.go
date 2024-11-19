package game

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	BallSize     = 10
	PaddleWidth  = 10
	PaddleHeight = 80
	PaddleSpeed  = 5
)

// Vector para posiciones y direcciones
type Vector struct {
	X, Y float64
}

// Estructura para representar el estado del juego
type GameStateStruct struct {
	BallPos       Vector  `json:"ball_position"`
	BallDirection Vector  `json:"ball_direction"`
	Paddle1Y      float64 `json:"paddle1_y"`
	Paddle2Y      float64 `json:"paddle2_y"`
}
type GameSession struct {
	ID        string
	GameState GameStateStruct
}

// Estado global del juego
var GameState = GameStateStruct{
	BallPos:       Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2},
	BallDirection: Vector{X: 1, Y: 1},
	Paddle1Y:      ScreenHeight / 2,
	Paddle2Y:      ScreenHeight / 2,
}

// UpdateBall actualiza la posición de la bola y maneja colisiones
func UpdateBall() {
	GameState.BallPos.X += GameState.BallDirection.X * 4
	GameState.BallPos.Y += GameState.BallDirection.Y * 4

	// Rebote en bordes superior e inferior
	if GameState.BallPos.Y <= 0 || GameState.BallPos.Y >= ScreenHeight-BallSize {
		GameState.BallDirection.Y *= -1
	}

	// Rebote en las palas
	checkPaddleCollision(GameState.Paddle1Y, true)  // Pala izquierda
	checkPaddleCollision(GameState.Paddle2Y, false) // Pala derecha
}

// checkPaddleCollision maneja las colisiones entre la bola y las palas
func checkPaddleCollision(paddleY float64, isLeftPaddle bool) {
	paddleX := 20.0
	if !isLeftPaddle {
		paddleX = ScreenWidth - 30
	}

	if (isLeftPaddle && GameState.BallPos.X <= paddleX+PaddleWidth) ||
		(!isLeftPaddle && GameState.BallPos.X+BallSize >= paddleX) {
		if GameState.BallPos.Y+BallSize >= paddleY && GameState.BallPos.Y <= paddleY+PaddleHeight {
			collisionPoint := (GameState.BallPos.Y - paddleY) / PaddleHeight

			// Rebote basado en la posición de colisión
			if collisionPoint <= 0.1 && GameState.BallDirection.Y > 0 {
				GameState.BallDirection.X *= -1
				GameState.BallDirection.Y *= -1
			} else if collisionPoint >= 0.9 && GameState.BallDirection.Y < 0 {
				GameState.BallDirection.X *= -1
				GameState.BallDirection.Y *= -1
			} else {
				GameState.BallDirection.X *= -1
			}

			// Ajustar la posición de la bola para evitar múltiples colisiones
			if isLeftPaddle {
				GameState.BallPos.X = paddleX + PaddleWidth
			} else {
				GameState.BallPos.X = paddleX - BallSize
			}
		}
	}
}

// UpdatePaddlePosition actualiza la posición de una pala según el jugador
func UpdatePaddlePosition(playerID string, newY float64) {
	if playerID == "1" {
		GameState.Paddle1Y = newY
	} else if playerID == "2" {
		GameState.Paddle2Y = newY
	}
}

// ResetGame reinicia el estado del juego
func ResetGame() {
	GameState.BallPos = Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2}
	GameState.BallDirection = Vector{X: 1, Y: 1}
	GameState.Paddle1Y = ScreenHeight / 2
	GameState.Paddle2Y = ScreenHeight / 2
}
