package apps

type Application struct {
	Path        string
	Name        string
	Description string
	GenericName string
	Icon        string
	IconBase64  string
	IconMime    string
	NoDisplay   bool
	Hidden      bool
	// Parameters needed for proper execution
	Exec     string
	Terminal bool
}
