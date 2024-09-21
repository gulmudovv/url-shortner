package urlfake

import (
	"fmt"
	"math/rand"
	"time"
)

func URLFake(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var letter = []rune(`abcdefghijklmnopqrstuvwxyz`)
	var org = []string{"ru", "com", "org", "en", "kz", "su", "uz"}
	x := make([]rune, size)
	for i := range x {
		x[i] = letter[rnd.Intn(len(letter))]
	}
	rand.Shuffle(len(org), func(i, j int) { org[i], org[j] = org[j], org[i] })
	y := org[rnd.Intn(len(org))]

	url := fmt.Sprintf("%s.%s", string(x), y)
	return url
}
