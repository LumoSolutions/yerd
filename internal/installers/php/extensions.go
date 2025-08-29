package php

import (
	"fmt"
	"slices"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
)

type ExtensionManager struct {
	Version  string
	Info     *config.PhpInfo
	Cached   bool
	Config   bool
	Rebuild  bool
	ToAdd    []string
	ToRemove []string
}

func NewExtensionManager(version string, data *config.PhpInfo, cached, config, rebuild bool) *ExtensionManager {
	return &ExtensionManager{
		Version: version,
		Info:    data,
		Cached:  cached,
		Config:  config,
		Rebuild: rebuild,
	}
}

func (ext *ExtensionManager) RunAction(action string, extensions []string) error {
	switch action {
	case "list":
		ext.listExtensions()
		return nil

	case "add":
		if err := ext.addExtensions(extensions); err != nil {
			return err
		}

	case "remove":
		if err := ext.removeExtensions(extensions); err != nil {
			return err
		}

	default:
		fmt.Printf("Error: Invalid action '%s'. Use 'add' or 'remove'\n", action)
		return fmt.Errorf("invalid action")
	}

	ext.saveConfig()
	if err := ext.handleRebuild(action, extensions); err != nil {
		fmt.Printf("Failed to rebuild PHP %s, please try again via command: \n", ext.Version)
		fmt.Printf(" sudo yerd php %s rebuild\n\n", ext.Version)
		return err
	}

	return nil
}

func (ext *ExtensionManager) listExtensions() error {
	fmt.Printf("PHP %s Extensions:\n\n", ext.Version)

	fmt.Println("✓ INSTALLED:")
	utils.PrintExtensionsGrid(ext.Info.Extensions)

	if len(ext.Info.AddExtensions) > 0 {
		fmt.Println("\n✓ TO BE ADDED:")
		utils.PrintExtensionsGrid(ext.Info.AddExtensions)
	}

	if len(ext.Info.RemoveExtensions) > 0 {
		fmt.Println("\n✗ TO BE REMOVED:")
		utils.PrintExtensionsGrid(ext.Info.RemoveExtensions)
	}

	fmt.Println("\nAVAILABLE:")
	all := constants.GetAvailableExtensions()
	all = utils.RemoveItems(all, ext.Info.Extensions...)
	all = utils.RemoveItems(all, ext.Info.AddExtensions...)
	utils.PrintExtensionsGrid(all)

	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Printf("  yerd php %s extensions add <extensions>        # Add Extensions\n", ext.Version)
	fmt.Printf("  yerd php %s extensions remove <extensions>     # Remove Extensions\n", ext.Version)
	fmt.Printf("  yerd php %s extensions add <extensions> -r     # Add Extensions & Rebuild PHP\n", ext.Version)

	return nil
}

func (ext *ExtensionManager) addExtensions(extensions []string) error {
	valid, invalid := constants.ValidateExtensions(extensions)
	if len(invalid) > 0 {
		utils.PrintInvalidExtensionsWithSuggestions(invalid)
		return fmt.Errorf("invalid extensions")
	}

	toAdd := []string{}

	for _, item := range valid {
		if slices.Contains(ext.Info.Extensions, item) && !slices.Contains(ext.Info.RemoveExtensions, item) {
			fmt.Printf("ℹ️  Extension %s is already installed\n", item)
		} else {
			toAdd = append(toAdd, item)
		}
	}

	if len(toAdd) > 0 {
		ext.displayInstallationNotes(toAdd)
	}

	ext.Info.AddExtensions = utils.AddUnique(ext.Info.AddExtensions, toAdd...)
	ext.Info.RemoveExtensions = utils.RemoveItems(ext.Info.RemoveExtensions, toAdd...)
	ext.Info.AddExtensions = utils.RemoveItems(ext.Info.AddExtensions, ext.Info.Extensions...)

	return nil
}

