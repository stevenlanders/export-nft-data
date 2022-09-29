package decorators

import (
	"context"
	"export-nft-data/client/centerdev"
	"export-nft-data/domain"
)

func TokenInfo(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	for _, c := range cs {
		collection, err := cfg.CenterDev.GetCollection(ctx, centerdev.NetworkEthereumMainnet, c.Address.Hex())
		if err != nil {
			continue
		}
		if collection == nil {
			continue
		}
		c.Name = collection.Name
	}
	return nil
}
