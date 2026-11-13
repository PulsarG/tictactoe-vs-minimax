package domain

import "errors"

var ErrNotFound = errors.New("запись не найдена")

type UserData struct {
	ID     string `json:"id"`
	Board  []int  `json:"board"`
	Status string `json:"status"`
}

type Storage interface {
	Create(board []int, status string) (UserData, error)
	Get(id string) (UserData, error)
	Update(user UserData) error
}
