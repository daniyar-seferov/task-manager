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
	Date   string `json:"date"`
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

	t, err := time.Parse(time.RFC3339, ver.Date)
	if err != nil {
		return fmt.Errorf("parse date: %w", err)
	}

	fmt.Println("api-server:", ver.Commit, timeAgo(t))
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

func timeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
	}
}
