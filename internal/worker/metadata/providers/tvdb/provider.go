package tvdb

import (
	"context"
	"fmt"

	"github.com/keithics/devops-dashboard/api/internal/worker/metadata"
)

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Search(ctx context.Context, query string, opts metadata.SearchOpts) ([]metadata.SearchHit, error) {
	_ = ctx
	_ = query
	_ = opts
	return nil, fmt.Errorf("tvdb provider Search is not implemented")
}

func (p *Provider) GetShow(ctx context.Context, externalID string) (metadata.Show, error) {
	_ = ctx
	_ = externalID
	return metadata.Show{}, fmt.Errorf("tvdb provider GetShow is not implemented")
}

func (p *Provider) ListEpisodes(ctx context.Context, externalID string, opts metadata.ListEpisodesOpts) ([]metadata.Episode, error) {
	_ = ctx
	_ = externalID
	_ = opts
	return nil, fmt.Errorf("tvdb provider ListEpisodes is not implemented")
}
