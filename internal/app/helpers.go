package app

import (
	"sort"
	"strings"
)

type scoredPaletteItem struct {
	item  paletteItem
	score int
}

func filterPaletteItems(items []paletteItem, query string) []paletteItem {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return items
	}

	scored := make([]scoredPaletteItem, 0, len(items))
	for _, it := range items {
		score, ok := paletteMatchScore(q, it)
		if !ok {
			continue
		}
		scored = append(scored, scoredPaletteItem{item: it, score: score})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].score < scored[j].score
	})

	filtered := make([]paletteItem, len(scored))
	for i, scoredItem := range scored {
		filtered[i] = scoredItem.item
	}
	return filtered
}

func paletteMatchScore(queryLower string, item paletteItem) (int, bool) {
	if queryLower == "" {
		return 0, true
	}

	label := strings.ToLower(item.label)
	desc := strings.ToLower(item.description)

	bestScore := 0
	matched := false

	if score, ok := fuzzyScoreLower(queryLower, label); ok {
		matched = true
		if strings.Contains(label, queryLower) {
			score -= 5
		}
		bestScore = score
	}

	if score, ok := fuzzyScoreLower(queryLower, desc); ok {
		score += 15
		if strings.Contains(desc, queryLower) {
			score -= 3
		}
		if !matched || score < bestScore {
			matched = true
			bestScore = score
		}
	}

	return bestScore, matched
}

func fuzzyScoreLower(query, target string) (int, bool) {
	if query == "" {
		return 0, true
	}

	qRunes := []rune(query)
	tRunes := []rune(target)
	if len(qRunes) == 0 {
		return 0, true
	}

	score := 0
	lastIdx := -1
	searchFrom := 0

	for _, qc := range qRunes {
		found := false
		for i := searchFrom; i < len(tRunes); i++ {
			if tRunes[i] == qc {
				if lastIdx >= 0 {
					gap := i - lastIdx - 1
					score += gap * 2
					if gap == 0 {
						score -= 1
					}
				} else {
					score += i * 2
				}
				lastIdx = i
				searchFrom = i + 1
				found = true
				break
			}
		}
		if !found {
			return 0, false
		}
	}

	return score, true
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
