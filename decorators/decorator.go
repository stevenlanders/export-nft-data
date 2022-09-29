package decorators

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"export-nft-data/client/centerdev"
	"export-nft-data/client/etherscan"
	"export-nft-data/domain"
	"export-nft-data/events"
)

type Config struct {
	// search space for new tokens
	StartBlock *big.Int
	NumBlocks  *big.Int
	// days to consider for owners
	Days int
	// prevents etherscan from getting too angry
	RateLimiter <-chan time.Time
	// data providers
	Stream    events.Stream
	Eth       *ethclient.Client
	EtherScan etherscan.Client
	CenterDev centerdev.Client
	// number of iterations on graph
	Iterations int
}

type Decorator func(ctx context.Context, c []*domain.Collection, cfg Config) error

var pipeline = []Decorator{
	TokenInfo,
	DeployInfo,
	Owners,
	Edges,
	MarkProcessed,
}

func filterUnprocessed(cs []*domain.Collection) []*domain.Collection {
	var result []*domain.Collection
	for _, c := range cs {
		if !c.Processed {
			result = append(result, c)
		}
	}
	return result
}

func addressMap(cs []*domain.Collection) map[string]bool {
	result := make(map[string]bool)
	for _, c := range cs {
		result[strings.ToLower(c.Address.Hex())] = true
	}
	return result
}

func RunDecorators(ctx context.Context, c []*domain.Collection, cfg Config) error {
	for i := 0; i < cfg.Iterations; i++ {
		processing := filterUnprocessed(c)
		for _, decorator := range pipeline {
			if err := decorator(ctx, processing, cfg); err != nil {
				return err
			}
		}

		// add new nodes to system here
		addresses := addressMap(c)
		for _, collection := range processing {
			for _, e := range collection.Edges {
				if _, ok := addresses[strings.ToLower(e.ToCollection.String())]; !ok {
					c = append(c, &domain.Collection{
						Address: e.ToCollection,
					})
					// avoid duplicates
					addresses[strings.ToLower(e.ToCollection.String())] = true
				}
			}
		}
	}
	// at end, fill in the data for the new nodes and stop
	unprocessed := filterUnprocessed(c)
	if err := TokenInfo(ctx, unprocessed, cfg); err != nil {
		return err
	}
	if err := DeployInfo(ctx, unprocessed, cfg); err != nil {
		return err
	}
	return nil
}
