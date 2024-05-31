package file

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"os"
)

func SaveJSON[T any](filename string, data T) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rewrite(filename, json)
}

func SaveGOB[T any](filename string, data T) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	buf := bufio.NewWriter(f)

	encoder := gob.NewEncoder(buf)
	return encoder.Encode(data)

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

func LoadGOB[T any](filename string, to *T) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(f)

	decoder := gob.NewDecoder(buf)
	return decoder.Decode(to)
}
