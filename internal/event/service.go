package event

import (
	"context"
	"time"

	db "github.com/doffy007/tlab-test/internal/database"
	"github.com/doffy007/tlab-test/internal/uid"
	"github.com/jackc/pgx/v5"
)

type EventService interface {
	CreateEvent(req *Event) (*Event, error)
	GetEvent(id uint64) (*Event, error)
}

type srv struct {
	db db.DBService
}

var (
	Service EventService
)

func init() {
	Service = New(db.Service)
}

func New(db db.DBService) EventService {
	return &srv{db}
}

func (s *srv) CreateEvent(req *Event) (*Event, error) {
	if req.ID == 0 {
		id, err := uid.New()
		if err != nil {
			return nil, err
		}
		req.ID = id
	}

	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO "event" (
			id,
			name,
			description,
			quota,
			start_at,
			ended_at,
			created_at,
			updated_at,
			deleted_at
		)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		req.ID,
		req.Name,
		req.Description,
		req.Quota,
		req.StartAt,
		req.EndedAt,
		time.Now(),
		time.Now(),
		req.DeletedAt,
	)
	if err != nil {
		tx.Rollback(context.Background())
		return nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *srv) GetEvent(id uint64) (*Event, error) {
	var result Event

	if err := s.db.QueryRow(
		`SELECT
			id,
			name,
			description,
			quota,
			start_at,
			ended_at,
			created_at,
			updated_at
			FROM "event"
			WHERE id = $1 AND deleted_at is null`,
		id,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.Quota,
		&result.StartAt,
		&result.EndedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
