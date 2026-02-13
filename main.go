package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/maicek/laluer/core/apps"
	"github.com/maicek/laluer/core/handler"
	"github.com/maicek/laluer/core/history"
	"github.com/maicek/laluer/gui"
	"github.com/maicek/laluer/gui/components/item"
)

//go:embed style.css
var styleCSS string

type LaluerApp struct {
	handler *handler.HandlerService
	ctx     context.Context
	app     *gtk.Application
	window  *gui.Laluer
}

func MakeApp(ctx context.Context) *LaluerApp {
	laluer := &LaluerApp{ctx: ctx}

	laluer.app = gtk.NewApplication("com.github.maicek.laluer", gio.ApplicationFlagsNone)

	laluer.app.ConnectActivate(func() {
		display := gdk.DisplayGetDefault()
		prov := loadCSSFromFileOrFallback("style.css", styleCSS)
		gtk.StyleContextAddProviderForDisplay(
			display, prov,
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
		)
		startCSSHotReload(laluer.ctx, display, prov, "style.css", styleCSS)

		go laluer.init()

		laluer.window = gui.NewLaluerWindow(ctx, laluer.app)
		laluer.window.Show()

		input := laluer.window.GetInput()
		input.ConnectChanged(func() {
			laluer.HandleInput(input.Text())
		})
	})

	return laluer
}

func (l *LaluerApp) init() {
	go apps.AppServiceInstance.Discover()
	go history.Init()
	// l.handler = &handler.HandlerService{}
}

func loadCSS(content string) *gtk.CSSProvider {
	prov := gtk.NewCSSProvider()
	prov.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		// Optional line parsing routine.
		loc := sec.StartLocation()
		lines := strings.Split(content, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})
	prov.LoadFromString(content)
	return prov
}

func loadCSSFromFileOrFallback(path string, fallback string) *gtk.CSSProvider {
	content, err := os.ReadFile(path)
	if err != nil {
		return loadCSS(fallback)
	}
	return loadCSS(string(content))
}

func startCSSHotReload(ctx context.Context, display *gdk.Display, prov *gtk.CSSProvider, path string, fallback string) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	lastMod := stat.ModTime()
	ticker := time.NewTicker(250 * time.Millisecond)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				st, err := os.Stat(path)
				if err != nil {
					continue
				}

				mod := st.ModTime()
				if !mod.After(lastMod) {
					continue
				}
				lastMod = mod

				content, err := os.ReadFile(path)
				if err != nil {
					content = []byte(fallback)
				}

				css := string(content)
				glib.IdleAdd(func() {
					prov.LoadFromString(css)
					gtk.StyleContextAddProviderForDisplay(display, prov, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
				})
			}
		}
	}()
}

func (l *LaluerApp) HandleInput(input string) {
	results, err := l.handler.Handle(handler.SearchParams{
		Query: input,
	})
	if err != nil {
		fmt.Println(err)
	}

	l.window.SetResults(l.HandleResultToGui(results.Items))
}

func (l *LaluerApp) HandleResultToGui(results []handler.Result) []item.ItemData {
	items := make([]item.ItemData, len(results))
	for i, result := range results {
		items[i] = item.ItemData{
			Name:        result.Label,
			Description: result.Subtitle,
		}
	}

	return items
}

func (l *LaluerApp) HandleHistoryToGui(results []history.HistoryEntry) []item.ItemData {
	items := make([]item.ItemData, len(results))

	return items
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := MakeApp(ctx)

	go func() {
		<-ctx.Done()
		glib.IdleAdd(app.app.Quit)
	}()

	if code := app.app.Run(os.Args); code > 0 {
		cancel()
		os.Exit(code)
	}
}
