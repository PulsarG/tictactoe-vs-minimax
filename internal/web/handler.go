package web

import (
	"errors"
	"fmt"
	"net/http"

	"tictactoe-vs-minimax/internal/domain"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/game", s.createGameHandler)
	r.POST("/game/:id", s.nextMoveHandler)
	r.GET("/game/:id", s.getUserHandler)
}

// @Tags Game
// @Summary      Create game
// @Description  Creates a new game board
// @Accept       json
// @Produce      json
// @Param        input body InputData true "Game board data"
// @Success      201 {string} string "Created"
// @Failure      400 {object} map[string]string
// @Router       /game [post]
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

// @Tags Game
// @Summary      Next Move
// @Description  Move Minimax
// @Accept       json
// @Produce      json
// @Param        id     path      string      true  "ID игры"
// @Param        input  body      InputData  true  "Состояние доски"
// @Success      201 {string} string "Created"
// @Failure      400 {object} map[string]string
// @Router       /game/{id} [post]
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

// @Tags Game
// @Summary      Get Board
// @Description  Show Board from id
// @Accept       json
// @Produce      json
// @Param        id path string true  "ID игры"
// @Success      201 {string} string "Created"
// @Failure      400 {object} map[string]string
// @Router       /game/{id} [get]
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
