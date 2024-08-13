package model

import (
	"time"
)

type Material struct {
	UUID      string    `json:"uuid"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MaterialFilter struct {
	Type     string
	DateFrom time.Time
	DateTo   time.Time
	Status   string
}
