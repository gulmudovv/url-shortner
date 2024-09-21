package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var letter = []rune(`abcdefghijklmnopqrstuvwxyz
	   ABCDEFGHIJKLMNOPQRSTUVWXYZ`)

	b := make([]rune, size)
	for i := range b {
		b[i] = letter[rnd.Intn(len(letter))]
	}
	return string(b)
}
