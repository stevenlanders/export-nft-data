package etherscan

type ContractCreation struct {
	ContractAddress string `json:"contractAddress"`
	ContractCreator string `json:"contractCreator"`
	TxHash          string `json:"txHash"`
}

type ContractCreationResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Result  []ContractCreation `json:"result"`
}
