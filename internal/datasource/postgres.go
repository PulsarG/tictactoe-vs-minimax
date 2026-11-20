package datasource

import (
	"context"
	"errors"
	"fmt"
	"tictactoe-vs-minimax/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Error create pool pgx: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("DB not access: %w", err)
	}
	initQuery := `
	CREATE TABLE IF NOT EXISTS games (
		id UUID PRIMARY KEY,
		board INTEGER[] NOT NULL,
		status VARCHAR(20) NOT NULL
	);`

	if _, err := pool.Exec(ctx, initQuery); err != nil {
		return nil, fmt.Errorf("Error init table: %w", err)
	}
	return &PostgresStorage{pool: pool}, nil
}

func (p *PostgresStorage) Close() {
	p.pool.Close()
}

func (p *PostgresStorage) Create(board []int, status string) (domain.UserData, error) {
	newID := uuid.NewString()
	query := `INSERT INTO games (id, board, status) VALUES ($1, $2, $3)`
	_, err := p.pool.Exec(context.Background(), query, newID, board, status)
	if err != nil {
		return domain.UserData{}, err
	}
	return domain.UserData{
		ID:     newID,
		Board:  board,
		Status: status,
	}, nil
}

func (p *PostgresStorage) Get(id string) (domain.UserData, error) {
	query := `SELECT id::text, board, status FROM games WHERE id = $1`
	rows, err := p.pool.Query(context.Background(), query, id)
	if err != nil {
		return domain.UserData{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.UserData])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.UserData{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.UserData{}, err
	}
	return user, nil
}

func (p *PostgresStorage) Update(user domain.UserData) error {
	query := `UPDATE games SET board = $1, status = $2 WHERE id = $3`

	cmdTag, err := p.pool.Exec(context.Background(), query, user.Board, user.Status, user.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
