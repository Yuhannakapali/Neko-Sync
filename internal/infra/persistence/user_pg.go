package persistence

import (
	"database/sql"
	"nekosync/internal/domain/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Create(u *user.User) error {
	_, err := repo.db.Exec("INSERT INTO users (email, name) VALUES ($1, $2)", u.Email, u.Name)
	return err
}

func (r *UserRepo) GetByID(id int) (*user.User, error) {
	var u user.User
	err := r.db.QueryRow("SELECT id, email, name FROM users WHERE id = $1", id).
		Scan(&u.Id, &u.Email, &u.Name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
