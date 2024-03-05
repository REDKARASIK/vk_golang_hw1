package user

type User struct {
	ID       uint32
	Login    string
	password string
}

type UserRepo interface {
	Authorize(login, pass string) (*User, error)
}
