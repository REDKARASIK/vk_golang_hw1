package storage

import (
	"errors"
	"log"
	"sync"
)

var NotFound = errors.New("not found")

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Price uint   `json:"price"`
}

type BookStore struct {
	books  map[int]Book
	mu     sync.RWMutex
	nextID int
}

func NewBookStore() *BookStore {
	return &BookStore{
		mu:    sync.RWMutex{},
		books: map[int]Book{},
	}
}

func (bs *BookStore) AddBook(in Book) (Book, error) {
	log.Println("AddBook called")

	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.nextID++
	in.ID = bs.nextID
	bs.books[in.ID] = in

	return in, nil
}

func (bs *BookStore) GetBooks() ([]Book, error) {
	log.Println("GetBooks called")

	bs.mu.RLock()
	defer bs.mu.RUnlock()

	result := make([]Book, 0, len(bs.books))
	for _, book := range bs.books {
		result = append(result, book)
	}

	return result, nil
}

func (bs *BookStore) GetBook(id int) (Book, error) {
	log.Println("GetBook called")

	bs.mu.RLock()
	defer bs.mu.RUnlock()

	book, ok := bs.books[id]
	if !ok {
		return Book{}, NotFound
	}

	return book, nil
}

func (bs *BookStore) Change(in Book) (Book, error) {
	log.Println("Change called")

	bs.mu.Lock()
	defer bs.mu.Unlock()

	_, ok := bs.books[in.ID]
	if !ok {
		return Book{}, NotFound
	}

	bs.books[in.ID] = in

	return in, nil
}
