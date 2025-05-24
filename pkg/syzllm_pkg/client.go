package syzllm_pkg

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

type client struct {
	client *http.Client
	once   sync.Once
}

var instance *client

func getClient() *client {
	if instance == nil {
		instance = &client{}
		instance.once.Do(func() {
			instance.client = &http.Client{}
		})
	}
	return instance
}

func SendPostRequestAsync(url string, jsonData []byte) (<-chan *http.Response, <-chan error) {
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("error creating request: %w", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := getClient().client.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("error executing request: %w", err)
			return
		}

		responseChan <- resp
	}()

	return responseChan, errorChan
}
