package export

import (
	"fmt"
	"os"
)

func DebugFile(filename string, content []byte) error {
	// Create (or truncate) the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
