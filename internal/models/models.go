package models

import "sync"

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

var (
	// In-memory storage
	UsersStore = struct {
		Sync   sync.RWMutex
		Data   map[string]User
		NextID int
	}{Data: make(map[string]User), NextID: 1}
)

// Response structure
type Response struct {
	Data    interface{} `json:"data"`
	IsOk    bool        `json:"isOk"`
	Message string      `json:"message"`
}
