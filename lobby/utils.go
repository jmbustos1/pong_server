package lobby

import (
	"crypto/rand"
	"fmt"
)

func generateLobbyID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("lobby-%x", b)
}
