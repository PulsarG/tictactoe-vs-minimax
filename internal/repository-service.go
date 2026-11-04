package internal

import (
	"errors"

	"github.com/google/uuid"
)

type gameRepository struct {
	storage *Storage
}

func NewGameRepository(storage *Storage) GameRepository {
	return &gameRepository{storage: storage}
}

func (r *gameRepository) Save(game Game) error {
	model := ToModel(game)
	r.storage.Save(model.ID, model)
	return nil
}

func (r *gameRepository) Get(id uuid.UUID) (Game, error) {
	model, ok := r.storage.Get(id.String())
	if !ok {
		return Game{}, errors.New("game not foound")
	}
	return ToDomain(model)
}
