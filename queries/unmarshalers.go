package queries

import (
	"encoding/json"
	"fmt"
)

type PostError struct {
	Errors []map[string]any
}

var (
	PostUnmarshaler = func(bytes []byte) (any, error) {
		var e PostError
		if err := json.Unmarshal(bytes, &e); err != nil {
			return nil, err
		}

		if len(e.Errors) > 0 {
			return nil, fmt.Errorf("unexpected API response: %s", e.Errors[0]["message"])
		}

		// { "data": { "post": { ... } } }
		var results map[string]map[string]Post
		if err := json.Unmarshal(bytes, &results); err != nil {
			return nil, err
		}

		return results["data"]["post"], nil
	}

	TimelineUnmarshaler = func(bytes []byte) (any, error) {
		// { "data": { "user": { "publication": { "posts": [ ... ] } } } }
		var results map[string]map[string]map[string]map[string][]Post
		if err := json.Unmarshal(bytes, &results); err != nil {
			return nil, err
		}

		return results["data"]["user"]["publication"]["posts"], nil
	}
)
