package tools

import (
	"encoding/json"
	"os"
)

func SaveJson(path string, object any) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(object); err != nil {
		return err
	}

	return nil
}
