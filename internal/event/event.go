package event

import "github.com/jackc/pgx/v5/pgtype"

type Event struct {
	ID          uint64             `json:"id,string"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Quota       int                `json:"quota"`
	StartAt     pgtype.Timestamptz `json:"start_at"`
	EndedAt     pgtype.Timestamptz `json:"ended_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
}
