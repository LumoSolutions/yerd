package composer

import (
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

func downloadComposer() error {
	return utils.DownloadFile(
		constants.ComposerDownloadUrl,
		constants.LocalComposerPath,
		nil,
	)
}

func InstallComposer() error {
	return utils.RunAll(
		func() error { return downloadComposer() },
		func() error { return utils.Chmod(constants.LocalComposerPath, 0755) },
		func() error { return utils.CreateSymlink(constants.LocalComposerPath, constants.GlobalComposerPath) },
	)
}

func RemoveComposer() error {
	return utils.RunAll(
		func() error { return utils.RemoveSymlink(constants.GlobalComposerPath) },
		func() error { return utils.RemoveFile(constants.LocalComposerPath) },
	)
}