func (ext *ExtensionManager) removeExtensions(extensions []string) error {
	valid, invalid := constants.ValidateExtensions(extensions)
	if len(invalid) > 0 {
		utils.PrintInvalidExtensionsWithSuggestions(invalid)
		return fmt.Errorf("invalid extensions")
	}

	toRemove := []string{}

	for _, item := range valid {
		if !slices.Contains(ext.Info.Extensions, item) && !slices.Contains(ext.Info.AddExtensions, item) {
			fmt.Printf("ℹ️  Extension %s is not installed\n", item)
		} else {
			toRemove = append(toRemove, item)
		}
	}

	ext.Info.RemoveExtensions = utils.AddUnique(ext.Info.RemoveExtensions, toRemove...)
	ext.Info.AddExtensions = utils.RemoveItems(ext.Info.AddExtensions, toRemove...)

	config.SetStruct(fmt.Sprintf("php.[%s]", ext.Info.Version), ext.Info)

	return nil
}

func (ext *ExtensionManager) saveConfig() {
	config.SetStruct(fmt.Sprintf("php.[%s]", ext.Info.Version), ext.Info)
}

func (ext *ExtensionManager) handleRebuild(action string, extensions []string) error {
	if ext.Rebuild {
		if len(ext.Info.AddExtensions) == 0 && len(ext.Info.RemoveExtensions) == 0 {
			fmt.Println("ℹ️  Nothing to add or remove, skipping rebuild")
			return nil
		}

		if len(ext.Info.AddExtensions) > 0 {
			fmt.Println("\n✓ TO BE ADDED:")
			utils.PrintExtensionsGrid(ext.Info.AddExtensions)
		}

		if len(ext.Info.RemoveExtensions) > 0 {
			fmt.Println("\n✗ TO BE REMOVED:")
			utils.PrintExtensionsGrid(ext.Info.RemoveExtensions)
		}

		if err := RunRebuild(ext.Info, ext.Cached, ext.Config); err != nil {
			fmt.Printf("Failed to rebuild PHP %s, %v\n", ext.Version, err)
			return err
		}
	} else {
		if action == "add" {
			fmt.Printf("These extensions will be added to PHP %s on the next rebuild\n", ext.Version)
		} else {
			fmt.Printf("These extensions will be removed from PHP %s on the next rebuild\n", ext.Version)
		}

		utils.PrintExtensionsGrid(extensions)
		fmt.Println()

		fmt.Println("ℹ️  These changes won't apply until PHP is rebuilt")
		fmt.Println("ℹ️  PHP can be rebuilt with the following command:")
		fmt.Printf("\n sudo yerd php %s rebuild\n\n", ext.Version)
	}

	return nil
}

func (ext *ExtensionManager) displayInstallationNotes(extensions []string) {
	dm, _ := manager.NewDependencyManager()
	pm := dm.GetPackageManager()

	hasSpecialNotes := false

	for _, extName := range extensions {
		if phpExt, exists := constants.GetExtension(extName); exists {
			for _, dep := range phpExt.Dependencies {
				if depConfig, exists := constants.GetDependencyConfig(dep); exists {
					if depConfig.RequiresSpecialHandling != nil && depConfig.RequiresSpecialHandling[pm] {
						if !hasSpecialNotes {
							fmt.Println("\n⚠️  SPECIAL INSTALLATION REQUIREMENTS:")
							fmt.Println("─────────────────────────────────────")
							hasSpecialNotes = true
						}

						fmt.Printf("\n📦 Extension '%s' (dependency: %s):\n", extName, dep)
						if note, exists := depConfig.InstallNotes[pm]; exists {
							fmt.Printf("   %s\n", note)
						}
					}
				}
			}
		}
	}

	if hasSpecialNotes {
		fmt.Println("\n─────────────────────────────────────")
		fmt.Println("ℹ️  Please install these dependencies before rebuilding PHP")
		fmt.Println()
	}
}
