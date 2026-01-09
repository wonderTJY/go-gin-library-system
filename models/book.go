package models

import "time"

type BookStatus string

const (
	BookStatusAvailable BookStatus = "available" // 在馆
	BookStatusBorrowed  BookStatus = "borrowed"  // 已借出
	BookStatusLost      BookStatus = "lost"      // 遗失
)

type Book struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stock     uint      `json:"stock"`

	Book_Students []Book_Student `json:"book_students"`
	Copies        []BookCopy     `json:"copies"`
}

type BookCopy struct {
	ID     uint       `gorm:"primaryKey" json:"id"`
	BookID uint       `json:"book_id"`
	Status BookStatus `json:"status"`
}
