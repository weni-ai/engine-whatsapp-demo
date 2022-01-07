package utils

import (
	"net/http"
	"sync"
	"time"
)

var (
	once   sync.Once
	client *http.Client
)

func GetHTTPClient() *http.Client {
	once.Do(func() {
		client = &http.Client{
			Timeout: 60 * time.Second,
		}
	})

	return client
}

func CopyHeader(dest, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dest.Add(k, v)
		}
	}
}
