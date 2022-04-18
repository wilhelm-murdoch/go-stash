package queries

import (
	"encoding/json"
)

var (
	PostUnmarshaler = func(bytes []byte) (any, error) {
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
