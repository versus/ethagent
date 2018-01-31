package main

import (
	"crypto/rand"
	"encoding/hex"
)

// Key return random string
func GenToken(key int) string {
	buf := make([]byte, key)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	return hex.EncodeToString(buf)
}
