package decorators

import (
	"context"
	"export-nft-data/domain"
	log "github.com/sirupsen/logrus"
)

func MarkProcessed(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	logger := log.WithFields(log.Fields{
		"step":  "MarkProcessed",
		"count": len(cs),
	})
	logger.Debug("start")
	for _, c := range cs {
		c.Processed = true
	}
	logger.Debug("complete")
	return nil
}
