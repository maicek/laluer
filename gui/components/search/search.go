package search

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type Search struct {
	*gtk.Box

	Input *gtk.Entry
}

func NewSearch() *Search {
	search := &Search{
		Box: gtk.NewBox(gtk.OrientationHorizontal, 5),
	}

	search.Input = gtk.NewEntry()

	search.Input.SetAlignment(0.5)
	search.Input.SetSizeRequest(200, 50)

	search.Input.SetPlaceholderText("Search...")
	search.Input.SetHExpand(true)
	search.Input.SetHAlign(gtk.AlignFill)
	search.Input.AddCSSClass("Header-Search")

	search.Box.Append(search.Input)

	return search
}
