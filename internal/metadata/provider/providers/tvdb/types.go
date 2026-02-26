package tvdb

import (
	"net/http"
	"sync"
	"time"
)

type Provider struct {
	baseURL string
	apiKey  string
	pin     string
	client  *http.Client
	debug   bool

	mu         sync.Mutex
	token      string
	tokenUntil time.Time
}
