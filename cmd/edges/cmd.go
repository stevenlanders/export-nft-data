package edges

import (
	"context"
	"encoding/json"
	"export-nft-data/domain"
	"export-nft-data/sales"
	"export-nft-data/utils"
	"fmt"
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

	s, err := sales.NewStream(sales.Config{JsonRpcUrl: cfg.JsonRpcUrl})
	if err != nil {
		return err
	}

	var edges []*domain.CollectionEdge

	// build lookup maps
	knownCollections := make(map[string]bool)
	buyerCollections := make(map[string][]*domain.Collection)
	for _, c := range collections {
		knownCollections[c.Address.Hex()] = true
		for _, o := range c.Owners {
			addr := o.Hex()
			if _, ok := buyerCollections[addr]; !ok {
				buyerCollections[addr] = []*domain.Collection{c}
			} else {
				buyerCollections[addr] = append(buyerCollections[addr], c)
			}
		}
	}

	endBlock := cfg.BlockStart + uint64(cfg.NumberOfBlocks)
	err = s.ForEachCollectionOrder(ctx, &sales.OrderFilter{
		BlockFilter: sales.BlockFilter{
			StartBlock: cfg.BlockStart,
			EndBlock:   &endBlock,
		},
		IgnoreTokens: tokenIgnoreList,
	}, func(o *sales.CollectionOrder) error {

		// if not buyer we care about, move on
		if _, ok := buyerCollections[o.Buyer.Hex()]; !ok {
			return nil
		}
		// if already known, move on
		if _, ok := knownCollections[o.Collection.Hex()]; ok {
			return nil
		}

		cs := buyerCollections[o.Buyer.Hex()]
		for _, c := range cs {
			edges = append(edges, &domain.CollectionEdge{
				FromCollection: c.Address,
				ToCollection:   o.Collection,
				Buyer:          o.Buyer,
				Price:          o.Price,
			})
		}
		return nil
	})
	if err != nil {
		return err
	}

	if edges == nil {
		edges = []*domain.CollectionEdge{}
	}

	fmt.Println(utils.ToJson(edges))
	return nil
}
