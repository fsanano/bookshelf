package models

import "sync"

// User model
type User struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// In-memory storage
var (
	UsersStore = struct {
		Sync sync.RWMutex
		Data map[string]User
	}{Data: make(map[string]User)}
)

// Response structure
type Response struct {
	Data    interface{} `json:"data"`
	IsOk    bool        `json:"isOk"`
	Message string      `json:"message"`
}
