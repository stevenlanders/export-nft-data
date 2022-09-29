package discover

import (
	"context"
	"encoding/json"
	"export-nft-data/client/centerdev"
	"export-nft-data/client/eth"
	"export-nft-data/client/etherscan"
	"export-nft-data/decorators"
	"export-nft-data/domain"
	"export-nft-data/events"
	"export-nft-data/utils"
	"fmt"
	"math/big"
	"os"
)

type Config struct {
	File           string
	JsonRpcUrl     string
	CenterDevKey   string
	EtherscanKey   string
	BlockStart     uint64
	NumberOfBlocks uint64
	OwnerDays      int
	Iterations     int
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

	s, err := events.NewStream(events.Config{JsonRpcUrl: cfg.JsonRpcUrl})
	if err != nil {
		return err
	}

	ec, err := eth.NewEthClient(cfg.JsonRpcUrl)
	if err != nil {
		return err
	}
	cd := centerdev.NewClient(cfg.CenterDevKey)

	e := etherscan.NewEtherscanClient(cfg.EtherscanKey)

	err = decorators.RunDecorators(ctx, collections, decorators.Config{
		StartBlock: big.NewInt(int64(cfg.BlockStart)),
		NumBlocks:  big.NewInt(int64(cfg.NumberOfBlocks)),
		Days:       cfg.OwnerDays,
		Stream:     s,
		Eth:        ec,
		EtherScan:  e,
		CenterDev:  cd,
		Iterations: cfg.Iterations,
	})

	if err != nil {
		return nil
	}

	fmt.Println(utils.ToJson(collections))
	return nil
}