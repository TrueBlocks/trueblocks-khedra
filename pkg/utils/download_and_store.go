package utils

import (
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return nil, err
	}

	_ = file.Touch(filename)
	return data, nil
}
