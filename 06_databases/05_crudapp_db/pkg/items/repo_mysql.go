package items

import (
	"database/sql"
)

type ItemMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo(db *sql.DB) *ItemMysqlRepository {
	return &ItemMysqlRepository{DB: db}
}

func (repo *ItemMysqlRepository) GetAll() ([]*Item, error) {
	items := []*Item{}
	rows, err := repo.DB.Query("SELECT id, title, updated FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // надо закрывать соединение, иначе будет течь
	for rows.Next() {
		post := &Item{}
		err = rows.Scan(&post.ID, &post.Title, &post.Updated)
		if err != nil {
			return nil, err
		}
		items = append(items, post)
	}
	return items, nil
}

func (repo *ItemMysqlRepository) GetByID(id int64) (*Item, error) {
	post := &Item{}
	// QueryRow сам закрывает коннект
	err := repo.DB.
		QueryRow("SELECT id, title, updated, description FROM items WHERE id = ?", id).
		Scan(&post.ID, &post.Title, &post.Updated, &post.Description)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *ItemMysqlRepository) Add(elem *Item) (int64, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO items (`title`, `description`) VALUES (?, ?)",
		elem.Title,
		elem.Description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (repo *ItemMysqlRepository) Update(elem *Item) (int64, error) {
	result, err := repo.DB.Exec(
		"UPDATE items SET"+
			"`title` = ?"+
			",`description` = ?"+
			",`updated` = ?"+
			"WHERE id = ?",
		elem.Title,
		elem.Description,
		"dmitry",
		elem.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (repo *ItemMysqlRepository) Delete(id int64) (int64, error) {
	result, err := repo.DB.Exec(
		"DELETE FROM items WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
