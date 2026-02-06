package apps

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
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
	HOME + "/.local/share/icons",
	HOME + "/.icons",
}

func (a *appService) DiscoverAppIcons() {
	// read icon cache
	// skip rn.

	// check system theme
	// temp
	// theme := "hicolor"

	for i, app := range a.Apps {
		a.Apps[i].IconBase64 = discoverappIcon(app)
	}
}

type SearchPattern struct {
	Size      string
	Extension string
}

var SEARCH_PATTERNS = []SearchPattern{
	{Size: "scalable", Extension: "svg"},
	{Size: "32x32", Extension: "png"},
	{Size: "48x48", Extension: "png"},
	{Size: "64x64", Extension: "png"},
	{Size: "128x128", Extension: "png"},
	{Size: "256x256", Extension: "png"},
	{Size: "512x512", Extension: "png"},
}

func discoverappIcon(app Application) string {
	theme := "hicolor"

	searchPath := func(path string) string {
		for _, pattern := range SEARCH_PATTERNS {
			iconPath := fmt.Sprintf("%s/%s/%s/apps/%s.%s", path, theme, pattern.Size, app.Icon, pattern.Extension)
			_, err := os.Stat(iconPath)

			if err == nil {
				return iconPath
			}
		}
		return ""
	}

	for _, path := range icon_paths {
		iconPath := searchPath(path)
		if iconPath == "" {
			continue
		}

		return iconToBase64(iconPath)
	}

	return ""
}

func iconToBase64(iconPath string) string {
	// get file description
	file, err := os.Open(iconPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	data, err := os.ReadFile(iconPath)
	if err != nil {
		return ""
	}

	if string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		fmt.Printf("Encode %s to base64\n", iconPath)
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Printf("Error decoding icon %s: %s\n", iconPath, err)
			return ""
		}

		return encodePNGToBase64(img)
	}

	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(data)
}

func encodePNGToBase64(img image.Image) string {
	var buf bytes.Buffer
	png.Encode(&buf, img)

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}
