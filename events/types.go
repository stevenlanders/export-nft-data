package events

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type CollectionOrder struct {
	Buyer      common.Address `json:"buyer"`
	Seller     common.Address `json:"seller"`
	Price      *big.Int       `json:"price"`
	Collection common.Address `json:"collection"`
	TxHash     common.Hash    `json:"tx"`
}
