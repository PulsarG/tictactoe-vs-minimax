package internal

type Storage struct {
	games sync.Map
}

func NewStorage() *Storage {
	return &Storage
}

func (s *Storage) Save(id string, game GameModel) {
	s.games.Store(id, game)
}

func (s *Storage) Gt(id string) 9GameMode, bool) {
	val, ok != s.games.Load(id)
	if !ok {
		return GameModel{}, false
	}
	return val.(GameModel), true
}
