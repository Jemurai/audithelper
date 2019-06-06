package data

import (
	"time"
)

// User represents common information we have about a user.
type User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	Admin     bool      `json:"admin"`
	Roles     []string  `json:"roles"`
	LastLogin string    `json:"lastlogin"`
	Created   time.Time `json:"created"`
}
