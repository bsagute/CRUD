package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/hashicorp/go-retryablehttp"
    "golang.org/x/oauth2"
)

// Client represents an HTTP client with a timeout
type Client struct {
    HTTPClient *http.Client
    timeout    time.Duration
}

// ClientOption defines a function type for setting client options
type ClientOption func(*Client)

// defaultClient returns a standard HTTP client with the given timeout
func defaultClient(timeout time.Duration) *http.Client {
    return &http.Client{
        Timeout: timeout,
    }
}

// passthroughErrorHandler is a basic error handler for retryablehttp
func passthroughErrorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
    // Log each retry attempt
    fmt.Printf("Attempt %d: Error %v\n", numTries, err)
    return resp, err
}

// initHTTP initializes a retryable HTTP client with the given retry max and logger
func initHTTP(retryMax int, leveledLogger retryablehttp.LeveledLogger) ClientOption {
    return func(c *Client) {
        // Create a new retryable HTTP client
        retryableHTTPClient := retryablehttp.NewClient()
        retryableHTTPClient.HTTPClient = defaultClient(c.timeout) // Initialize with default client

        // Set retry max and logger
        retryableHTTPClient.RetryMax = retryMax
        retryableHTTPClient.Logger = leveledLogger

        // Set the error handler
        retryableHTTPClient.ErrorHandler = passthroughErrorHandler
        fmt.Printf("== HTTPCLIENT passthroughErrorHandler.retryMax: %d\n", retryMax)

        // Assign the retryable client to the main client
        c.HTTPClient = retryableHTTPClient.StandardClient()
        c.HTTPClient.Timeout = c.timeout
    }
}

func main() {
    // Initialize the main client with a timeout
    client := &Client{
        timeout: 10 * time.Second,
    }

    // Apply the retryable HTTP client option
    option := initHTTP(3, retryablehttp.NewLogger(retryablehttp.INFO))
    option(client)

    // Create a new request
    req, err := retryablehttp.NewRequest("GET", "https://example.com", nil)
    if err != nil {
        fmt.Printf("Error creating request: %v\n", err)
        return
    }

    // Perform the request using the retryable HTTP client
    resp, err := client.HTTPClient.Do(req)
    if err != nil {
        fmt.Printf("Request failed: %v\n", err)
        return
    }
    defer resp.Body.Close()

    // Print the status of the response
    fmt.Println("Request succeeded with status:", resp.Status)
}
