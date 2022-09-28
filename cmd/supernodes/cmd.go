package supernodes

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"export-nft-data/client/centerdev"
	"export-nft-data/client/eth"
	"export-nft-data/client/etherscan"
	"export-nft-data/domain"
	"export-nft-data/utils"
)

type Config struct {
	File         string
	CenterDevKey string
	EtherscanKey string
	JsonRpcUrl   string
}

func Run(ctx context.Context, cfg Config) error {
	b, err := os.ReadFile("./seed-names.csv")
	if err != nil {
		return err
	}
	etherum, err := eth.NewEthClient(cfg.JsonRpcUrl)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	cl := centerdev.NewClient(cfg.CenterDevKey)

	var output []*domain.Collection

	var collections []*centerdev.Collection
	var addresses []string
	results := make(map[string][]*centerdev.Collection)
	for _, l := range lines {
		cs, err := cl.GetCollections(ctx, centerdev.NetworkEthereumMainnet, l)
		if err != nil {
			return err
		}
		collections = append(collections, cs...)
		for _, c := range cs {
			addresses = append(addresses, c.Address)
			output = append(output, &domain.Collection{
				Name:    c.Name,
				Address: common.HexToAddress(c.Address),
			})
		}

		results[l] = cs
	}

	e := etherscan.NewEtherscanClient(cfg.EtherscanKey)
	creations, err := e.GetContractCreations(ctx, addresses)
	if err != nil {
		return err
	}
	for _, c := range creations {
		for _, o := range output {
			if strings.EqualFold(o.Address.Hex(), c.ContractAddress) {
				tr, err := etherum.TransactionReceipt(ctx, common.HexToHash(c.TxHash))
				if err != nil {
					return err
				}
				o.DeployBlock = tr.BlockNumber
			}
		}
	}

	fmt.Println(utils.ToJson(output))
	return nil
}
