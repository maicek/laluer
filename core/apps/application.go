package apps

type Application struct {
	Path        string
	Name        string
	Description string
	GenericName string
	Icon        string
	IconBase64  string
	NoDisplay   bool
	// Parameters needed for proper execution
	Exec     string
	Terminal bool
}
