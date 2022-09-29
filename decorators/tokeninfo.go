package decorators

import (
	"context"
	"export-nft-data/client/centerdev"
	"export-nft-data/domain"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func TokenInfo(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	logger := log.WithFields(log.Fields{
		"step":  "TokenInfo",
		"count": len(cs),
	})
	logger.Debug("start")
	defer logger.Debug("complete")

	grp, ctx := errgroup.WithContext(ctx)
	ch := make(chan *domain.Collection)

	for i := 0; i < 10; i++ {
		grp.Go(func() error {
			for c := range ch {
				if c.Name != "" {
					continue
				}
				logger.WithFields(log.Fields{
					"collection": c.Name,
					"address":    c.Address,
				}).Debug("getting token info")
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
		})
	}

	grp.Go(func() error {
		defer close(ch)
		for _, c := range cs {
			ch <- c
		}
		return nil
	})

	return grp.Wait()
}
