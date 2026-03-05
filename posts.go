package main

import "time"

type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    string    `json:"date"`
	Author  string    `json:"author,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

type PostRepository interface {
	Create(post *Post) error
	GetByID(id int) (*Post, error)
	GetAll() ([]*Post, error)
	Update(post *Post) error
	Delete(id int) error
	Exists(id int) bool
	GetNextID() int
}
