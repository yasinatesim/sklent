package enhance

import (
	"sort"
	"strings"
)

type CategoryOption struct {
	ID   int64
	Name string
}

type scored struct {
	option CategoryOption
	score  int
}

// FilterCandidates ranks options by title-token hits so the LLM only ever sees real category ids.
func FilterCandidates(title string, options []CategoryOption, limit int) []CategoryOption {
	tokens := significantTokens(title)
	if len(tokens) == 0 || len(options) == 0 {
		return nil
	}
	ranked := make([]scored, 0, len(options))
	for _, o := range options {
		name := strings.ToLower(o.Name)
		s := 0
		for _, t := range tokens {
			if strings.Contains(name, t) {
				s++
			}
		}
		if s > 0 {
			ranked = append(ranked, scored{o, s})
		}
	}
	if len(ranked) == 0 {
		return nil
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}
	out := make([]CategoryOption, len(ranked))
	for i, r := range ranked {
		out[i] = r.option
	}
	return out
}

func significantTokens(s string) []string {
	fields := strings.Fields(strings.ToLower(s))
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.Trim(f, ".,()[]{}!?\"'")
		if len(f) >= 3 {
			out = append(out, f)
		}
	}
	return out
}
