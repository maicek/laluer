package results

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/maicek/laluer/gui/components/item"
)

type Results struct {
	*gtk.Box

	ActiveIndex  int
	items        []*item.ResultItem
	ActiveItem   *item.ResultItem
	ItemSelected func(item *item.ResultItem)
}

type ResultsInit struct {
	ItemSelected func(item *item.ResultItem)
}

func NewResults(init ResultsInit) *Results {
	results := &Results{
		Box:          gtk.NewBox(gtk.OrientationVertical, 5),
		ActiveIndex:  0,
		ItemSelected: init.ItemSelected,
	}

	results.Box.AddCSSClass("Body-Results")

	return results
}

func (r *Results) SetResults(results []item.ItemData) {
	for _, item := range r.items {
		r.Box.Remove(item)
	}

	r.ActiveIndex = 0
	r.items = make([]*item.ResultItem, len(results))

	for i, itm := range results {
		r.items[i] = item.NewResultItem(itm)
		r.Box.Append(r.items[i])
	}

	if len(r.items) > 0 {
		r.ActiveItem = r.items[0]
		r.ActiveItem.Select()
	} else {
		r.ActiveItem = nil
	}
}

func (r *Results) Next() {
	if r.ActiveIndex < len(r.items)-1 {
		r.ActiveItem.Deselect()
		r.ActiveIndex++
		r.ActiveItem = r.items[r.ActiveIndex]
		r.ActiveItem.Select()
	}
}

func (r *Results) Previous() {
	if r.ActiveIndex > 0 {
		r.ActiveItem.Deselect()
		r.ActiveIndex--
		r.ActiveItem = r.items[r.ActiveIndex]
		r.ActiveItem.Select()
	}
}

func (r *Results) Select() {
	if r.ActiveItem != nil {
		if r.ItemSelected != nil {
			r.ItemSelected(r.ActiveItem)
		}
	}
}
