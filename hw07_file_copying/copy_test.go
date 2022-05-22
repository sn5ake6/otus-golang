package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	fromPath = "testdata/input.txt"
	tempPath = os.TempDir() + "/hw07_tmp.out"
)

func TestCopy(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		var file *os.File
		file, _ = os.Open(fromPath)
		defer file.Close()

		var fileInfo os.FileInfo
		fileInfo, _ = file.Stat()

		fileSize := fileInfo.Size()
		fileSize++

		err := Copy(fromPath, "", fileSize, 0)

		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("", "", 0, 0)

		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("copy full", func(t *testing.T) {
		var file *os.File
		file, _ = os.Open(fromPath)
		defer file.Close()

		var fileInfo os.FileInfo
		fileInfo, _ = file.Stat()

		defer os.Remove(tempPath)

		err := Copy(fromPath, tempPath, 0, 0)

		var copiedFile *os.File
		copiedFile, _ = os.Open(tempPath)
		defer copiedFile.Close()

		var copiedFileInfo os.FileInfo
		copiedFileInfo, _ = copiedFile.Stat()

		require.Equal(t, nil, err)
		require.Equal(t, fileInfo.Size(), copiedFileInfo.Size())
	})

	t.Run("copy partial", func(t *testing.T) {
		var file *os.File
		file, _ = os.Open(fromPath)
		defer file.Close()

		var fileInfo os.FileInfo
		fileInfo, _ = file.Stat()

		fileSize := fileInfo.Size()
		countToRead := fileSize / 4
		offset := countToRead

		defer os.Remove(tempPath)

		err := Copy(fromPath, tempPath, offset, countToRead)

		var copiedFile *os.File
		copiedFile, _ = os.Open(tempPath)
		defer copiedFile.Close()

		var copiedFileInfo os.FileInfo
		copiedFileInfo, _ = copiedFile.Stat()

		require.Equal(t, nil, err)
		require.Equal(t, countToRead, copiedFileInfo.Size())
	})

	t.Run("limit exceeds file size", func(t *testing.T) {
		var file *os.File
		file, _ = os.Open(fromPath)
		defer file.Close()

		var fileInfo os.FileInfo
		fileInfo, _ = file.Stat()
		fileSize := fileInfo.Size()

		defer os.Remove(tempPath)

		err := Copy(fromPath, tempPath, 0, (fileSize * 2))

		var copiedFile *os.File
		copiedFile, _ = os.Open(tempPath)
		defer copiedFile.Close()

		var copiedFileInfo os.FileInfo
		copiedFileInfo, _ = copiedFile.Stat()

		require.Equal(t, nil, err)
		require.Equal(t, fileSize, copiedFileInfo.Size())
	})
}
