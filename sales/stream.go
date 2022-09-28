package sales

import (
	"context"
	"errors"
	"export-nft-data/client/eth"
	"export-nft-data/contracts/seaport"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strings"
)

const (
	defaultAddress = "0x00000000006c3852cbef3e08e8df289169ede581"
	ether          = "0x0000000000000000000000000000000000000000"
)

type BlockFilter struct {
	StartBlock uint64
	EndBlock   *uint64
}

type OrderFilter struct {
	BlockFilter
	IgnoreTokens []string
}

type Stream interface {
	ForEachCollectionOrder(ctx context.Context, of *OrderFilter, handler func(o *CollectionOrder) error) error
}

type stream struct {
	f *seaport.SeaportFilterer
}

func toMap(l []string) map[string]bool {
	res := make(map[string]bool)
	for _, item := range l {
		res[strings.ToLower(item)] = true
	}
	return res
}

func calculatePrice(items []seaport.ReceivedItem) (*big.Int, error) {
	price := big.NewInt(0)
	for _, si := range items {
		// only consider ether payments
		if si.Token.Hex() != ether {
			return nil, errors.New("not ether")
		}
		price = price.Add(price, si.Amount)
	}
	return price, nil
}

func (s *stream) ForEachCollectionOrder(ctx context.Context, of *OrderFilter, handler func(o *CollectionOrder) error) error {
	log.WithFields(log.Fields{
		"filter": of,
	}).Info("ForEachCollectionOrder")

	iterator, err := s.f.FilterOrderFulfilled(&bind.FilterOpts{
		Start:   of.StartBlock,
		End:     of.EndBlock,
		Context: ctx,
	}, nil, nil)

	if err != nil {
		return err
	}

	tokenIgnoreList := toMap(of.IgnoreTokens)

	for iterator.Next() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		o := iterator.Event

		// skip complicated orders (price isn't clear per item when bundled)
		if len(o.Offer) != 1 {
			return nil
		}
		var token common.Address
		for _, si := range o.Offer {
			if _, ok := tokenIgnoreList[strings.ToLower(si.Token.Hex())]; ok {
				return nil
			}
			token = si.Token
		}

		price, err := calculatePrice(o.Consideration)
		if err != nil {
			continue
		}

		if err := handler(&CollectionOrder{
			Buyer:      o.Recipient,
			Seller:     o.Offerer,
			Price:      price,
			Collection: token,
			TxHash:     o.Raw.TxHash,
		}); err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	JsonRpcUrl      string
	ContractAddress string
}

func NewStream(cfg Config) (Stream, error) {
	addr := defaultAddress
	if cfg.ContractAddress != "" {
		addr = cfg.ContractAddress
	}
	c, err := eth.NewEthClient(cfg.JsonRpcUrl)
	if err != nil {
		return nil, err
	}

	f, err := seaport.NewSeaportFilterer(common.HexToAddress(addr), c)
	if err != nil {
		return nil, err
	}

	return &stream{
		f: f,
	}, nil
}
