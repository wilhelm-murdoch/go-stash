package queries

import "encoding/json"

var (
	PostsUnmarshaler = func(bytes []byte) (any, error) {
		// { "data": { "user": { "publication": { "posts": [ ... ] } } } }
		var results map[string]map[string]map[string]map[string][]Post
		if err := json.Unmarshal(bytes, &results); err != nil {
			return nil, err
		}

		return results["data"]["user"]["publication"]["posts"], nil
	}

	AuthorUnmarshaler = func(bytes []byte) (any, error) {
		// { "data": { "user": { ... } } }
		var results map[string]map[string]Author
		if err := json.Unmarshal(bytes, &results); err != nil {
			return nil, err
		}

		return results["data"]["user"], nil
	}
)
