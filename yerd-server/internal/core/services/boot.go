package services

import (
	"fmt"
	"path/filepath"

	"github.com/lumosolutions/yerd/server/internal/constants"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

func Boot() error {
	return utils.Run(
		func() error { return createDirectories() },
	)
}

func createDirectories() error {
	elevated := utils.IsRunningElevated()

	workingPath := constants.WorkingPath
	if elevated {
		workingPath = constants.WorkingElevatedPath
	}

	workingPath = constants.ExpandPath(workingPath)
	reqDirectories := []string{
		workingPath,
		filepath.Join(workingPath, constants.LogPathRelative),
		filepath.Join(workingPath, constants.EtcPathRelative),
		filepath.Join(workingPath, constants.PhpPathRelative),
		filepath.Join(workingPath, constants.NginxPathRelative),
		filepath.Join(workingPath, constants.SrcPathRelative),
	}

	for _, directory := range reqDirectories {
		if utils.IsDirectory(directory) {
			continue
		}

		fmt.Printf("Creating directory: %s\n", directory)
		err := utils.CreateDirectory(directory)
		if err != nil {
			fmt.Printf("Error creating directory: %s\n", directory)
			return err
		}
	}

	return nil
}
