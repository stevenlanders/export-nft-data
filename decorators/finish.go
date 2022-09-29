package decorators

import (
	"context"
	"export-nft-data/domain"
)

func MarkProcessed(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	for _, c := range cs {
		c.Processed = true
	}
	return nil
}