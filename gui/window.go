package gui

import (
	"context"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/maicek/laluer/gui/components/item"
	"github.com/maicek/laluer/gui/components/results"
	"github.com/maicek/laluer/gui/components/search"
)

type Laluer struct {
	*gtk.Application
	window *gtk.ApplicationWindow

	view *LaluerView
}

type LaluerView struct {
	*gtk.Box
	Header struct {
		*gtk.Box
		Search *search.Search
	}
	Body *gtk.Box

	Results *results.Results
}

func NewLaluerWindow(ctx context.Context, app *gtk.Application) *Laluer {
	l := &Laluer{
		Application: app,
	}

	l.window = gtk.NewApplicationWindow(app)
	l.window.SetDefaultSize(500, 400)
	l.window.SetResizable(false)
	l.window.SetTitle("Laluer")
	l.window.SetIconName("com.github.maicek.laluer")
	l.window.SetDecorated(false)

	l.view = newLaluerView()
	l.window.SetChild(l.view)

	return l
}

func newLaluerView() *LaluerView {
	view := LaluerView{
		Box:  gtk.NewBox(gtk.OrientationVertical, 5),
		Body: gtk.NewBox(gtk.OrientationVertical, 5),
	}

	view.Box.AddCSSClass("App")

	header := struct {
		*gtk.Box
		Search *search.Search
	}{
		Box: gtk.NewBox(gtk.OrientationHorizontal, 5),
	}
	header.Search = search.NewSearch()

	view.Header = header

	scroller := gtk.NewScrolledWindow()
	scroller.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scroller.SetChild(view.Body)
	scroller.SetVExpand(true)

	view.Header.Append(header.Search)
	view.Results = results.NewResults()
	view.Body.Append(view.Results)

	view.Box.Append(view.Header)
	view.Box.Append(scroller)

	return &view
}

func (l *Laluer) Show() {
	l.window.Widget.SetVisible(true)
}

func (l *Laluer) SetResults(results []item.ItemData) {
	l.view.Results.SetResults(results)
}

func (l *Laluer) GetInput() *gtk.Entry {
	return l.view.Header.Search.Input
}
