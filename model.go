package main

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title    string `json:"title"`
	Author   string `json:"author"`
	Year     int    `json:"year"`
	Category string `json:"category"`
}
