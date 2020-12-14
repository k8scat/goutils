package random

import (
	"math/rand"
	"time"
)

const (
	DigitsAndAsciiLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVXWYZabcdefghijklmnopqrstuvxwyz"
	LowLetters            = "abcdefghijklmnopqrstuvxwyz"
	UpperLetters          = "ABCDEFGHIJKLMNOPQRSTUVXWYZ"
	Digits                = "0123456789"
	AsciiLetters          = "ABCDEFGHIJKLMNOPQRSTUVXWYZabcdefghijklmnopqrstuvxwyz"
)

func String(s string, n int) string {
	rand.Seed(time.Now().UnixNano())
	buffer := make([]byte, n)
	for i := 0; i < n; i++ {
		buffer[i] = s[rand.Intn(len(s))]
	}
	return string(buffer)
}
