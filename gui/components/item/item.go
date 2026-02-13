package item

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

type ResultItem struct {
	*gtk.Box

	Name        *gtk.Label
	Description *gtk.Label
	Icon        *gtk.Image
}

type ItemData struct {
	Name        string
	Description string
	Icon        string
}

func NewResultItem(data ItemData) *ResultItem {
	item := &ResultItem{
		Box:         gtk.NewBox(gtk.OrientationVertical, 0),
		Name:        gtk.NewLabel(data.Name),
		Description: gtk.NewLabel(data.Description),
	}

	item.SetSizeRequest(0, 50)
	// item.SetHExpand(true)

	item.AddCSSClass("Item")

	item.Name.AddCSSClass("Item__name")
	item.Name.SetHAlign(gtk.AlignStart)
	item.Name.SetHExpand(true)
	item.Name.SetSingleLineMode(true)
	item.Name.SetEllipsize(pango.EllipsizeEnd)

	// if data.Description == "" {
	// 	item.Description.SetVisible(false)
	// }

	item.Description.AddCSSClass("Item__description")
	item.Description.SetHAlign(gtk.AlignStart)
	item.Description.SetHExpand(true)
	item.Description.SetWrap(true)
	item.Description.SetWrapMode(pango.WrapWordChar)
	item.Description.SetEllipsize(pango.EllipsizeEnd)

	// layout
	left := gtk.NewBox(gtk.OrientationHorizontal, 6)
	right := gtk.NewBox(gtk.OrientationVertical, 2)

	right.SetHExpand(true)

	right.Append(item.Name)
	right.Append(item.Description)

	item.Box.Append(left)
	item.Box.Append(right)

	return item
}
