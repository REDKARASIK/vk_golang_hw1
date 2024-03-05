package user

import "errors"

var (
	ErrNoUser  = errors.New("No user found")
	ErrBadPass = errors.New("Invald password")
)

type UserMemoryRepository struct {
	data map[string]*User
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]*User{
			"rvasily": &User{
				ID:       1,
				Login:    "rvasily",
				password: "love",
			},
		},
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
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
