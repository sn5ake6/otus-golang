package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var file *os.File
	var err error

	file, err = os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer file.Close()

	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	countToRead := fileSize - offset
	if limit > 0 && (offset+limit) < fileSize {
		countToRead = limit
	}

	if offset > 0 {
		file.Seek(offset, io.SeekStart)
	}

	reader := io.LimitReader(file, countToRead)

	bar := pb.Full.Start64(countToRead)
	barReader := bar.NewProxyReader(reader)

	var fileTo *os.File
	fileTo, err = os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	_, err = io.Copy(fileTo, barReader)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}
