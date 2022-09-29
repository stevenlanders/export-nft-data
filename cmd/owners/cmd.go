package owners

import (
	"context"
	"encoding/json"
	"export-nft-data/client/eth"
	"export-nft-data/events"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
	"os"

	log "github.com/sirupsen/logrus"

	"export-nft-data/domain"
	"export-nft-data/utils"
)

// seconds per block
const mainnetRate = 13

type Config struct {
	File       string
	JsonRpcUrl string
	Days       int
}

func Run(ctx context.Context, cfg Config) error {
	b, err := os.ReadFile(cfg.File)
	if err != nil {
		return err
	}
	var collections []*domain.Collection
	if err := json.Unmarshal(b, &collections); err != nil {
		return err
	}

	et, err := eth.NewEthClient(cfg.JsonRpcUrl)
	if err != nil {
		return err
	}

	latest, err := et.BlockByNumber(ctx, nil)
	if err != nil {
		log.WithError(err).Error("error getting latest block number")
		return err
	}

	s, err := events.NewStream(events.Config{
		JsonRpcUrl: cfg.JsonRpcUrl,
	})
	if err != nil {
		return err
	}

	seconds := cfg.Days * 24 * 60 * 60
	blockOffset := big.NewInt(int64(math.Floor(float64(seconds / mainnetRate))))

	for _, c := range collections {
		targetBlock := big.NewInt(0)
		targetBlock.Add(blockOffset, c.DeployBlock)
		if targetBlock.Cmp(latest.Number()) == 1 {
			targetBlock = latest.Number()
		}
		tb := targetBlock.Uint64()
		var owners []common.Address
		err := s.ForEachOwner(ctx, &events.OwnerFilter{
			BlockFilter: events.BlockFilter{
				StartBlock: c.DeployBlock.Uint64(),
				EndBlock:   tb,
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

	fmt.Println(utils.ToJson(collections))
	return nil
}
