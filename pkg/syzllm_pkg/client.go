package syzllm_pkg

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

func SendCoverAsync(cover uint64) {
	// Convert cover to bytes
	numberBytes := []byte(strconv.Itoa(int(cover)))

	// Get server info
	serverHostInDocker := os.Getenv("SERVER_HOST")
	serverPortInDocker := os.Getenv("SERVER_PORT")
	serverInfo, err := GetServerInfo("test")
	if err != nil {
		fmt.Printf("Error getting server info: %v\n", err)
		return
	}
	if serverHostInDocker != "" && serverPortInDocker != "" {
		serverInfo.Host = serverHostInDocker
		serverInfo.Port = serverPortInDocker
	}
	url := fmt.Sprintf("http://%s:%s/cover", serverInfo.Host, serverInfo.Port)

	// Create HTTP client and request
	client := getClient().client
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(numberBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	// Send async POST request
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending cover to server: %v\n", err)
			return
		}
		resp.Body.Close()
	}()
}

func ExtractCoverage(log string) (int, error) {
	// Regular expression to match "coverage=<number>"
	re := regexp.MustCompile(`coverage=(\d+)`)
	matches := re.FindStringSubmatch(log)

	if len(matches) < 2 {
		return 0, fmt.Errorf("coverage value not found in log")
	}

	// Convert the captured number to an integer
	coverage, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse coverage value: %v", err)
	}

	return coverage, nil
}
