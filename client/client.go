package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wilhelm-murdoch/go-stash/queries"
)

// Client
type Client struct {
	client  http.Client
	baseUrl string
}

// NewClient
func New() *Client {
	c := http.Client{
		Timeout: 60 * time.Second,
	}

	return &Client{
		client:  c,
		baseUrl: "https://api.hashnode.com",
	}
}

// Execute
func (c *Client) Execute(query queries.Query) (any, error) {
	response, err := c.client.Post(c.baseUrl, "application/json", strings.NewReader(query.String()))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 429 {
		if query.Backoff.Attempt() >= 10 {
			return nil, fmt.Errorf("maximum attempts reached; skipping %s", query.Name)
		}
		d := query.Backoff.Duration()
		log.Printf("rate limit detected for %s; retrying in %s for attempt %d/%d\n", query.Name, d, int(query.Backoff.Attempt()), 10)
		time.Sleep(d)
		return c.Execute(query)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected a 200 response, but got %d instead", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return query.Unmarshaler(body)
}
