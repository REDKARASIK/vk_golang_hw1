package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/pkg/domain"
)

type repository struct{}

func (r repository) Create(thread domain.Thread) error {
	reqBody, err := json.Marshal(thread)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://vk-golang.ru:15000/thread", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("failed to create thread remotely")
	}

	return nil
}

func (r repository) Get(id string) (domain.Thread, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://vk-golang.ru:15000/thread?id=%s", id), nil)
	if err != nil {
		return domain.Thread{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.Thread{}, err
	}

	if resp.StatusCode != 200 {
		return domain.Thread{}, errors.New("failed to fetch thread remotely")
	}

	var thread domain.Thread
	err = json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		return domain.Thread{}, err
	}

	return thread, nil
}

func NewRepository() domain.ThreadRepository {
	return repository{}
}
