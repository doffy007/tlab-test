CREATE TABLE IF NOT EXISTS bookings (
    id BIGINT PRIMARY KEY,
    event_id BIGINT,
    start_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ 
)