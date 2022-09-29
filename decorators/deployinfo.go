package decorators

import (
	"context"
	"strings"

	"export-nft-data/domain"
	"github.com/ethereum/go-ethereum/common"
)

// DeployInfo marks the deployment block of each address
func DeployInfo(ctx context.Context, cs []*domain.Collection, cfg Config) error {
	var addresses []string
	for _, c := range cs {
		addresses = append(addresses, c.Address.Hex())
	}
	creations, err := cfg.EtherScan.GetContractCreations(ctx, addresses)
	if err != nil {
		return err
	}
	for _, c := range cs {
		for _, creation := range creations {
			if strings.EqualFold(c.Address.Hex(), creation.ContractAddress) {
				tr, err := cfg.Eth.TransactionReceipt(ctx, common.HexToHash(creation.TxHash))
				if err != nil {
					return err
				}
				c.DeployBlock = tr.BlockNumber
			}
		}
	}

	return nil
}
