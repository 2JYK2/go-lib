package random

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
)

func GetRandomString(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	randomString := generateRandomString(min, max)
	return randomString
}

func generateRandomString(min, max int) string {
	b := make([]byte, rand.Intn(max-min+1)+min)
	for i := range b {
		if i%2 == 0 {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else {
			b[i] = numberBytes[rand.Intn(len(numberBytes))]
		}
	}
	return string(b)
}
