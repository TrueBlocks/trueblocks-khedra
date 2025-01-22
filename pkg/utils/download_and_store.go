package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

func DownloadAndStore(url, filename string, dur time.Duration) ([]byte, error) {
	if file.FileExists(filename) {
		lastModDate, err := file.GetModTime(filename)
		if err != nil {
			return nil, err
		}
		if time.Since(lastModDate) < dur {
			data, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prettyData []byte
	if json.Valid(rawData) {
		var jsonData interface{}
		err := json.Unmarshal(rawData, &jsonData)
		if err != nil {
			return nil, err
		}
		prettyData, err = json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		// If the data is not valid JSON, write it as-is
		prettyData = rawData
	}

	err = os.WriteFile(filename, prettyData, 0644)
	if err != nil {
		return nil, err
	}

	_ = file.Touch(filename)
	return prettyData, nil
}
