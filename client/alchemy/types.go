package alchemy

import "github.com/ethereum/go-ethereum/common"

type OwnerResult struct {
	OwnerAddresses []common.Address `json:"ownerAddresses"`
	PageKey        string           `json:"pageKey"`
}
