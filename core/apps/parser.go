package apps

import (
	"gopkg.in/ini.v1"
)

type DesktopEntry struct {
	Name        string `ini:"Name"`
	Comment     string `ini:"Comment"`
	GenericName string `ini:"GenericName"`
	Exec        string `ini:"Exec"`
	Icon        string `ini:"Icon"`
	NoDisplay   bool   `ini:"NoDisplay"`
	Hidden      bool   `ini:"Hidden"`
	Terminal    bool   `ini:"Terminal"`
}

type AppDotDesktop struct {
	DesktopEntry DesktopEntry `ini:"Desktop Entry"`
}

func ParseAppFile(file []byte) (Application, error) {
	iniFile, err := ini.Load(file)

	if err != nil {
		return Application{}, err
	}

	appDotDesktop := AppDotDesktop{}
	iniFile.MapTo(&appDotDesktop)

	app := Application{
		Name:        appDotDesktop.DesktopEntry.Name,
		Description: appDotDesktop.DesktopEntry.Comment,
		GenericName: appDotDesktop.DesktopEntry.GenericName,
		Exec:        appDotDesktop.DesktopEntry.Exec,
		Icon:        appDotDesktop.DesktopEntry.Icon,
		IconBase64:  "",
		NoDisplay:   appDotDesktop.DesktopEntry.NoDisplay,
		Hidden:      appDotDesktop.DesktopEntry.Hidden,
		Terminal:    appDotDesktop.DesktopEntry.Terminal,
	}

	return app, nil
}
