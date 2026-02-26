package metadata

import "context"

func NewService(registry *Registry) *Service {
	return &Service{registry: registry}
}

func (s *Service) Search(ctx context.Context, providerName ProviderName, query string, opts SearchOpts) ([]SearchHit, error) {
	provider, err := s.registry.Provider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Search(ctx, query, opts)
}

func (s *Service) GetShow(ctx context.Context, providerName ProviderName, externalID string) (Show, error) {
	provider, err := s.registry.Provider(providerName)
	if err != nil {
		return Show{}, err
	}
	return provider.GetShow(ctx, externalID)
}

func (s *Service) ListEpisodes(ctx context.Context, providerName ProviderName, externalID string, opts ListEpisodesOpts) ([]Episode, error) {
	provider, err := s.registry.Provider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.ListEpisodes(ctx, externalID, opts)
}
