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
	CrtFile      string
	KeyFile      string
}

func NewSiteManager() (*SiteManager, error) {
	s := utils.NewSpinner("Managing Sites...")
	s.SetDelay(150)

	var webConfig *config.WebConfig
	if err := config.GetStruct("web", &webConfig); err != nil {
		webConfig = &config.WebConfig{}
	}

	if !webConfig.Installed {
		return nil, fmt.Errorf("web not installed")
	}

	return &SiteManager{
		Spinner:   s,
		WebConfig: webConfig,
	}, nil
}

func (sm *SiteManager) ListSites() {
	if len(sm.WebConfig.Sites) == 0 {
		sm.Spinner.AddWarningStatus("No Sites Created")
		sm.Spinner.AddInfoStatus("Create a site with one of the following commands:")
		sm.Spinner.AddInfoStatus("'sudo yerd sites add .'  # use current directory")
		sm.Spinner.AddInfoStatus("'sudo yerd sites add relative/folder'")
		sm.Spinner.AddInfoStatus("'sudo yerd sites add /home/user/absolute/folder'")
		return
	}

	for _, site := range sm.WebConfig.Sites {
		fmt.Printf("🌐 Site: %s  (PHP %s)\n", site.Domain, site.PhpVersion)
		fmt.Printf("├─ Secure Link: https://%s/\n", site.Domain)
		fmt.Printf("└─ Directory: %s\n\n", site.RootDirectory)
	}
}

func (sm *SiteManager) SetValue(name, value, site string) error {
	sm.Spinner.UpdatePhrase("Updating Site...")
	sm.Spinner.Start()

	if !sm.identifySite(site) {
		sm.Spinner.StopWithError("Unable to identify site")
		return fmt.Errorf("unable to identify site")
	}

	sm.Spinner.AddSuccessStatus("Identified Site")
	sm.Spinner.AddInfoStatus("Domain: %s", sm.Domain)
	sm.Spinner.AddInfoStatus("Directory: %s", sm.Directory)

	switch strings.ToLower(name) {
	case "php":
		return sm.updatePhp(value)
	default:
		sm.Spinner.StopWithError("Unknown setting name %s", name)
		return fmt.Errorf("unknown setting name")
	}
}

func (sm *SiteManager) updatePhp(version string) error {
	sm.PhpVersion = version
	if err := sm.validatePhpVersion(); err != nil {
		sm.Spinner.StopWithError("Failed to update site")
		return err
	}

	if err := sm.createSiteConfig(); err != nil {
		sm.Spinner.StopWithError("Failed to update site")
		return err
	}

	if err := sm.restartNginx(); err != nil {
		sm.Spinner.StopWithError("Failed to restart nginx")
		return err
	}

	config.SetStringData(fmt.Sprintf("web.sites.[%s].php_version", sm.Domain), sm.PhpVersion)

	sm.Spinner.AddInfoStatus("Updated to PHP %s", sm.PhpVersion)
	sm.Spinner.StopWithSuccess("Update Successful")

	return nil
}

func (sm *SiteManager) RemoveSite(identifier string) error {
	sm.Spinner.UpdatePhrase("Removing site")
	sm.Spinner.Start()

	sm.Spinner.AddInfoStatus("Using identifier: %s", identifier)
	if !sm.identifySite(identifier) {
		sm.Spinner.StopWithError("Unable to identify site")
		return fmt.Errorf("unable to identify site")
	}

	sm.Spinner.AddSuccessStatus("Identified Site")
	sm.Spinner.AddInfoStatus("Domain: %s", sm.Domain)
	sm.Spinner.AddInfoStatus("Directory: %s", sm.Directory)

	nginxPath := filepath.Join(constants.YerdWebDir, "nginx")
	files := []string{
		filepath.Join(constants.YerdWebDir, "certs", "sites", sm.Domain+".key"),
		filepath.Join(constants.YerdWebDir, "certs", "sites", sm.Domain+".crt"),
		filepath.Join(nginxPath, "sites-enabled", sm.Domain+".conf"),
	}

	utils.SystemdStopService("yerd-nginx")

	for _, file := range files {
		if err := utils.RemoveFile(file); err != nil {
			sm.Spinner.AddInfoStatus("Unable to remove %s", filepath.Base(file))
		} else {
			sm.Spinner.AddSuccessStatus("Removed %s", filepath.Base(file))
		}
	}

	if err := utils.SystemdStartService("yerd-nginx"); err != nil {
		sm.Spinner.AddInfoStatus("Unable to restart nginx")
	} else {
		sm.Spinner.AddSuccessStatus("Restarted Nginx")
	}

	hm := utils.NewHostsManager()
	hm.Remove(sm.Domain)

	config.Delete(fmt.Sprintf("web.sites.[%s]", sm.Domain))

	sm.Spinner.StopWithSuccess("Site Removed")
	return nil
}

