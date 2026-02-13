package handler

import (
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

		// if app.IconBase64 == "" {
		// 	continue
		// }

		rank := fuzzy.RankMatchNormalized(queryLowercase, strings.ToLower(app.Name))

		if rank >= 0 {
			var lastUsed int64
			if entry, ok := recentMap[app.Path]; ok {
				rank = adjustRankWithHistory(rank, entry, now)
				lastUsed = entry.LastUsed
			}
			results = append(results, Result{
				Label:      app.Name,
				Rank:       rank,
				Icon:       app.Icon,
				IconBase64: app.IconBase64,
				LastUsed:   lastUsed,
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
		if results[i].Rank != results[j].Rank {
			return results[i].Rank < results[j].Rank
		}
		return results[i].LastUsed > results[j].LastUsed
	})

	return HandlerResult{
		Items: results,
	}, nil
}

func adjustRankWithHistory(rank int, entry history.HistoryEntry, now int64) int {
	age := time.Unix(now, 0).Sub(time.Unix(entry.LastUsed, 0))
	boost := 0
	// recency buckets (lower rank is better)
	switch {
	case age <= time.Hour:
		boost += 40
	case age <= 24*time.Hour:
		boost += 25
	case age <= 7*time.Hour*24:
		boost += 12
	case age <= 30*time.Hour*24:
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

func LoadRecent() []history.HistoryEntry {
	entries, _ := history.Service.GetLast()

	return entries
}
