package secrets

import (
	"context"
	"fmt"
)

type vaultClient interface {
	ReadSecret(
		ctx context.Context,
		path string,
		dst any,
	) error
}

type Paths struct {
	Postgres string
	Hasher   string
}

type SecretsProvider struct {
	client vaultClient
	paths  Paths
}

func New(client vaultClient, paths Paths) *SecretsProvider {
	return &SecretsProvider{
		client: client,
		paths:  paths,
	}
}

func (sp *SecretsProvider) Load(ctx context.Context) (*Secrets, error) {
	var secrets Secrets

	secretLoaders := []struct {
		name string
		path string
		dst  any
	}{
		{
			name: "postgres",
			path: sp.paths.Postgres,
			dst:  &secrets.Postgres,
		},
		{
			name: "hasher",
			path: sp.paths.Hasher,
			dst:  &secrets.Hasher,
		},
	}

	for _, loader := range secretLoaders {
		if err := sp.client.ReadSecret(
			ctx,
			loader.path,
			loader.dst,
		); err != nil {
			return nil, fmt.Errorf(
				"reading %s secrets: %w",
				loader.name,
				err,
			)
		}
	}

	return &secrets, nil
}
