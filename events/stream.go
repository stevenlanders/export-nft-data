package events

import (
	"context"
	"errors"
	"export-nft-data/client/eth"
	"export-nft-data/contracts/erc721"
	"export-nft-data/contracts/seaport"
	"export-nft-data/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"math"
	"math/big"
	"strings"
)

const (
	defaultAddress = "0x00000000006c3852cbef3e08e8df289169ede581"
	ether          = "0x0000000000000000000000000000000000000000"
)

type BlockFilter struct {
	StartBlock uint64
	EndBlock   uint64
}

type OrderFilter struct {
	BlockFilter
	IgnoreTokens []string
}

type OwnerFilter struct {
	BlockFilter
	Token common.Address
}

type Stream interface {
	ForEachCollectionOrder(ctx context.Context, of *OrderFilter, handler func(o *CollectionOrder) error) error
	ForEachOwner(ctx context.Context, of *OwnerFilter, handler func(owner common.Address) error) error
}

type stream struct {
	f *seaport.SeaportFilterer
	e *ethclient.Client
}

func toMap(l []string) map[string]bool {
	res := make(map[string]bool)
	for _, item := range l {
		res[strings.ToLower(item)] = true
	}
	return res
}

var ErrNotEther = errors.New("not ether")

func calculatePrice(items []seaport.ReceivedItem) (*big.Int, error) {
	price := big.NewInt(0)
	for _, si := range items {
		// only consider ether payments
		if si.Token.Hex() != ether {
			return nil, ErrNotEther
		}
		price = price.Add(price, si.Amount)
	}
	return price, nil
}

func (s *stream) ForEachOwner(ctx context.Context, of *OwnerFilter, handler func(owner common.Address) error) error {
	ef, err := erc721.NewERC721Filterer(of.Token, s.e)
	if err != nil {
		return err
	}

	latest, err := s.e.BlockByNumber(ctx, nil)
	if err != nil {
		log.WithError(err).Error("error with BlockByNumber")
		return err
	}

	endBlock := uint64(math.Min(float64(latest.NumberU64()), float64(of.EndBlock)))

	return utils.WithPages(of.StartBlock, endBlock, 2000, func(start, end uint64) error {

		iterator, err := ef.FilterTransfer(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: ctx,
		}, nil, nil, nil)
		if err != nil {
			return err
		}

		for iterator.Next() {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err := handler(iterator.Event.To); err != nil {
				return err
			}
		}
		return nil
	})
}

type blockRange struct {
	start uint64
	end   uint64
}

func (s *stream) ForEachCollectionOrder(ctx context.Context, of *OrderFilter, handler func(o *CollectionOrder) error) error {
	tokenIgnoreList := toMap(of.IgnoreTokens)

	logger := log.WithFields(log.Fields{
		"method": "ForEachCollectionOrder",
	})
	pageCh := make(chan blockRange)
	grp, ctx := errgroup.WithContext(ctx)

	for i := 0; i < 10; i++ {
		grp.Go(func() error {
			f, err := seaport.NewSeaportFilterer(common.HexToAddress(defaultAddress), s.e)
			if err != nil {
				logger.WithError(err).Error("error creating filterer")
				return err
			}
			for p := range pageCh {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				logger := logger.WithFields(log.Fields{
					"start": p.start,
					"end":   p.end,
				})

				logger.Debug("page")

				iterator, err := f.FilterOrderFulfilled(&bind.FilterOpts{
					Start:   p.start,
					End:     &p.end,
					Context: ctx,
				}, nil, nil)
				if err != nil {
					logger.WithError(err).Error("FilterOrderFulfilled error")
					return err
				}

				for iterator.Next() {

					if ctx.Err() != nil {
						return ctx.Err()
					}
					o := iterator.Event

					// skip complicated orders (price isn't clear per item when bundled)
					if len(o.Offer) != 1 {
						continue
					}

					var token common.Address

					if _, ok := tokenIgnoreList[strings.ToLower(o.Offer[0].Token.Hex())]; ok {
						continue
					}
					token = o.Offer[0].Token

					price, err := calculatePrice(o.Consideration)
					if err == ErrNotEther {
						continue
					}
					if err != nil {
						logger.WithError(err).Error("error calculating price")
						return err
					}

					if err := handler(&CollectionOrder{
						Buyer:      o.Recipient,
						Seller:     o.Offerer,
						Price:      price,
						Collection: token,
						TxHash:     o.Raw.TxHash,
					}); err != nil {
						logger.WithError(err).Error("handler error")
						return err
					}
				}
			}
			return nil
		})
	}

	grp.Go(func() error {
		defer close(pageCh)

		latest, err := s.e.BlockByNumber(ctx, nil)
		if err != nil {
			log.WithError(err).Error("error with BlockByNumber")
			return err
		}

		endBlock := uint64(math.Min(float64(latest.NumberU64()), float64(of.EndBlock)))
		return utils.WithPages(of.StartBlock, endBlock, 2000, func(start, end uint64) error {
			select {
			case <-ctx.Done():
				logger.WithError(ctx.Err()).Error("returning from WithPages")
				return ctx.Err()
			default:
				pageCh <- blockRange{
					start: start,
					end:   end,
				}
			}
			return nil
		})
	})
	if err := grp.Wait(); err != nil {
		logger.WithError(err).Error("errgroup error")
		return err
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
		e: c,
	}, nil
}
