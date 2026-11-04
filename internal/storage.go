package internal

import "sync"

type Storage struct {
	games sync.Map
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Save(id string, game GameModel) {
	s.games.Store(id, game)
}

func (s *Storage) Get(id string) (GameModel, bool) {
	val, ok := s.games.Load(id)
	if !ok {
		return GameModel{}, false
	}
	return val.(GameModel), true
}
