package apps

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var iconBaseDirs = []string{
	HOME + "/.local/share/icons",
	HOME + "/.icons",
	"/usr/share/icons",
}

var pixmapDirs = []string{
	"/usr/share/pixmaps",
	HOME + "/.local/share/pixmaps",
}

var iconExtensions = []string{"svg", "png", "xpm"}

// getIconThemeName detects the current icon theme via gsettings
func getIconThemeName() string {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme").Output()
	if err == nil {
		name := strings.TrimSpace(string(out))
		name = strings.Trim(name, "'\"")
		if name != "" {
			return name
		}
	}
	return "hicolor"
}

// getThemeInherits parses index.theme to get parent themes
func getThemeInherits(themeName string) []string {
	for _, baseDir := range iconBaseDirs {
		indexPath := filepath.Join(baseDir, themeName, "index.theme")
		f, err := os.Open(indexPath)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Inherits=") {
				val := strings.TrimPrefix(line, "Inherits=")
				parts := strings.Split(val, ",")
				var result []string
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p != "" {
						result = append(result, p)
					}
				}
				return result
			}
		}
	}
	return nil
}

// getThemeAppDirs returns all "apps" subdirectories for a given theme
// sorted by preference: scalable first, then larger sizes
func getThemeAppDirs(themeName string) []string {
	var dirs []string

	// preferred subdirs ordered by quality
	preferredSubdirs := []string{
		"apps/scalable",
		"apps/scalable@2",
		"apps/symbolic",
		"scalable/apps",
		"256x256/apps",
		"128x128/apps",
		"96x96/apps",
		"64x64/apps",
		"48x48/apps",
		"32x32/apps",
		"24x24/apps",
		"22x22/apps",
		"16x16/apps",
	}

	for _, baseDir := range iconBaseDirs {
		themeDir := filepath.Join(baseDir, themeName)
		for _, subdir := range preferredSubdirs {
			fullPath := filepath.Join(themeDir, subdir)
			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				dirs = append(dirs, fullPath)
			}
		}
	}

	return dirs
}

// buildThemeChain creates the full inheritance chain: current -> parents -> hicolor
func buildThemeChain(startTheme string) []string {
	seen := make(map[string]bool)
	var chain []string

	var walk func(theme string)
	walk = func(theme string) {
		if seen[theme] {
			return
		}
		seen[theme] = true
		chain = append(chain, theme)
		for _, parent := range getThemeInherits(theme) {
			walk(parent)
		}
	}

	walk(startTheme)

	// ensure hicolor is always last fallback
	if !seen["hicolor"] {
		chain = append(chain, "hicolor")
	}

	return chain
}

// lookupIcon finds an icon file for the given icon name using the XDG icon theme spec
func lookupIcon(iconName string, themeChain []string, themeAppDirsCache map[string][]string) string {
	// if icon is an absolute path, use it directly
	if filepath.IsAbs(iconName) {
		if _, err := os.Stat(iconName); err == nil {
			return iconName
		}
		return ""
	}

	// strip extension if someone put one in the .desktop file
	iconName = strings.TrimSuffix(iconName, ".svg")
	iconName = strings.TrimSuffix(iconName, ".png")
	iconName = strings.TrimSuffix(iconName, ".xpm")

	// search through theme chain
	for _, theme := range themeChain {
		dirs, ok := themeAppDirsCache[theme]
		if !ok {
			dirs = getThemeAppDirs(theme)
			themeAppDirsCache[theme] = dirs
		}
		for _, dir := range dirs {
			for _, ext := range iconExtensions {
				path := filepath.Join(dir, iconName+"."+ext)
				if _, err := os.Stat(path); err == nil {
					return path
				}
			}
		}
	}

	// fallback: search pixmap directories
	for _, dir := range pixmapDirs {
		for _, ext := range iconExtensions {
			path := filepath.Join(dir, iconName+"."+ext)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
		// some pixmaps don't have standard extensions
		path := filepath.Join(dir, iconName)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func iconToBase64(iconPath string) (string, string) {
	if iconPath == "" {
		return "", ""
	}

	data, err := os.ReadFile(iconPath)
	if err != nil {
		return "", ""
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	ext := strings.ToLower(filepath.Ext(iconPath))
	var mime string
	switch ext {
	case ".svg":
		mime = "image/svg+xml"
	case ".png":
		mime = "image/png"
	case ".xpm":
		mime = "image/x-xpixmap"
	default:
		mime = "image/png"
	}

	return encoded, mime
}

func (a *appService) DiscoverAppIcons() {
	themeName := getIconThemeName()
	fmt.Printf("Detected icon theme: %s\n", themeName)

	themeChain := buildThemeChain(themeName)
	fmt.Printf("Theme chain: %v\n", themeChain)

	// pre-build the theme app dirs cache (single-threaded, no race)
	themeAppDirsCache := make(map[string][]string)
	for _, theme := range themeChain {
		themeAppDirsCache[theme] = getThemeAppDirs(theme)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range a.Apps {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			app := &a.Apps[idx]
			if app.Icon == "" {
				return
			}

			iconPath := lookupIcon(app.Icon, themeChain, themeAppDirsCache)
			if iconPath == "" {
				return
			}

			encoded, mime := iconToBase64(iconPath)
			mu.Lock()
			app.IconBase64 = encoded
			app.IconMime = mime
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("Icon discovery complete\n")
}
