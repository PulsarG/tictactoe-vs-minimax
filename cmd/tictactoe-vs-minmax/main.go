package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var ErrNotFound = errors.New("Data not found")

type UserData struct {
	ID     string `json:"id"`
	Board  []int  `json:"board"`
	Status string `json:"status"`
}

type InputData struct {
	Board []int `json:"board" binding:"required,len=9"`
}

type Storage interface {
	Create(b []int, status string) (UserData, error)
	Get(id string) (UserData, error)
	Update(u UserData) error
}

type ServerStorage struct {
	mu     sync.RWMutex
	lastID int
	users  map[string]UserData
}

func NewServerStorage() *ServerStorage {
	return &ServerStorage{
		users: make(map[string]UserData),
	}
}

func (ss *ServerStorage) Create(b []int, status string) (UserData, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.lastID++
	newID := strconv.Itoa(ss.lastID)
	bCopy := make([]int, len(b))
	copy(bCopy, b)
	user := UserData{
		ID:     newID,
		Board:  bCopy,
		Status: status,
	}
	ss.users[newID] = user
	return user, nil
}

func (ss *ServerStorage) Get(id string) (UserData, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	user, exists := ss.users[id]
	if !exists {
		return UserData{}, ErrNotFound
	}
	bCopy := make([]int, len(user.Board))
	copy(bCopy, user.Board)
	user.Board = bCopy
	return user, nil
}

func (ss *ServerStorage) Update(u UserData) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, exists := ss.users[u.ID]; !exists {
		return ErrNotFound
	}
	bCopy := make([]int, len(u.Board))
	copy(bCopy, u.Board)
	u.Board = bCopy
	ss.users[u.ID] = u
	return nil
}

type Server struct {
	storage Storage
}

func NewServer(storage Storage) *Server {
	return &Server{storage: storage}
}

func NewRouter() *gin.Engine {
	return gin.Default()
}

func startServer(lc fx.Lifecycle, r *gin.Engine, srv *Server) {
	r.POST("/user", srv.createGameHandler)
	r.POST("/user/:id", srv.nextMoveHandler)
	r.GET("/user/:id", srv.getUserHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("UberFX start server :8080")
			go func() {
				if err := httpServer.ListenAndServe(); err != nil {
					log.Fatalf("Error start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("UberFx stopped server")
			return httpServer.Shutdown(ctx)
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			fx.Annotate(NewServerStorage, fx.As(new(Storage))),
			NewServer,
			NewRouter,
		),
		fx.Invoke(startServer),
	).Run()
}

func (s *Server) createGameHandler(c *gin.Context) {
	var input InputData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect data" + err.Error()})
		return
	}
	empBoard := make([]int, 9)
	if !isValidMove(empBoard, input.Board) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not valid board"})
		return
	}
	b := make([]int, 9)
	copy(b, input.Board)
	makeServerMove(b)
	status := getGameStatus(b)

	user, err := s.storage.Create(b, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error save game"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Session start",
		"user":    user,
	})
	printBoard(c, user.Board)
}

func (s *Server) nextMoveHandler(c *gin.Context) {
	uID := c.Param("id")
	var input InputData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect data" + err.Error()})
		return
	}

	user, err := s.storage.Get(uID)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game id " + uID + " not found"})
		return
	}

	if user.Status != "play" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Game Over! You " + user.Status,
			"current_board": user.Board,
		})
		return
	}

	if !isValidMove(user.Board, input.Board) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User id: " + uID + " cheat",
		})
		return
	}

	copy(user.Board[:], input.Board)
	user.Status = getGameStatus(user.Board)
	if user.Status == "play" {
		makeServerMove(user.Board)
		user.Status = getGameStatus(user.Board)
	}

	if err := s.storage.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data saved",
		"result":  user,
	})

	printBoard(c, user.Board)
	if isBoardFull(user.Board) {
		c.String(http.StatusOK, "Game Over! Board is full!")
	}
}

func printBoard(c *gin.Context, b []int) {
	for i := range 9 {
		n := ""
		if i == 0 || i == 3 || i == 6 {
			n = "\n"
		}
		c.String(http.StatusOK, fmt.Sprintf(n+"%d ", b[i]))
	}
}

func getGameStatus(b []int) string {
	winner := checkWinner(b)
	if winner == 1 {
		return "win"
	}
	if winner == 2 {
		return "lose"
	}
	if isBoardFull(b) {
		return "draw"
	}
	return "play"
}

func (s *Server) getUserHandler(c *gin.Context) {
	uID := c.Param("id")

	user, err := s.storage.Get(uID)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gamae id " + uID + " not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func makeServerMove(b []int) {
	bestScore := -1000
	bestMove := -1

	for i := range 9 {
		if b[i] == 0 {
			b[i] = 2
			score := minimax(b, 0, false)
			b[i] = 0

			if score > bestScore {
				bestScore = score
				bestMove = i
			}
		}
	}

	if bestMove != -1 {
		b[bestMove] = 2
	}
}

func minimax(b []int, depth int, isMaximizing bool) int {
	winner := checkWinner(b)
	if winner == 2 {
		return 10 - depth
	}
	if winner == 1 {
		return depth - 10
	}
	if isBoardFull(b) {
		return 0
	}

	if isMaximizing {
		bestScore := -1000
		for i := range 9 {
			if b[i] == 0 {
				b[i] = 2
				score := minimax(b, depth+1, false)
				b[i] = 0
				if score > bestScore {
					bestScore = score
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for i := range 9 {
			if b[i] == 0 {
				b[i] = 1
				score := minimax(b, depth+1, true)
				b[i] = 0
				if score < bestScore {
					bestScore = score
				}
			}
		}
		return bestScore
	}
}

func checkWinner(b []int) int {
	lines := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, line := range lines {
		if b[line[0]] != 0 && b[line[0]] == b[line[1]] && b[line[1]] == b[line[2]] {
			return b[line[0]]
		}
	}
	return 0
}

func isBoardFull(b []int) bool {
	for _, cell := range b {
		if cell == 0 {
			return false
		}
	}
	return true
}

func isValidMove(oldB, newB []int) bool {
	newMovesCount := 0

	for i := range 9 {
		if oldB[i] == newB[i] {
			continue
		}

		if oldB[i] == 0 && newB[i] == 1 {
			newMovesCount++
		} else {
			return false
		}
	}
	return newMovesCount == 1
}
