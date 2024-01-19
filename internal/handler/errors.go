package handler

const (
	ErrBadRequest          = "Invalid request body"
	ErrMissingTaskID       = "Missing attribute: id"
	ErrInvalidTaskID       = "Invalid attribute: id"
	ErrNotAllowTaskID      = "Not allowed attribute: id"
	ErrMissingTaskName     = "Missing attribute: name"
	ErrTaskNotFound        = "Task not found"
	ErrInternalServerError = "Internal Server Error"
	ErrUnauthorized        = "Unauthorized"
)
