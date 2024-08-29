package storage

import (
	"go-fiber-server/models"

	"github.com/jmoiron/sqlx"
)

type TodoStorage struct {
	db *sqlx.DB
}

func NewTodoStorage(db *sqlx.DB) *TodoStorage {
	return &TodoStorage{
		db: db,
	}
}

func (ts *TodoStorage) GetTodos(userID int) ([]models.Todo, error) {
	var todos = make([]models.Todo, 0)
	err := ts.db.Select(&todos, "SELECT * FROM todos WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (ts *TodoStorage) GetTodoByID(userID, todoID int) (*models.Todo, error) {
	var todo models.Todo
	err := ts.db.Get(&todo, "SELECT * FROM todos WHERE user_id=? AND id=?", userID, todoID)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (ts *TodoStorage) CreateTodo(todoPayload *models.TodoPayload, userID int) (*models.Todo, error) {
	var todo models.Todo
	err := ts.db.Get(&todo, "INSERT INTO todos (title, user_id, completed) VALUES (?, ?, ?) RETURNING *", todoPayload.Title, userID, todoPayload.Completed)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (ts *TodoStorage) UpdateTodo(todoPayload *models.TodoPayload, todoID, userID int) error {
	_, err := ts.db.Exec("UPDATE todos SET title=?, completed=? WHERE id=? AND user_id=?", todoPayload.Title, todoPayload.Completed, todoID, userID)
	return err
}

func (ts *TodoStorage) DeleteTodo(todoID, userID int) error {
	_, err := ts.db.Exec("DELETE FROM todos WHERE id=? AND user_id=?", todoID, userID)
	return err
}
