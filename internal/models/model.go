package models

type Task struct {
	Id          int
	Title       string
	Description string
	Completed   bool
	Deleted     bool
}
