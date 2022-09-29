package edges

import (
	"context"
	"encoding/json"
	"export-nft-data/decorators"
	"export-nft-data/domain"
	"export-nft-data/events"
	"export-nft-data/utils"
	"fmt"
	"math/big"
	"os"
)

var tokenIgnoreList = []string{
	"0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85", //END
}

type Config struct {
	File           string
	JsonRpcUrl     string
	BlockStart     uint64
	NumberOfBlocks int
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

	err = decorators.Edges(ctx, collections, decorators.Config{
		Stream:     s,
		StartBlock: big.NewInt(int64(cfg.BlockStart)),
		NumBlocks:  big.NewInt(int64(cfg.NumberOfBlocks)),
	})
	if err != nil {
		return err
	}

	var edges []*domain.CollectionEdge

	for _, c := range collections {
		for _, e := range c.Edges {
			edges = append(edges, e)
		}
	}

	fmt.Println(utils.ToJson(edges))

	return nil
}
