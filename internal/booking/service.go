package booking

import (
	"context"
	"errors"
	"time"

	db "github.com/doffy007/tlab-test/internal/database"
	"github.com/doffy007/tlab-test/internal/event"
	"github.com/doffy007/tlab-test/internal/uid"
)

type BookingService interface {
	CreateBooking(req *Booking) (*Booking, error)
	GetQuotaEvent(id uint64) bool
}

type srv struct {
	db db.DBService
}

var (
	Service BookingService
)

func init() {
	Service = New(db.Service)
}

func New(db db.DBService) BookingService {
	return &srv{db}
}

func (s *srv) CreateBooking(req *Booking) (*Booking, error) {
	if req.ID == 0 {
		id, err := uid.New()
		if err != nil {
			return nil, err
		}
		req.ID = id
	}

	booking, err := event.Service.GetEvent(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	existQuota := s.GetQuotaEvent(req.EventID)
	if !existQuota {
		return nil, errors.New("quota is limited")
	}

	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO "event" (
			id,
			event_id,
			start_at,
			ended_at,
			create_at,
			updated_at,
			deleted_at
		)	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		req.ID,
		booking.ID,
		booking.StartAt,
		booking.EndedAt,
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

func (s *srv) GetQuotaEvent(id uint64) bool {
	booking, err := event.Service.GetEvent(id)
	if err != nil {
		return false
	}

	var quotaBooking int

	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM booking WHERE event_id = $1 `,
		id,
	).Scan(&quotaBooking); err != nil {
		return false
	}

	if quotaBooking == booking.Quota {
		return false
	}

	return true
}
