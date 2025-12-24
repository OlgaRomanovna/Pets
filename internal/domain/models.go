package domain

import "time"

type Pet struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Species   string    `json:"species"`
	CreatedAt time.Time `json:"created_at"`
}
