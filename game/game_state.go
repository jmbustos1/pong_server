package game

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

// Vector para posiciones y direcciones
type Vector struct {
	X, Y float64
}

// Estado global del juego
var GameState = struct {
	BallPos  Vector
	Paddle1Y float64
	Paddle2Y float64
}{
	BallPos:  Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2},
	Paddle1Y: ScreenHeight / 2,
	Paddle2Y: ScreenHeight / 2,
}

// UpdateBallPosition actualiza la posición de la bola
func UpdateBallPosition(newPosition Vector) {
	GameState.BallPos = newPosition
}

// UpdatePaddlePosition actualiza la posición de una pala
func UpdatePaddlePosition(playerID int, newY float64) {
	if playerID == 1 {
		GameState.Paddle1Y = newY
	} else if playerID == 2 {
		GameState.Paddle2Y = newY
	}
}

// ResetGame reinicia el estado del juego
func ResetGame() {
	GameState.BallPos = Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2}
	GameState.Paddle1Y = ScreenHeight / 2
	GameState.Paddle2Y = ScreenHeight / 2
}
