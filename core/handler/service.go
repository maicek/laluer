package handler

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/maicek/laluer/core/apps"
	"github.com/maicek/laluer/core/history"
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
	recentEntries, _ := history.Service.GetLast()
	recentMap := make(map[string]history.HistoryEntry, len(recentEntries))
	for _, entry := range recentEntries {
		if entry.Type != history.ENTRY_TYPE_APP {
			continue
		}
		recentMap[entry.EntryName] = entry
	}
	now := time.Now().Unix()

	for _, app := range apps.AppServiceInstance.Apps {
		// Skip hidden apps
		if app.NoDisplay {
			continue
		}

		rank := fuzzy.RankMatchNormalized(queryLowercase, strings.ToLower(app.Name))

		if rank >= 0 {
			if entry, ok := recentMap[app.Path]; ok {
				rank = adjustRankWithHistory(rank, entry, now)
			}
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

func adjustRankWithHistory(rank int, entry history.HistoryEntry, now int64) int {
	ageSeconds := now - entry.LastUsed
	boost := 0
	// recency buckets (lower rank is better)
	switch {
	case ageSeconds <= 3600:
		boost += 40
	case ageSeconds <= 24*3600:
		boost += 25
	case ageSeconds <= 7*24*3600:
		boost += 12
	case ageSeconds <= 30*24*3600:
		boost += 6
	}
	if entry.UseCount > 0 {
		if entry.UseCount > 6 {
			boost += 18
		} else {
			boost += entry.UseCount * 3
		}
	}
	if rank-boost < 0 {
		return 0
	}
	return rank - boost
}

func LoadRecent() {
	entries, _ := history.Service.GetLast()

	fmt.Printf("Last: %+v \n", entries)

}
