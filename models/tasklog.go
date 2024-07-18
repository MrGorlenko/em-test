package models

import "time"

type TaskLog struct {
	ID 					uint 				`gorm:"primaryKey" json:"id"`
	TaskID 			uint 				`json:"task_id" validate:"required"`
	UserID 			uint 				`json:"user_id" validate:"required"`
	StartTime 	time.Time 	`json:"start_time"`
  EndTime   	time.Time 	`json:"end_time"`
  CreatedAt 	time.Time 	`json:"created_at"`
  UpdatedAt 	time.Time 	`json:"updated_at"`
}