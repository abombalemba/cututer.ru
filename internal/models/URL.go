package models

import (
	"time"
)

type URL struct {
	Id uint8
	OriginalUrl string
	ShortUrl string
	SpawnDate time.Time
}
