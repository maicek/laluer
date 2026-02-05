package handler

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/maicek/laluer/core/apps"
)

type HandlerService struct{}

type SearchParams struct {
	Query string `json:"query"`
}

type HandlerResult struct {
	Items []Result `json:"items"`
}

// ranked search

func (h *HandlerService) Handle(searchParams SearchParams) (HandlerResult, error) {
	results := make([]Result, 0)
	queryLowercase := strings.ToLower(searchParams.Query)

	for _, app := range apps.AppServiceInstance.Apps {
		// Skip hidden apps
		if app.NoDisplay {
			continue
		}

		rank := fuzzy.RankMatchNormalized(queryLowercase, strings.ToLower(app.Name))

		if rank >= 0 {
			results = append(results, Result{
				Label:      app.Name,
				Rank:       rank,
				Icon:       app.Icon,
				IconBase64: app.IconBase64,
				Subtitle:   app.Description,
				Action: Action{
					Event: "run",
					Payload: struct {
						Path string `json:"path"`
					}{
						Path: app.Path,
					},
				},
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Rank < results[j].Rank
	})

	return HandlerResult{
		Items: results,
	}, nil
}
