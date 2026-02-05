package apps

import (
	"encoding/base64"
	"fmt"
	"os"
)

// todo: enhance icon discovery
// todo: support other icon sizes
// todo: support fetching proper system theme
// todo: support caching
// todo: support other icon formats
// todo: Replace base64 encoding with direct file serving via providing embedded FS to frontend (maicek)

var icon_paths = []string{
	"/usr/share/icons",
	"~/.local/share/icons",
	"~/.icons",
}

func (a *appService) DiscoverAppIcons() {
	// read icon cache
	// skip rn.

	// check system theme
	// temp
	theme := "hicolor"

	for i, app := range a.Apps {
		// app.Icon
		iconPath := fmt.Sprintf("%s/%s/%s/apps/%s.%s", icon_paths[0], theme, "scalable", app.Icon, "svg")
		fmt.Printf("App icon path for %s is %s\n", app.Name, iconPath)
		a.Apps[i].IconBase64 = iconToBase64(iconPath)
	}
}

func iconToBase64(iconPath string) string {
	data, err := os.ReadFile(iconPath)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}
