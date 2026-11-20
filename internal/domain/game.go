package domain

import "errors"

var ErrNotFound = errors.New("запись не найдена")

type UserData struct {
	ID     string `json:"id" db:"id"`
	Board  []int  `json:"board" db:"board"`
	Status string `json:"status" db:"status"`
}

type Storage interface {
	Create(board []int, status string) (UserData, error)
	Get(id string) (UserData, error)
	Update(user UserData) error
}
