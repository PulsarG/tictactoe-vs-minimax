package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"tictactoe-vs-minimax/internal/domain"
)

type Server struct {
	storage domain.Storage
}

func NewServer(storage domain.Storage) *Server {
	return &Server{storage: storage}
}

func NewRouter() *gin.Engine {
	return gin.Default()
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.POST("/game", s.createGameHandler)
	r.POST("/game/:id", s.nextMoveHandler)
	r.GET("/game/:id", s.getUserHandler)
}

func (s *Server) createGameHandler(c *gin.Context) {
	var input InputData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ncorrect board" + err.Error()})
		return
	}

	emptyBoard := make([]int, 9)
	if !domain.IsValidMove(emptyBoard, input.Board) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect move"})
		return
	}

	board := make([]int, 9)
	copy(board, input.Board)

	domain.MakeServerMove(board)
	status := domain.GetGameStatus(board)

	user, err := s.storage.Create(board, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error save game"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":  "Game start",
		"game": user,
	})
	printBoard(c, user.Board)

}

func (s *Server) nextMoveHandler(c *gin.Context) {
	userID := c.Param("id")

	var input InputData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect board"})
		return
	}

	user, err := s.storage.Get(userID)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game " + userID + " not found"})
		return
	}

	if user.Status != "play" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You: " + user.Status,
			"game":  user,
		})
		return
	}

	if !domain.IsValidMove(user.Board, input.Board) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Cheat move - error",
			"current_board": user.Board,
		})
		return
	}

	copy(user.Board, input.Board)
	user.Status = domain.GetGameStatus(user.Board)

	if user.Status == "play" {
		domain.MakeServerMove(user.Board)
		user.Status = domain.GetGameStatus(user.Board)
	}

	if err := s.storage.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": user,
	})
	printBoard(c, user.Board)
}

func (s *Server) getUserHandler(c *gin.Context) {
	userID := c.Param("id")

	user, err := s.storage.Get(userID)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game " + userID + " not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func printBoard(c *gin.Context, b []int) {
	for i := range 9 {
		n := ""
		if i == 0 || i == 3 || i == 6 {
			n = "\n"
		}
		cell := "."
		switch b[i] {
		case 1:
			cell = "X"
		case 2:
			cell = "O"
		}
		c.String(http.StatusOK, fmt.Sprintf(n+"%s ", cell))
	}
}
