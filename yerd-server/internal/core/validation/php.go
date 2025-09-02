package validation

import (
	"slices"

	"github.com/lumosolutions/yerd/server/internal/config"
	"github.com/lumosolutions/yerd/server/internal/constants"
)

func IsPhpVersionValid(version string) bool {
	if slices.Contains(constants.PhpVersions, version) {
		return true
	}

	return false
}

func IsPhpInstalled(version string) (*config.PhpInfo, bool) {
	appConfig, _ := config.GetConfig()
	info, exists := appConfig.Php[version]
	return info, exists
}
