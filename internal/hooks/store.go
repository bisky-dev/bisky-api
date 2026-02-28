package hooks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Event     Event     `json:"event"`
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var allEvents = []Event{
	EventShowCreatePre,
	EventShowCreatePost,
	EventShowUpdatePre,
	EventShowUpdatePost,
	EventShowDeletePre,
	EventShowDeletePost,
	EventEpisodeCreatePre,
	EventEpisodeCreatePost,
	EventEpisodeUpdatePre,
	EventEpisodeUpdatePost,
	EventEpisodeDeletePre,
	EventEpisodeDeletePost,
}

func AllEvents() []Event {
	out := make([]Event, len(allEvents))
	copy(out, allEvents)
	return out
}

func IsValidEvent(event Event) bool {
	for _, candidate := range allEvents {
		if candidate == event {
			return true
		}
	}
	return false
}

func EnsureStore(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, event := range allEvents {
		if _, err := tx.Exec(ctx, `
INSERT INTO hook_settings (event_name, target_url)
VALUES ($1, '')
ON CONFLICT (event_name) DO NOTHING
`, string(event)); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func ListConfigs(ctx context.Context, pool *pgxpool.Pool) ([]Config, error) {
	rows, err := pool.Query(ctx, `
SELECT event_name, target_url, updated_at
FROM hook_settings
ORDER BY event_name ASC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Config, 0, len(allEvents))
	for rows.Next() {
		var item Config
		var eventName string
		if err := rows.Scan(&eventName, &item.URL, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.Event = Event(eventName)
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func UpsertConfig(ctx context.Context, pool *pgxpool.Pool, event Event, rawURL string) (Config, error) {
	if !IsValidEvent(event) {
		return Config{}, fmt.Errorf("invalid hook event %q", event)
	}

	url := strings.TrimSpace(rawURL)
	var updated Config
	var eventName string
	err := pool.QueryRow(ctx, `
INSERT INTO hook_settings (event_name, target_url, updated_at)
VALUES ($1, $2, now())
ON CONFLICT (event_name)
DO UPDATE SET target_url = EXCLUDED.target_url, updated_at = now()
RETURNING event_name, target_url, updated_at
`, string(event), url).Scan(&eventName, &updated.URL, &updated.UpdatedAt)
	if err != nil {
		return Config{}, err
	}

	updated.Event = Event(eventName)
	return updated, nil
}
