package decorators

import (
	"context"
	log "github.com/sirupsen/logrus"
	"strings"

	"export-nft-data/domain"
	"github.com/ethereum/go-ethereum/common"
)

// DeployInfo marks the deployment block of each address
func DeployInfo(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	logger := log.WithFields(log.Fields{
		"step":  "DeployInfo",
		"count": len(cs),
	})

	logger.Debug("start")
	var addresses []string
	for _, c := range cs {
		if c.DeployBlock == nil {
			addresses = append(addresses, c.Address.Hex())
		}
	}
	creations, err := cfg.EtherScan.GetContractCreations(ctx, addresses)
	if err != nil {
		return err
	}
	for _, creation := range creations {
		for _, c := range cs {
			if strings.EqualFold(c.Address.Hex(), creation.ContractAddress) {
				tr, err := cfg.Eth.TransactionReceipt(ctx, common.HexToHash(creation.TxHash))
				if err != nil {
					return err
				}
				c.DeployBlock = tr.BlockNumber
			}
		}
	}

	logger.Debug("complete")
	return nil
}
