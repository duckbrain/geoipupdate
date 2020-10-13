// Package internal is none of your business
package internal

import (
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

// MaybeRetryRequest is an internal implementation detail of this module. It
// shouldn't be used by users of the geoipupdate Go library. You can use the
// RetryFor field of geoipupdate.Config if you'd like to retry failed requests
// when using the library directly.
func MaybeRetryRequest(c *http.Client, retryFor time.Duration, req *http.Request) (*http.Response, error) {
	if retryFor < 0 {
		return nil, errors.New("negative retry duration")
	}
	if req.Body != nil {
		return nil, errors.New("can't retry requests with bodies")
	}
	exp := backoff.NewExponentialBackOff()
	exp.MaxElapsedTime = retryFor
	var resp *http.Response
	err := backoff.Retry(
		func() error {
			var err error
			resp, err = c.Do(req) // nolint: bodyclose
			if err != nil {
				return errors.Wrap(err, "error performing http request")
			}
			return err
		},
		exp,
	)
	return resp, err
}
