package models

import "time"

type BorrowStatus string

const (
	BorrowStatusBorrowed BorrowStatus = "borrowed"
	BorrowStatusReturned BorrowStatus = "returned"
	BorrowStatusLost     BorrowStatus = "lost"
)

type Book_Student struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	BookID     uint         `json:"book_id"`
	StudentID  uint         `json:"student_id"`
	BorrowedAt time.Time    `json:"borrowed_time"`
	ReturnedAt time.Time    `json:"return_time"`
	Status     BorrowStatus `json:"status"`
}
