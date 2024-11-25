package ws

import (
	"crypto/rand"
	"fmt"
)

func generatePlayerID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
