package utils

import (
	"errors"
	"fmt"
	"os"
)

func CreateWriteFile(filePath string, content string) error {
	newF, err := os.Create(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("creating %s file error: %w", filePath, err)
		}
	}
	defer newF.Close()

	if _, err = newF.WriteString(content); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("writing %s file error: %w", filePath, err)
		}
	}
	return nil
}
