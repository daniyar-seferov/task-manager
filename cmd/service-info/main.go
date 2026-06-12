package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("too few arguments provided")
	}
	cmd := os.Args[1]

	switch cmd {
	case "versions":
		return runVersions()
	case "is-deployed":
		return runIsDeployed()
	}

	return nil
}

type VersionResponse struct {
	Commit string `json:"commit"`
}

func runVersions() error {
	url := "http://localhost:8090/debug/info"
	resp, err := sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("get debug info: %w", err)
	}

	var ver VersionResponse
	if err := json.Unmarshal(resp, &ver); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	fmt.Println("api-server:", ver.Commit)
	return nil
}

func runIsDeployed() error {
	return nil
}

func sendRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return respBody, nil
}
