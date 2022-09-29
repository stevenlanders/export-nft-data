package decorators

import (
	"context"
	"export-nft-data/domain"
	"export-nft-data/events"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"math"
	"math/big"
)

// seconds per block
const mainnetRate = 13

func Owners(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	logger := log.WithFields(log.Fields{
		"step": "Owners",
	})
	logger.Debug("start")

	seconds := cfg.Days * 24 * 60 * 60
	blockOffset := big.NewInt(int64(math.Floor(float64(seconds / mainnetRate))))

	latest, err := cfg.Eth.BlockByNumber(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("error getting latest block number")
		return err
	}

	grp, ctx := errgroup.WithContext(ctx)
	ch := make(chan *domain.Collection)

	for i := 0; i < 5; i++ {
		grp.Go(func() error {
			for c := range ch {
				targetBlock := big.NewInt(0)
				targetBlock.Add(blockOffset, c.DeployBlock)
				if targetBlock.Cmp(latest.Number()) == 1 {
					targetBlock = latest.Number()
				}
				tb := targetBlock.Uint64()
				var owners []common.Address

				err = cfg.Stream.ForEachOwner(ctx, &events.OwnerFilter{
					BlockFilter: events.BlockFilter{
						StartBlock: c.DeployBlock.Uint64(),
						EndBlock:   tb,
					},
					Token: c.Address,
				}, func(owner common.Address) error {
					owners = append(owners, owner)
					return nil
				})

				logger.WithFields(log.Fields{
					"collection": c.Name,
					"address":    c.Address.Hex(),
					"owners":     len(owners),
				}).Debug("found owners")

				if err != nil {
					log.WithError(err).Error("error in ForEachOwner")
					return err
				}
				c.OwnerBlock = targetBlock
				c.Owners = owners
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

	if err := grp.Wait(); err != nil {
		return err
	}

	logger.Debug("complete")

	return nil
}
