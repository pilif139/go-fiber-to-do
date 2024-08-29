package storage

import (
	"go-fiber-server/models"

	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := u.db.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserStorage) Create(userPayload *models.RegisterPayload) (*models.User, error) {
	var user models.User
	err := u.db.Get(&user, "INSERT INTO users (username, email, password) VALUES (?, ?, ?) RETURNING *", userPayload.Username, userPayload.Email, userPayload.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserStorage) GetUsers() ([]models.User, error) {
	var users = make([]models.User, 0)

	err := u.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}
