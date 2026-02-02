package apps

import (
	"fmt"
	"os"
)

var discovery_paths = []string{
	"/usr/share/applications",
	"$HOME/.local/share/applications",
	"/var/lib/flatpak/exports/share/applications",
}

type appService struct {
	Apps []Application
}

var AppServiceInstance = &appService{}

func (a *appService) Discover() (any, error) {
	discoveredApps := []Application{}

	for _, path := range discovery_paths {
		apps, err := a.discoverPath(path)
		if err != nil {
			fmt.Printf("Error discovering path %s: %s\n", path, err)
			continue
		}

		fmt.Printf("Apps in path %s: %d\n", path, len(apps))

		discoveredApps = append(discoveredApps, apps...)
	}

	a.Apps = discoveredApps

	fmt.Printf("All apps discovered: %d\n", len(discoveredApps))

	return nil, nil
}

func (a *appService) discoverPath(path string) ([]Application, error) {
	apps := make([]Application, 0)
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		file, err := os.ReadFile(path + "/" + entry.Name())

		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", path+"/"+entry.Name(), err)
			continue
		}

		application, err := ParseAppFile(file)
		if err != nil {
			fmt.Printf("Error parsing app file %s: %s\n", path+"/"+entry.Name(), err)
			continue
		}

		application.Path = path + "/" + entry.Name()
		apps = append(apps, application)
	}

	return apps, nil
}

func GetApplcationByPath(path string) *Application {
	for _, app := range AppServiceInstance.Apps {
		if app.Path == path {
			return &app
		}
	}

	return nil
}
