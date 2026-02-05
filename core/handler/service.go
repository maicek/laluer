package handler

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/maicek/laluer/core/apps"
)

type HandlerService struct{}

func bestSubtitle(genericName, description string) string {
	if genericName != "" {
		return genericName
	}
	return description
}

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
		if app.NoDisplay || app.Hidden {
			continue
		}

		nameRank := fuzzy.RankMatchNormalized(queryLowercase, strings.ToLower(app.Name))
		genericRank := fuzzy.RankMatchNormalized(queryLowercase, strings.ToLower(app.GenericName))
		descRank := fuzzy.RankMatchNormalized(queryLowercase, strings.ToLower(app.Description))

		// pick the best (lowest) rank among name, generic name, description
		rank := nameRank
		if genericRank >= 0 && (rank < 0 || genericRank < rank) {
			rank = genericRank
		}
		if descRank >= 0 && (rank < 0 || descRank < rank) {
			rank = descRank
		}

		if rank >= 0 {
			results = append(results, Result{
				Label:      app.Name,
				Rank:       rank,
				Icon:       app.Icon,
				IconBase64: app.IconBase64,
				IconMime:   app.IconMime,
				Subtitle:   bestSubtitle(app.GenericName, app.Description),
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
