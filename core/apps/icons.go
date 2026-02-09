package apps

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"
)

// todo: support fetching proper system theme
// todo: support caching
// todo: Replace base64 encoding with direct file serving via providing embedded FS to frontend (maicek)

type IconPathSearchPattern struct {
	Path   string
	Format string
}

type IconFormatSearchPattern struct {
	Size      string
	Extension string
}

var ICON_PATHS = []IconPathSearchPattern{
	{Path: "/usr/share/icons", Format: "%s/%s/%s/apps/%s.%s"},
	{Path: HOME + "/.local/share/icons", Format: "%s/%s/%s/apps/%s.%s"},
	{Path: HOME + "/.icons", Format: "%s/%s/%s/apps/%s.%s"},
	{Path: "/var/lib/flatpak/exports/share/icons", Format: "%s/%s/%s/apps/%s.%s"},
	{Path: "/usr/share/pixmaps", Format: "%[1]s/%[4]s.%[5]s"},
}

var SEARCH_PATTERNS = []IconFormatSearchPattern{
	{Size: "scalable", Extension: "svg"},
	{Size: "32x32", Extension: "png"},
	{Size: "48x48", Extension: "png"},
	{Size: "64x64", Extension: "png"},
	{Size: "128x128", Extension: "png"},
	{Size: "256x256", Extension: "png"},
	{Size: "512x512", Extension: "png"},
}

func (a *appService) DiscoverAppIcons() {
	var wg sync.WaitGroup
	for i, app := range a.Apps {
		wg.Add(1)
		go func(app Application) {
			a.Apps[i].IconBase64 = discoverappIcon(app)
			wg.Done()
		}(app)
	}

	wg.Wait()
}

func searchPathForIcons(path string, format string, app Application) string {
	theme := "hicolor"

	for _, pattern := range SEARCH_PATTERNS {
		iconPath := fmt.Sprintf(format, path, theme, pattern.Size, app.Icon, pattern.Extension)
		_, err := os.Stat(iconPath)

		if err == nil {
			return iconPath
		}
	}
	return ""
}

func discoverappIcon(app Application) string {
	for _, path := range ICON_PATHS {
		iconPath := searchPathForIcons(path.Path, path.Format, app)
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
