package models

import "sync"

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Book struct {
	ID            int    `json:"id"`
	ISBN          string `json:"isbn"`
	Title         string `json:"title"`
	Status        int    `json:"status"` // 0-"new", 1-"reading", 2-"finished"
	Author        string `json:"author"`
	Cover         string `json:"cover"`
	PublishedYear int    `json:"published"`
	Pages         int    `json:"pages"`
	User          string `json:"-"`
}

var (
	// In-memory storage
	UsersStore = struct {
		Sync   sync.RWMutex
		Data   map[string]User
		NextID int
	}{Data: make(map[string]User), NextID: 1}

	BooksStore = struct {
		Sync sync.RWMutex
		Data []Book
		Next int
	}{Data: []Book{}, Next: 1}
)

// Response structure
type Response struct {
	Data    interface{} `json:"data"`
	IsOk    bool        `json:"isOk"`
	Message string      `json:"message"`
}
