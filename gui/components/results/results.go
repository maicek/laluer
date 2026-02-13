package results

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/maicek/laluer/gui/components/item"
)

type Results struct {
	*gtk.Box

	activeIndex int
	items       []*item.ResultItem
}

func NewResults() *Results {
	results := &Results{
		Box:         gtk.NewBox(gtk.OrientationVertical, 5),
		activeIndex: 0,
	}

	results.Box.AddCSSClass("Body-Results")

	return results
}

func (r *Results) SetResults(results []item.ItemData) {
	for _, item := range r.items {
		r.Box.Remove(item)
	}

	r.activeIndex = 0
	r.items = make([]*item.ResultItem, len(results))

	for i, itm := range results {
		r.items[i] = item.NewResultItem(itm)
		r.Box.Append(r.items[i])
	}
}
