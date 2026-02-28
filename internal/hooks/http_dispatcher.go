package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const hookRequestTimeout = 5 * time.Second

type HTTPDispatcher struct {
	pool   *pgxpool.Pool
	client *http.Client
}

type dispatchBody struct {
	Event     Event     `json:"event"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

func NewHTTPDispatcher(pool *pgxpool.Pool) (*HTTPDispatcher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := EnsureStore(ctx, pool); err != nil {
		return nil, err
	}

	return &HTTPDispatcher{
		pool: pool,
		client: &http.Client{
			Timeout: hookRequestTimeout,
		},
	}, nil
}

func (d *HTTPDispatcher) DispatchPre(ctx context.Context, event Event, payload any) error {
	return d.dispatch(ctx, event, payload)
}

func (d *HTTPDispatcher) DispatchPost(ctx context.Context, event Event, payload any) {
	if err := d.dispatch(ctx, event, payload); err != nil {
		log.Printf("hook dispatch failed event=%s err=%v", event, err)
	}
}

func (d *HTTPDispatcher) dispatch(ctx context.Context, event Event, payload any) error {
	if !IsValidEvent(event) {
		return fmt.Errorf("invalid hook event %q", event)
	}

	var targetURL string
	if err := d.pool.QueryRow(ctx, `
SELECT target_url
FROM hook_settings
WHERE event_name = $1
`, string(event)).Scan(&targetURL); err != nil {
		return err
	}
	targetURL = strings.TrimSpace(targetURL)
	if targetURL == "" {
		return nil
	}

	bodyBytes, err := json.Marshal(dispatchBody{
		Event:     event,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hook-Event", string(event))

	res, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("hook endpoint returned status %d", res.StatusCode)
	}
	return nil
}
