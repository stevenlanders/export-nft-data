package decorators

import (
	"context"
	"export-nft-data/domain"
	"export-nft-data/events"
	log "github.com/sirupsen/logrus"
)

var tokenIgnoreList = []string{
	"0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85", // ignore ENS
}

func Edges(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	logger := log.WithFields(log.Fields{
		"step":  "Edges",
		"count": len(cs),
	})

	logger.Debug("start")

	var edges []*domain.CollectionEdge

	// build lookup maps
	knownCollections := make(map[string]bool)
	buyerCollections := make(map[string][]*domain.Collection)
	for _, c := range cs {
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

	eb := cfg.StartBlock.Uint64() + cfg.NumBlocks.Uint64()
	err := cfg.Stream.ForEachCollectionOrder(ctx, &events.OrderFilter{
		BlockFilter: events.BlockFilter{
			StartBlock: cfg.StartBlock.Uint64(),
			EndBlock:   eb,
		},
		IgnoreTokens: tokenIgnoreList,
	}, func(o *events.CollectionOrder) error {

		// if not buyer we care about, move on
		if _, ok := buyerCollections[o.Buyer.Hex()]; !ok {
			return nil
		}
		// if already known, move on
		// commenting out to allow revisits
		//if _, ok := knownCollections[o.Collection.Hex()]; ok {
		//	return nil
		//}

		cs := buyerCollections[o.Buyer.Hex()]

		for _, c := range cs {
			c.Edges = append(c.Edges, &domain.CollectionEdge{
				FromCollection: &c.Address,
				ToCollection:   &o.Collection,
				Buyer:          &o.Buyer,
				Price:          o.Price,
			})
		}
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("error getting edges")
		return err
	}

	if edges == nil {
		edges = []*domain.CollectionEdge{}
	}
	logger.Debug("complete")

	return nil
}
