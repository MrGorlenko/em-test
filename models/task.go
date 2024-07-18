package models

import "time"

type Task struct {
	ID 					uint			`gorm:"primaryKey" json:"id"`
	Title       string    `json:"title" validate:"required"`
  Description string    `json:"description" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
}