func (sm *SiteManager) identifySite(identifier string) bool {
	path, _ := filepath.Abs(identifier)
	for _, site := range sm.WebConfig.Sites {
		if site.Domain == identifier || site.RootDirectory == path {
			sm.Domain = site.Domain
			sm.Directory = site.RootDirectory
			sm.PublicFolder = site.PublicDirectory
			sm.PhpVersion = site.PhpVersion
			sm.CrtFile = filepath.Join(constants.CertsDir, "sites", site.Domain+".crt")
			sm.KeyFile = filepath.Join(constants.CertsDir, "sites", site.Domain+".key")

			return true
		}
	}

	return false
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
		func() error { return siteManager.createCertificate() },
		func() error { return siteManager.createSiteConfig() },
		func() error { return siteManager.createHostsEntry() },
		func() error { return siteManager.restartNginx() },
		func() error { return siteManager.addToConfig() },
	)

	if err != nil {
		siteManager.Spinner.StopWithError("Failed to add site")
		return err
	}

	siteManager.Spinner.StopWithSuccess("Site Created!  https://%s", siteManager.Domain)

	return nil
}

func (sm *SiteManager) addToConfig() error {
	siteConfig := &config.SiteConfig{
		RootDirectory:   sm.Directory,
		PublicDirectory: sm.PublicFolder,
		PhpVersion:      sm.PhpVersion,
		Domain:          sm.Domain,
	}

	config.SetStruct(fmt.Sprintf("web.sites.[%s]", sm.Domain), siteConfig)

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
		}
	}

	if !utils.IsDirectory(filepath.Join(abs, siteManager.PublicFolder)) {
		siteManager.Spinner.AddErrorStatus("Provided public path is not a directory")
		return fmt.Errorf("public folder does not exist")
	}

	if siteManager.WebConfig.Sites != nil {
		for _, site := range siteManager.WebConfig.Sites {
			if site.RootDirectory == siteManager.Directory {
				siteManager.Spinner.AddErrorStatus("Directory is already registered")
				return fmt.Errorf("directory in use")
			}
		}
	}

	siteManager.Spinner.AddInfoStatus("Serving from /%s", siteManager.PublicFolder)

	return nil
}

func (sm *SiteManager) createCertificate() error {
	cm := NewCertificateManager()
	keyFile, certFile, err := cm.GenerateCert(sm.Domain, "yerd")
	if err != nil {
		sm.Spinner.AddErrorStatus("Unable to secure site")
		return err
	}

	sm.CrtFile = certFile
	sm.KeyFile = keyFile

	sm.Spinner.AddSuccessStatus("Site Secured Successfully")
	sm.Spinner.AddInfoStatus("https://%s/", sm.Domain)

	return nil
}

func (siteManager *SiteManager) validateDomain() error {
	if siteManager.Domain == "" {
		base := filepath.Base(siteManager.Directory)
		base = strings.ToLower(base) + ".test"
		siteManager.Domain = base
	}

	siteManager.Domain = strings.ToLower(siteManager.Domain)

	if siteManager.WebConfig.Sites != nil {
		for _, site := range siteManager.WebConfig.Sites {
			if site.Domain == siteManager.Domain {
				siteManager.Spinner.AddErrorStatus("Domain is already in use")
				siteManager.Spinner.AddInfoStatus("Use the -d flag to specify a custom domain")
				return fmt.Errorf("domain in use")
			}
		}
	}

	siteManager.Spinner.AddInfoStatus("Domain: %s", siteManager.Domain)

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
				break
			}
		}
	}

	_, installed := config.GetInstalledPhpInfo(siteManager.PhpVersion)
	if !installed {
		siteManager.Spinner.AddErrorStatus("PHP %s is not installed, or is not valid", siteManager.PhpVersion)
		return fmt.Errorf("php version not installed")
	}

	siteManager.Spinner.AddInfoStatus("Using PHP %s", siteManager.PhpVersion)

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
		"cert":        siteManager.CrtFile,
		"key":         siteManager.KeyFile,
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
