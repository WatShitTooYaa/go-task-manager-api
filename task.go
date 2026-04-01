package main

type Task struct {
	Id          uint16 `json:"id"`
	Content     string `json:"content" validate:"required,min=3,max=200"`
	IsCompleted bool   `json:"isCompleted"`
	Time        string `json:"timestamp"`
	Priority    string `json:"priority" validate:"required,oneof=low medium high"`
}

type CreateTaskRequest struct {
	Content  string `json:"content" validate:"required,min=3,max=200"`
	Priority string `json:"priority" validate:"required,oneof=low medium high"`
}

type UpdateTaskRequest struct {
	Content   string `json:"content" validate:"required,min=3,max=200"`
	Priority  string `json:"priority" validate:"required,oneof=low medium high"`
	Completed bool   `json:"completed"`
}
