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

const CreateTableSQL = `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		short_url TEXT NOT NULL UNIQUE,
		spawn_date TEXT NOT NULL
	);
`

const CheckUrlSQL = `
	SELECT COUNT(*) FROM urls WHERE short_url == ?
`