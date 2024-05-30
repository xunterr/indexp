package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Info os.FileInfo
	Data []byte
}

var ErrNotSupported = errors.New("File type is not supported")

func ReadFile(path string) (*File, error) {
	fd, err := os.Open(path)
	defer fd.Close()
	if err != nil {
		return nil, err
	}

	extension := filepath.Ext(path)
	var data []byte
	switch extension {
	case ".txt", ".md", ".xml", ".csv", ".json":
		data, err = bufferedRead(fd)
	default:
		return nil, ErrNotSupported
	}

	if err != nil {
		return nil, err
	}

	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	return &File{
		Info: info,
		Data: data,
	}, nil
}

func bufferedRead(fd *os.File) ([]byte, error) {
	buff := bufio.NewReader(fd)
	data, err := io.ReadAll(buff)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func rewrite(filename string, data []byte) error {
	file, err := os.Create(filename)
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Fatalln(closeErr.Error())
			return
		}
	}()
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	_, err = w.Write(data)
	return err
}

func SaveJSON[T any](filename string, data T) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rewrite(filename, json)
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
