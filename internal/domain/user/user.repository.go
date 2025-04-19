package user

type Repository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
}
