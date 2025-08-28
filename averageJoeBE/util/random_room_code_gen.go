package util

import "math/rand/v2"

const roomCodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRoomCode(length int) string {
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = roomCodeChars[rand.IntN(len(roomCodeChars))]
	}
	return string(code)
}
