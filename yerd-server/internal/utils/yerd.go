package utils

import (
	"os"
	"path/filepath"

	"github.com/lumosolutions/yerd/server/internal/constants"
)

func IsRunningElevated() bool {
	return os.Getuid() == 0
}

func GetWorkingPath() string {
	if IsRunningElevated() {
		return constants.ExpandPath(constants.WorkingElevatedPath)
	}

	return constants.ExpandPath(constants.WorkingPath)
}

func GetGlobalBinPath() string {
	if IsRunningElevated() {
		return constants.ExpandPath(constants.BinElevatedPath)
	}
	return constants.ExpandPath(constants.BinPath)
}

func GetYerdLogPath() string {
	return filepath.Join(GetWorkingPath(), constants.LogPathRelative)
}

func GetYerdPhpPath() string {
	return filepath.Join(GetWorkingPath(), constants.PhpPathRelative)
}

func GetYerdEtcPath() string {
	return filepath.Join(GetWorkingPath(), constants.EtcPathRelative)
}

func GetYerdBinPath() string {
	return filepath.Join(GetWorkingPath(), constants.BinPathRelative)
}
