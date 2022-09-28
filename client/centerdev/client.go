package centerdev

import (
	"context"
	"export-nft-data/client/utils"
	"fmt"
	"strings"
)

const (
	pageLimit      = 100
	minRelevance   = 10
	apiHost        = "https://api.center.dev"
	apiCollections = "/v1/%s/search?query=%s&type=Collection"
)

type Client interface {
	GetCollections(ctx context.Context, network Network, name string) ([]*Collection, error)
}

type client struct {
	Headers map[string]string
}

func NewClient(key string) Client {
	return &client{Headers: map[string]string{
		"X-API-Key": key,
	}}
}

func (c *client) GetCollections(ctx context.Context, network Network, name string) ([]*Collection, error) {
	return getListAPI[Collection](ctx, fmt.Sprintf(apiCollections, network, name), c.Headers)
}

func getListAPI[T ResultData](ctx context.Context, urlPath string, headers map[string]string) ([]*T, error) {
	var result []*T
	var resp Result

	path := strings.ReplaceAll(fmt.Sprintf("%s%s", apiHost, urlPath), " ", "%20")
	if err := utils.Get(ctx, path, headers, &resp); err != nil {
		return nil, err
	}
	for _, item := range resp.Results {
		b, err := item.MarshalJSON()
		if err != nil {
			return nil, err
		}
		res, err := unmarshalAny[T](b)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if len(result) == 0 {
		return make([]*T, 0), nil
	}
	return result, nil
}

// getPageableList iterates over all pages of the API
func getPageableList[T ResultData](ctx context.Context, urlPath string, headers map[string]string) ([]*T, error) {
	var result []*T
	var offset int64
	for {
		sym := "?"
		if strings.Contains(urlPath, "?") {
			sym = "&"
		}
		path := fmt.Sprintf("%s%soffset=%d&limit=%d", urlPath, sym, offset, pageLimit)
		items, err := getListAPI[T](ctx, path, headers)
		if err != nil {
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		result = append(result, items...)
		offset += int64(len(items))
	}
	if len(result) == 0 {
		return make([]*T, 0), nil
	}
	return result, nil
}
