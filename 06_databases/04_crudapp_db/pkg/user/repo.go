package user

import "errors"

type UserMemoryRepo struct {
	data map[string]*User
}

func NewMemoryRepo() *UserMemoryRepo {
	return &UserMemoryRepo{
		data: map[string]*User{
			"rvasily": &User{
				ID:       1,
				Login:    "rvasily",
				password: "love",
			},
		},
	}
}

var (
	ErrNoUser  = errors.New("User not found")
	ErrBadPass = errors.New("Invalid password")
)

func (repo *UserMemoryRepo) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	// dont do this un production :)
	if u.password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}
