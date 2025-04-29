package bpebuilder

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/unixpickle/essentials"
)

// LoadJSONData reads all of the JSON files inside a directory and parses each
// one as a []string. It concatenates the results.
func LoadJSONData(pathToDir string) ([]string, error) {
	entries, err := os.ReadDir(pathToDir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		fullPath := filepath.Join(pathToDir, entry.Name())
		paths = append(paths, fullPath)
	}

	var results []string
	var decodeErr error
	essentials.ReduceConcurrentMap(0, len(paths), func() (func(int), func()) {
		var localResults []string
		var localErr error
		return func(i int) {
				data, err := os.ReadFile(paths[i])
				if err != nil {
					localErr = err
					return
				}
				if err = json.Unmarshal(data, &localResults); err != nil {
					localErr = err
				}
			}, func() {
				if localErr != nil && decodeErr == nil {
					decodeErr = localErr
				}
				results = append(results, localResults...)
			}
	})

	return results, decodeErr
}
