package apps

import (
	"encoding/base64"
	"fmt"
	"os"
)

var icon_paths = []string{
	"/usr/share/icons/",
}

func FindIcon(iconName string) string {
	if iconName == "" {
		return ""
	}

	for _, basePath := range icon_paths {
		result := findIconRecursive(iconName, basePath)
		if result != "" {
			return result
		}
	}

	extensions := []string{".png", ".svg", ".xpm"}
	for _, ext := range extensions {
		nameWithExt := iconName + ext
		for _, basePath := range icon_paths {
			result := findIconRecursive(nameWithExt, basePath)
			if result != "" {
				return result
			}
		}
	}

	return ""
}

func findIconRecursive(iconName string, dirPath string) string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		fullPath := dirPath + entry.Name()
		if entry.IsDir() {
			result := findIconRecursive(iconName, fullPath+"/")
			if result != "" {
				return result
			}
		} else if entry.Name() == iconName {
			fmt.Println("Found icon:", fullPath)
			return iconToBase64(fullPath)
		}
	}
	return ""
}

func FindIconInPath(iconName string, path string) string {
	icon_path := path + iconName
	if _, err := os.Stat(icon_path); err == nil {
		fmt.Println("Found icon:", icon_path, "in path:", path)
		return iconToBase64(icon_path)
	}
	return ""
}

func iconToBase64(iconPath string) string {
	data, err := os.ReadFile(iconPath)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}
