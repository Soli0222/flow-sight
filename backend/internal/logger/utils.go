package logger

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateRequestID generates a unique request ID
func GenerateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
