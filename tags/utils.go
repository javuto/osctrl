package tags

import (
	"fmt"
	"math/rand"
	"time"
)

// Helper to generate a random color in hex for HTML
func RandomColor() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	red := r.Intn(255)
	green := r.Intn(255)
	blue := r.Intn(255)

	return "#" + GetHex(red) + GetHex(green) + GetHex(blue)
}

// Helper to get the hex from an int
func GetHex(num int) string {
	hex := fmt.Sprintf("%x", num)
	if len(hex) == 1 {
		hex = "0" + hex
	}
	return hex
}
