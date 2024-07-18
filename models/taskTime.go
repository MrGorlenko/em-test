package models

type TaskTime struct {
	TaskID  uint   `json:"task_id"`
	Title   string `json:"title"`
	Hours   int    `json:"hours"`
	Minutes int    `json:"minutes"`
}
