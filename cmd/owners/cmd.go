package owners

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"

	log "github.com/sirupsen/logrus"

	"export-nft-data/client/alchemy"
	"export-nft-data/client/eth"
	"export-nft-data/domain"
	"export-nft-data/utils"
)

// seconds per block
const mainnetRate = 13

type Config struct {
	File       string
	AlchemyKey string
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

	et, err := eth.NewEthClient(fmt.Sprintf("https://eth-mainnet.alchemyapi.io/v2/%s", cfg.AlchemyKey))
	if err != nil {
		log.WithError(err).Error("error creating eth client")
		return err
	}

	latest, err := et.BlockByNumber(ctx, nil)
	if err != nil {
		log.WithError(err).Error("error getting latest block number")
		return err
	}

	a := alchemy.NewNFTClient(cfg.AlchemyKey)
	seconds := cfg.Days * 24 * 60 * 60
	blockOffset := big.NewInt(int64(math.Floor(float64(seconds / mainnetRate))))

	for _, c := range collections {
		targetBlock := big.NewInt(0)
		targetBlock.Add(blockOffset, c.DeployBlock)
		if targetBlock.Cmp(latest.Number()) == 1 {
			targetBlock = latest.Number()
		}
		owners, err := a.GetOwners(ctx, c.Address, targetBlock)
		if err != nil {
			return err
		}
		c.OwnerBlock = targetBlock
		c.Owners = owners
	}

	fmt.Println(utils.ToJson(collections))
	return nil
}
