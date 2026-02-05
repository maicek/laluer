package apps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var HOME = os.Getenv("HOME")

var discovery_paths = []string{
	"/usr/share/applications",
	HOME + "/.local/share/applications",
	HOME + "/.local/share/applications",
	"/var/lib/flatpak/exports/share/applications",
}

type appService struct {
	Apps []Application
	// map app name to icon path
	Icons map[string]string
}

var AppServiceInstance = &appService{}

func (a *appService) Discover() (any, error) {
	discoveredApps := []Application{}
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for _, path := range discovery_paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			apps, err := a.discoverPath(path)
			if err != nil {
				fmt.Printf("Error discovering path %s: %s\n", path, err)
				return
			}

			fmt.Printf("Apps in path %s: %d\n", path, len(apps))

			mu.Lock()
			discoveredApps = append(discoveredApps, apps...)
			mu.Unlock()
		}(path)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		apps := a.discoverWinePath()
		mu.Lock()
		discoveredApps = append(discoveredApps, apps...)
		mu.Unlock()
	}()

	wg.Wait()

	a.Apps = discoveredApps

	fmt.Printf("All apps discovered: %d\n", len(discoveredApps))

	a.DiscoverAppIcons()

	return nil, nil
}

func (a *appService) discoverWinePath() []Application {
	WINE_PATH := HOME + "/.local/share/applications/wine/Programs"

	entries, err := os.ReadDir(WINE_PATH)

	if err != nil {
		fmt.Printf("Error reading wine path %s: %s\n", WINE_PATH, err)
		return []Application{}
	}

	discoveredApps := make([]Application, 0)
	// iterate over folders
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirPath := filepath.Join(WINE_PATH, entry.Name())
		info, err := os.Stat(dirPath)
		if err != nil || !info.IsDir() {
			fmt.Printf("Skipping wine path %s: %v\n", dirPath, err)
			continue
		}
		apps, err := a.discoverPath(dirPath)
		if err != nil {
			fmt.Printf("Error discovering path %s: %s\n", dirPath, err)
			continue
		}
		discoveredApps = append(discoveredApps, apps...)
	}

	// TODO: parse wine entries
	fmt.Printf("entries: %+v", entries)

	return discoveredApps
}

func (a *appService) discoverPath(path string) ([]Application, error) {
	apps := make([]Application, 0)
	mu := sync.RWMutex{}
	var wg sync.WaitGroup
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".desktop") {
			continue
		}
		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()
			file, err := os.ReadFile(path + "/" + entry.Name())

			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", path+"/"+entry.Name(), err)
				return
			}

			application, err := ParseAppFile(file)
			if err != nil {
				fmt.Printf("Error parsing app file %s: %s\n", path+"/"+entry.Name(), err)
				return
			}

			application.Path = path + "/" + entry.Name()
			application.IconBase64 = ""
			mu.Lock()
			apps = append(apps, application)
			mu.Unlock()
		}(entry)
	}

	wg.Wait()
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
