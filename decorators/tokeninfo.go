package decorators

import (
	"context"
	"export-nft-data/client/centerdev"
	"export-nft-data/domain"
	log "github.com/sirupsen/logrus"
)

func TokenInfo(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	logger := log.WithFields(log.Fields{
		"step": "TokenInfo",
	})
	logger.Debug("start")
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
	logger.Debug("complete")
	return nil
}
