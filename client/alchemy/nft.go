package alchemy

import (
	"context"
	"export-nft-data/client/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const alchemyApi = "https://eth-mainnet.g.alchemy.com/nft/v2/%s/getOwnersForCollection"

type NFTClient interface {
	GetOwners(ctx context.Context, address common.Address, block *big.Int) ([]common.Address, error)
}

type nftClient struct {
	url string
}

func (c *nftClient) GetOwners(ctx context.Context, address common.Address, block *big.Int) ([]common.Address, error) {
	url := fmt.Sprintf("%s?contractAddress=%s&block=%d", c.url, address.Hex(), block.Int64())
	var result []common.Address
	var page OwnerResult
	if err := utils.Get(ctx, url, nil, &page); err != nil {
		return nil, err
	}
	result = append(result, page.OwnerAddresses...)

	for page.PageKey != "" {
		pageUrl := fmt.Sprintf("%s&pageKey=%s", url, page.PageKey)
		if err := utils.Get(ctx, pageUrl, nil, &page); err != nil {
			return nil, err
		}
		result = append(result, page.OwnerAddresses...)
	}

	return result, nil
}

func NewNFTClient(key string) NFTClient {
	return &nftClient{
		url: fmt.Sprintf(alchemyApi, key),
	}
}
