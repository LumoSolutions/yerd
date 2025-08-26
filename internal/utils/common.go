package utils

import (
	"fmt"

	"github.com/lumosolutions/yerd/internal/constants"
)

const context = "utils"

func EnsureYerdDirectories() error {
	dirs := []string{
		constants.YerdBaseDir,
		constants.YerdBinDir,
		constants.YerdPHPDir,
		constants.YerdEtcDir,
		constants.YerdWebDir,
	}

	for _, dir := range dirs {
		if err := CreateDirectory(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

func RunAll(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
