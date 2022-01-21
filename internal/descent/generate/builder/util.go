package builder

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

func LoadBuilder(path string) ([]*CaveBuilder, error) {
	errMsg := "load image"
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var caveBuilder []*CaveBuilder
	err = decoder.Decode(&caveBuilder)
	return caveBuilder, nil
}
