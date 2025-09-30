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

func NewRepository() domain.CommentRepository {
	return repository{}
}

func (r repository) Create(comment domain.Comment) error {
	reqBody, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://vk-golang.ru:16000/comment", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("failed to create comment remotely")
	}

	return nil
}

func (r repository) Like(commentID string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://vk-golang.ru:16000/comment/like?cid=%s", commentID),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("failed to like comment remotely")
	}

	return nil
}
