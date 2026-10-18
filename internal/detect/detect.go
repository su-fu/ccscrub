package detect

import (
	"regexp"
)

type Rule struct {
	ID      string
	Pattern *regexp.Regexp
}

var rules = []Rule{
	{ID: "github-token", Pattern: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36}`)},
}

type Match struct {
	RuleID string
	Start  int
	End    int
}

func Find(line []byte) []Match {
	var result []Match
	for _, rule := range rules {
		locs := rule.Pattern.FindAllIndex(line, -1)
		for _, loc := range locs {
			result = append(result, Match{RuleID: rule.ID, Start: loc[0], End: loc[1]})
		}
	}
	return result
}
