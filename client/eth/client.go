package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"net"
	"net/http"
	"time"
)

func NewEthClient(url string) (*ethclient.Client, error) {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			Timeout:   5 * time.Second,
		}).DialContext,
		IdleConnTimeout:       5 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	reth, err := rpc.DialHTTPWithClient(url, &http.Client{Transport: tr})

	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(reth), nil
}
