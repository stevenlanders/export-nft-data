package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var ErrRateLimit = errors.New("rate limit")

func WithBackoff(
	ctx context.Context, operation func(ctx context.Context) error,
) error {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 5 * time.Second
	bo.MaxInterval = 30 * time.Second

	err := backoff.Retry(func() error {
		if ctx.Err() != nil {
			return backoff.Permanent(ctx.Err())
		}
		err := operation(ctx)
		if err == nil {
			return nil
		} else if ctx.Err() != nil {
			return backoff.Permanent(ctx.Err())
		} else if err == ErrRateLimit {
			log.Warnf("%s (retryable)", err.Error())
			return err
		}
		return backoff.Permanent(err)
	}, bo)
	if err != nil {
		log.Errorf("failed with error: %s", err.Error())
	}
	return err
}

func Get(ctx context.Context, path string, headers map[string]string, target interface{}) error {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	for h, v := range headers {
		req.Header.Add(h, v)
	}

	return WithBackoff(ctx, func(ctx context.Context) error {
		req = req.WithContext(ctx)

		dc := http.DefaultClient
		resp, err := dc.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || resp.StatusCode == 502 {
			return ErrRateLimit
		}
		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("error code %d, body=%s, path=%s", resp.StatusCode, string(b), path)
		}
		return json.NewDecoder(resp.Body).Decode(target)
	})
}
