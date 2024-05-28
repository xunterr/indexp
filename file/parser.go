package file

import (
	"bufio"
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
		data, err = parseSimple(fd)
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

func parseSimple(fd *os.File) ([]byte, error) {
	buff := bufio.NewReader(fd)
	data, err := io.ReadAll(buff)
	if err != nil {
		log.Fatalf("Error reading file: %s", err.Error())
		return nil, err
	}
	return data, nil
}
