package domain

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Collection struct {
	Name        string            `json:"name"`
	Address     common.Address    `json:"address"`
	DeployBlock *big.Int          `json:"deployBlock,omitempty"`
	OwnerBlock  *big.Int          `json:"ownerBlock,omitempty"`
	Owners      []common.Address  `json:"owners,omitempty"`
	Edges       []*CollectionEdge `json:"edges"`
	Processed   bool              `json:"processed"`
}

type CollectionEdge struct {
	FromCollection *common.Address `json:"fromCollection"`
	ToCollection   *common.Address `json:"toCollection"`
	Buyer          *common.Address `json:"buyer"`
	Price          *big.Int        `json:"price"`
}
