package booking

import "github.com/jackc/pgx/v5/pgtype"

type Booking struct {
	ID        uint64             `json:"id,string"`
	EventID   uint64             `json:"event_id,string"`
	StartAt   pgtype.Timestamptz `json:"start_at"`
	EndedAt   pgtype.Timestamptz `json:"ended_at"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}
