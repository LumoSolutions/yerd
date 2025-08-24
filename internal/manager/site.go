package manager

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

type SiteManager struct {
	Directory    string
	PhpVersion   string
	Domain       string
	Spinner      *utils.Spinner
	WebConfig    *config.WebConfig
	PublicFolder string
}

func NewSiteManager() (*SiteManager, error) {
	var webConfig *config.WebConfig
	if err := config.GetStruct("web", webConfig); err != nil {
		webConfig = &config.WebConfig{}
	}

	s := utils.NewSpinner("Managing Sites...")
	s.SetDelay(150)

	return &SiteManager{
		Spinner:   s,
		WebConfig: webConfig,
	}, nil
}

func (siteManager *SiteManager) AddSite(directory, domain, publicFolder, phpVersion string) error {
	siteManager.Directory = directory
	siteManager.Domain = domain
	siteManager.PhpVersion = phpVersion
	siteManager.PublicFolder = publicFolder

	siteManager.Spinner.Start()

	err := utils.RunAll(
		func() error { return siteManager.validateDirectory() },
		func() error { return siteManager.validateDomain() },
		func() error { return siteManager.validatePhpVersion() },
		func() error { return siteManager.createSiteConfig() },
		//func() error { return siteManager.createHostsEntry() },
		func() error { return siteManager.restartNginx() },
	)

	if err != nil {
		siteManager.Spinner.StopWithError("Failed to add site")
		return err
	}

	siteManager.Spinner.StopWithSuccess("Finished")

	return nil
}

func (siteManager *SiteManager) validateDirectory() error {
	siteManager.Spinner.UpdatePhrase("Validating Path...")
	abs, err := filepath.Abs(siteManager.Directory)
	if err != nil {
		siteManager.Spinner.AddErrorStatus("Directory provided is invalid")
		return err
	}

	if !utils.IsDirectory(abs) {
		siteManager.Spinner.AddErrorStatus("Path provided is not a directory")
		return fmt.Errorf("path not a directory")
	}

	siteManager.Directory = abs
	siteManager.Spinner.AddInfoStatus("Directory: %s", abs)

	if siteManager.PublicFolder == "" {
		if utils.IsDirectory(filepath.Join(abs, "public")) {
			siteManager.PublicFolder = "public"
			siteManager.Spinner.AddInfoStatus("Public directory discovered")
			siteManager.Spinner.AddInfoStatus("Serving from: /public")
		}

		return nil
	}

	if !utils.IsDirectory(filepath.Join(abs, siteManager.PublicFolder)) {
		siteManager.Spinner.AddErrorStatus("Provided public path is not a directory")
		return fmt.Errorf("public folder does not exist")
	}

	siteManager.Spinner.AddInfoStatus("Serving from /%s", siteManager.PublicFolder)

	return nil
}

func (siteManager *SiteManager) validateDomain() error {
	if siteManager.Domain == "" {
		base := filepath.Base(siteManager.Directory)
		base = strings.ToLower(base) + ".test"
		siteManager.Domain = base
		siteManager.Spinner.AddInfoStatus("Domain: %s", base)
	}

	if siteManager.WebConfig.Sites != nil {
		for _, site := range siteManager.WebConfig.Sites {
			if site.Domain == siteManager.Domain {
				siteManager.Spinner.AddErrorStatus("Domain is already in use")
				siteManager.Spinner.AddInfoStatus("- Use the -d flag to specify a custom domain")
				return fmt.Errorf("domain in use")
			}
		}
	}

	return nil
}

func (siteManager *SiteManager) validatePhpVersion() error {
	siteManager.Spinner.UpdatePhrase("Validating PHP Version...")

	if siteManager.PhpVersion == "" {
		versions := constants.GetAvailablePhpVersions()
		slices.Reverse(versions)

		for _, version := range versions {
			if _, installed := config.GetInstalledPhpInfo(version); installed {
				siteManager.PhpVersion = version
				siteManager.Spinner.AddInfoStatus("Using PHP %s", version)
				break
			}
		}
	}

	_, installed := config.GetInstalledPhpInfo(siteManager.PhpVersion)
	if !installed {
		siteManager.Spinner.AddErrorStatus("PHP%s is not installed, or is not valid", siteManager.PhpVersion)
		return fmt.Errorf("php version not installed")
	}

	return nil
}

func (siteManager *SiteManager) createSiteConfig() error {
	siteManager.Spinner.UpdatePhrase("Downloading site.conf...")
	content, err := utils.FetchFromGitHub("nginx", "site.conf")
	if err != nil {
		siteManager.Spinner.AddErrorStatus("Unable to download site.conf")
	}

	projectPath := filepath.Join(siteManager.Directory, siteManager.PublicFolder)
	content = utils.Template(content, utils.TemplateData{
		"domain":      siteManager.Domain,
		"path":        projectPath,
		"php_version": siteManager.PhpVersion,
	})

	path := filepath.Join(constants.YerdWebDir, "nginx", "sites-enabled", siteManager.Domain+".conf")
	err = utils.WriteStringToFile(
		path,
		content,
		constants.FilePermissions,
	)

	if err != nil {
		siteManager.Spinner.AddErrorStatus("Unable to save %s", filepath.Base(path))
		return err
	}

	siteManager.Spinner.AddSuccessStatus("Created Nginx Configration (%s)", filepath.Base(path))

	return nil
}

func (siteManager *SiteManager) createHostsEntry() error {
	hostManager := utils.NewHostsManager()
	if err := hostManager.Add(siteManager.Domain); err != nil {
		siteManager.Spinner.AddErrorStatus("Unable to add hosts entry")
		utils.LogError(err, "hosts")
		return err
	}

	return nil
}

func (siteManager *SiteManager) restartNginx() error {
	siteManager.Spinner.UpdatePhrase("Restarting Nginx...")
	utils.SystemdStopService("yerd-nginx")
	if err := utils.SystemdStartService("yerd-nginx"); err != nil {
		utils.LogError(err, "nginx")
		siteManager.Spinner.AddErrorStatus("Failed to restart nginx")
		return err
	}

	return nil
}
