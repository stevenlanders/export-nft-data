package etherscan

import (
	"context"
	"export-nft-data/client/utils"
	u "export-nft-data/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const apiUrl = "https://api.etherscan.io/api?module=contract&action=getcontractcreation&contractaddresses=%s&apikey=%s"

type Client interface {
	GetContractCreations(ctx context.Context, addresses []string) ([]ContractCreation, error)
}

type client struct {
	key string
}

func (c *client) GetContractCreations(ctx context.Context, addresses []string) ([]ContractCreation, error) {
	var result []ContractCreation
	chunks := u.ChunkBy[string](addresses, 5)
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}
		addressList := strings.Join(chunk, ",")
		url := fmt.Sprintf(apiUrl, addressList, c.key)
		var page ContractCreationResponse
		if err := utils.Get(ctx, url, nil, &page); err != nil {
			log.WithFields(log.Fields{
				"url": url,
			}).WithError(err).Error("error getting contract")
			return nil, err
		}
		time.Sleep(1 * time.Second)
		result = append(result, page.Result...)
	}
	return result, nil
}

func NewEtherscanClient(key string) Client {
	return &client{
		key: key,
	}
}
