package datasource

import (
	"sync"

	"tictactoe-vs-minimax/internal/domain"

	"github.com/google/uuid"
)

type ServerStorage struct {
	mu     sync.RWMutex
	lastID int
	users  map[string]domain.UserData
}

func NewServerStorage() *ServerStorage {
	return &ServerStorage{
		users: make(map[string]domain.UserData),
	}
}

func (m *ServerStorage) Create(board []int, status string) (domain.UserData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastID++
	// newID := strconv.Itoa(m.lastID)
	newID := uuid.NewString()

	boardCopy := make([]int, len(board))
	copy(boardCopy, board)

	user := domain.UserData{
		ID:     newID,
		Board:  boardCopy,
		Status: status,
	}
	m.users[newID] = user
	return user, nil
}

func (m *ServerStorage) Get(id string) (domain.UserData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return domain.UserData{}, domain.ErrNotFound
	}

	boardCopy := make([]int, len(user.Board))
	copy(boardCopy, user.Board)
	user.Board = boardCopy

	return user, nil
}

func (m *ServerStorage) Update(user domain.UserData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID]; !exists {
		return domain.ErrNotFound
	}

	boardCopy := make([]int, len(user.Board))
	copy(boardCopy, user.Board)
	user.Board = boardCopy

	m.users[user.ID] = user
	return nil
}
