package services

import (
	"time"
	"math/rand"

	"cututer/internal/config"
)

func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, config.LengthShortUrl)

	for i := range str {
		str[i] = config.Symbols[rand.Intn(len(config.Symbols))]
	}

	return string(str)
}