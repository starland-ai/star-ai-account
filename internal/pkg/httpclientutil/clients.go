package httpclientutil

import (
	"net/http"
	"sync"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
	"go.uber.org/zap"
)

var (
	httpClientPool = sync.Pool{
		New: func() interface{} {
			return httpclient.NewClient(
				httpclient.WithHTTPClient(MustStdClient()),
				//nolint:gomnd
				httpclient.WithHTTPTimeout(5*time.Second),
				httpclient.WithRetrier(heimdall.NewRetrier(
					//nolint:gomnd
					heimdall.NewConstantBackoff(1*time.Second, 500*time.Millisecond),
				)),
				httpclient.WithRetryCount(2),
			)
		},
	}
)

func NewStdClient() (*http.Client, error) {
	var (
		transport = &http.Transport{}
	)

	// read proxy from env
	transport.Proxy = http.ProxyFromEnvironment

	return &http.Client{Transport: transport}, nil
}

func MustStdClient() *http.Client {
	client, err := NewStdClient()
	if err != nil {
		zap.S().Error(err)
	}
	return client
}

func GetHttpClient() *httpclient.Client {
	return httpClientPool.Get().(*httpclient.Client)
}
