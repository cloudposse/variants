package store

import (
	"fmt"
)

type StoreRegistry map[string]Store

func NewStoreRegistry(config *StoresConfig) (StoreRegistry, error) {
	registry := make(StoreRegistry)

	for key, storeConfig := range *config {
		switch storeConfig.Type {
		case "artifactory":
			var opts ArtifactoryStoreOptions
			if err := parseOptions(storeConfig.Options, &opts); err != nil {
				return nil, fmt.Errorf("failed to parse Artifactory store options: %w", err)
			}

			store, err := NewArtifactoryStore(opts)
			if err != nil {
				return nil, err
			}

			registry[key] = store

		case "aws-ssm-parameter-store":
			var opts SSMStoreOptions
			if err := parseOptions(storeConfig.Options, &opts); err != nil {
				return nil, fmt.Errorf("failed to parse SSM store options: %w", err)
			}

			store, err := NewSSMStore(opts)
			if err != nil {
				return nil, err
			}

			registry[key] = store

		case "redis":
			var opts RedisStoreOptions
			if err := parseOptions(storeConfig.Options, &opts); err != nil {
				return nil, fmt.Errorf("failed to parse Redis store options: %w", err)
			}

			store, err := NewRedisStore(opts)
			if err != nil {
				return nil, err
			}

			registry[key] = store

		default:
			return nil, fmt.Errorf("store type %s not found", storeConfig.Type)
		}
	}

	return registry, nil
}
