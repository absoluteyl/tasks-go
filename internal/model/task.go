package model

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}
