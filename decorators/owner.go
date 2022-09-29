package decorators

import (
	"context"
	"export-nft-data/domain"
	"export-nft-data/events"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
)

// seconds per block
const mainnetRate = 13

func Owners(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	seconds := cfg.Days * 24 * 60 * 60
	blockOffset := big.NewInt(int64(math.Floor(float64(seconds / mainnetRate))))

	latest, err := cfg.Eth.BlockByNumber(ctx, nil)
	if err != nil {
		log.WithError(err).Error("error getting latest block number")
		return err
	}

	for _, c := range cs {
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
				EndBlock:   &tb,
			},
			Token: c.Address,
		}, func(owner common.Address) error {
			owners = append(owners, owner)
			return nil
		})
		if err != nil {
			return err
		}
		c.OwnerBlock = targetBlock
		c.Owners = owners
	}
	return nil
}
