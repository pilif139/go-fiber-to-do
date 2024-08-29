package models

type Todo struct {
	ID         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Completed  bool   `json:"completed" db:"completed"`
	Created_at string `json:"created_at" db:"created_at"`
	User_ID    int    `json:"user_id" db:"user_id"`
}

type TodoPayload struct {
	Title     string `json:"title" db:"title"`
	Completed bool   `json:"completed" db:"completed"`
}
