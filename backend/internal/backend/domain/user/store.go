package user

type Store interface {
	CreateUser(user *User) error

	GetUserByLogin(login string) (*User, error)
}
