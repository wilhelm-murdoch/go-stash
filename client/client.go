package client

import (
	"fmt"
	"io/ioutil"
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
func (c *Client) Execute(query *queries.Query) (any, error) {
	response, err := c.client.Post(c.baseUrl, "application/json", strings.NewReader(query.String()))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected a 200 response, but got %d instead", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return query.Unmarshaler(body)
}
