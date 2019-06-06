package data

import (
	"time"
)

// Resource represents common information we have about a thing (file, etc.).
type Resource struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Public      bool      `json:"public"`
	Dataclass   string    `json:"dataclass"`
	Source      string    `json:"source"`
	Created     time.Time `json:"created"`
	Touched     time.Time `json:"touched"`
}
