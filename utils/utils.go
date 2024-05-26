package utils

import (
	"bufio"
	"encoding/json"
	"os"
)

func SaveJSON[T any](filename string, data T) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	_, err = w.Write(json)
	return err
}

func LoadJSON[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result T
	err = json.Unmarshal(data, &result)
	return &result, err
}
