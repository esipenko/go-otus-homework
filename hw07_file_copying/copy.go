package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrPathsAreEqual         = errors.New("paths are equal")
)

func arePathEqual(fromPath, toPath string) (bool, error) {
	fromAbs, err := filepath.Abs(fromPath)
	if err != nil {
		return false, err
	}
	toAbs, err := filepath.Abs(toPath)
	if err != nil {
		return false, err
	}

	return fromAbs == toAbs, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	areEqual, err := arePathEqual(fromPath, toPath)
	if err != nil {
		return err
	}

	if areEqual {
		return ErrPathsAreEqual
	}

	inFile, err := os.Open(fromPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer inFile.Close()

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer outFile.Close()

	fileInfo, err := inFile.Stat()
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()

	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		_, err = inFile.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	bytesToCopy := fileSize - offset

	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	bar := pb.Default.Start64(bytesToCopy)

	defer bar.Finish()

	barReader := bar.NewProxyReader(inFile)

	_, err = io.CopyN(outFile, barReader, bytesToCopy)
	if err != nil {
		removeErr := os.Remove(to)
		if removeErr != nil {
			return fmt.Errorf("copy error: %w, remove error: %w", err, removeErr)
		}

		return err
	}

	return nil
}
