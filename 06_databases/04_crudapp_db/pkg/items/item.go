package items

import (
	"database/sql"
	"strconv"
)

type Item struct {
	ID          uint32
	Title       string
	Description string
	Updated     sql.NullString
}

// позволяет items handlers не импортировать sql
func (it *Item) SetUpdated(val uint32) {
	it.Updated = sql.NullString{String: strconv.Itoa(int(val))}
}

// go install github.com/golang/mock/mockgen@v1.6.0
//go:generate mockgen -source=item.go -destination=repo_mock.go -package=items ItemRepo
type ItemRepo interface {
	GetAll() ([]*Item, error)
	GetByID(int64) (*Item, error)
	Add(*Item) (int64, error)
	Update(*Item) (int64, error)
	Delete(int64) (int64, error)
}
