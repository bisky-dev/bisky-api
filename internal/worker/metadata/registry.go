package metadata

import "fmt"

func NewRegistry(providers map[ProviderName]Provider) *Registry {
	copied := make(map[ProviderName]Provider, len(providers))
	for name, provider := range providers {
		copied[name] = provider
	}
	return &Registry{providers: copied}
}

func (r *Registry) Provider(name ProviderName) (Provider, error) {
	provider, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("metadata provider %q is not supported", name)
	}
	return provider, nil
}
