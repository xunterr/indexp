package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Info  os.FileInfo
	Title string
	Data  []byte
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
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrNotSupported
	}

	info, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	return &File{
		Info:  info,
		Title: getTitle(info),
		Data:  data,
	}, nil
}

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func GetSnippet(path string, lines []int) (string, error) {
	snippet := ""
	fd, err := os.Open(path)
	defer fd.Close()
	if err != nil {
		return "", err
	}
	for _, l := range lines {
		line, _, err := ReadLine(fd, l)
		if err != nil {
			return "", err
		}

		snippet = fmt.Sprintf("%s... %s", snippet, line)
	}
	return snippet, nil
}

func getTitle(fi os.FileInfo) string {
	fileName := fi.Name()
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
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
