package centerdev

import "encoding/json"

type Network string

const NetworkEthereumMainnet Network = "ethereum-mainnet"

type Collection struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	PreviewImageUrl string  `json:"previewImageUrl"`
	Relevance       float64 `json:"relevance"`
	Url             string  `json:"url"`
	Address         string  `json:"address"`
	Type            string  `json:"type"`
}

type ResultData interface {
	Collection
}

type Result struct {
	Results []*json.RawMessage `json:"results"`
}

func unmarshalAny[T any](bytes []byte) (*T, error) {
	out := new(T)
	if err := json.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}
