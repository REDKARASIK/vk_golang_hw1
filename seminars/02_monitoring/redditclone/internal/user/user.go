package user

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	password string
}

type Repo interface {
	Register(username, password string) (*User, error)
	Authorize(username, password string) (*User, error)
}

type MemoryRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		users: make(map[string]*User),
	}
}

func (r *MemoryRepo) Register(username, password string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[username]; ok {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:       uuid.NewString(),
		Username: username,
		password: string(hashedPassword),
	}
	r.users[username] = u
	return u, nil
}

func (r *MemoryRepo) Authorize(username, password string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return u, nil
}
