package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/segfaultax/go-nagios"
	"github.com/spf13/pflag"
)

func TestMain(t *testing.T) {
	// Save the original command-line arguments
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set up the test command-line arguments
	os.Args = []string{"go-check-prometheus", "-H", "localhost", "-q", "test_query", "-w", "10", "-c", "100"}

	// Set up the test flags
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// Set up the test output buffer
	output := new(bytes.Buffer)
	printUsageErrorAndExit = func(code int, err error) {
		fmt.Fprintf(output, "execution failed: %s\n", err)
		printUsage()
		os.Exit(code)
	}

	// Set up the test Prometheus server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/query" {
			fmt.Fprint(w, `{"status": "success", "data": {"resultType": "vector", "result": []}}`)
		}
	}))
	defer server.Close()

	// Set up the test Prometheus client
	client, _ := api.NewClient(api.Config{
		Address: server.URL,
		RoundTripper: (&http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 30 * time.Second,
		}),
	})
	v1api := v1.NewAPI(client)

	// Set up the test context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set up the test check
	check := nagios.NewRangeCheck()

	// Set up the test result
	result := &v1.Result{
		Type: v1.ValVector,
	}

	// Set up the test warnings
	warnings := []string{"warning"}

	// Override the functions used in the main function
	checkRequiredOptions = func() error {
		return nil
	}
	v1api.Query = func(ctx context.Context, query string, ts time.Time) (v1.Value, []string, error) {
		return result, warnings, nil
	}
	runCheck = func(c *nagios.RangeCheck, result v1.Value) {
		c.Status = nagios.StatusOK
		c.SetMessage("Test message")
	}

	// Run the main function
	main()

	// Check the test output
	expectedOutput := "The query did not return any result\n"
	if output.String() != expectedOutput {
		t.Errorf("Expected output: %s, but got: %s", expectedOutput, output.String())
	}
}
