package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("directories should throw unsupported error", func(t *testing.T) {
		err := Copy("testdata/testDir", "out.txt", 0, 0)
		require.Error(t, err)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("directories should throw unsupported error", func(t *testing.T) {
		err := Copy("/dev/urandom", "out.txt", 0, 0)
		require.Error(t, err)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("offset more than file length should throw error", func(t *testing.T) {
		const fromPath = "testdata/simple_test.txt"
		const toPath = "text.txt"
		inFile, err := os.Open(fromPath)

		if err != nil {
			return
		}

		fileInfo, _ := inFile.Stat()
		fileSize := fileInfo.Size()

		err = Copy(fromPath, toPath, fileSize+1, 0)
		require.Error(t, err)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
		os.Remove(toPath)
	})

	t.Run("offset equal with file length should create empty file", func(t *testing.T) {
		const fromPath = "testdata/simple_test.txt"
		const toPath = "text.txt"
		inFile, err := os.Open(fromPath)

		if err != nil {
			return
		}

		fileInfo, _ := inFile.Stat()
		fileSize := fileInfo.Size()

		Copy(fromPath, toPath, fileSize+1, 0)

		outFile, _ := os.Open(toPath)
		outFileInfo, _ := outFile.Stat()

		require.Equal(t, outFileInfo.Size(), int64(0))
		os.Remove(toPath)
	})

	t.Run("offset property works", func(t *testing.T) {
		const fromPath = "testdata/simple_test.txt"
		const toPath = "text.txt"

		inFile, err := os.Open(fromPath)

		if err != nil {
			return
		}

		var offset int64 = 5

		fileInfo, _ := inFile.Stat()
		fileSize := fileInfo.Size()

		Copy(fromPath, toPath, offset, 0)

		outFile, _ := os.Open(toPath)
		outFileInfo, _ := outFile.Stat()

		require.Equal(t, outFileInfo.Size(), fileSize-offset)
		os.Remove(toPath)
	})

	t.Run("limit property works", func(t *testing.T) {
		const fromPath = "testdata/simple_test.txt"
		const toPath = "text.txt"

		inFile, err := os.Open(fromPath)

		if err != nil {
			return
		}

		var limit int64 = 5

		fileInfo, _ := inFile.Stat()
		fileSize := fileInfo.Size()

		Copy(fromPath, toPath, 0, limit)

		outFile, _ := os.Open(toPath)
		outFileInfo, _ := outFile.Stat()

		require.Equal(t, outFileInfo.Size(), fileSize-limit)
		os.Remove(toPath)
	})
}
