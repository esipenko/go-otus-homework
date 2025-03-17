package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inFile, err := os.Open(fromPath)
	defer inFile.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	outFile, err := os.Create(toPath)
	defer outFile.Close()

	if err != nil {
		return err
	}

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
		return err
	}

	return nil

}
