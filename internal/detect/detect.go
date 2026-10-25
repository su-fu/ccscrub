package detect

import (
	"cmp"
	"regexp"
	"slices"
)

// previewKeep and previewMinLen control how Preview masks a match.
const (
	previewKeep   = 4
	previewMinLen = 12
)

type Rule struct {
	ID      string
	Pattern *regexp.Regexp
}

var rules = []Rule{
	// Format: https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/
	{ID: "github-token", Pattern: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36}`)},
	{ID: "github-token", Pattern: regexp.MustCompile(`github_pat_[A-Za-z0-9]{22}_[A-Za-z0-9]{59}`)},
	// Format: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
	{ID: "aws-access-key-id", Pattern: regexp.MustCompile(`A[KS]IA[A-Z0-9]{16}`)},
	// Based on gitleaks (MIT): https://github.com/gitleaks/gitleaks
	{ID: "anthropic-api-key", Pattern: regexp.MustCompile(`sk-ant-api03-[a-zA-Z0-9_\-]{93}AA`)},
	{ID: "anthropic-admin-api-key", Pattern: regexp.MustCompile(`sk-ant-admin01-[a-zA-Z0-9_\-]{93}AA`)},
	// Based on gitleaks (MIT): https://github.com/gitleaks/gitleaks
	{ID: "openai-api-key", Pattern: regexp.MustCompile(`sk-(?:proj|svcacct|admin)-(?:[a-zA-Z0-9_-]{74}|[a-zA-Z0-9_-]{58})T3BlbkFJ(?:[a-zA-Z0-9_-]{74}|[a-zA-Z0-9_-]{58})`)},
	{ID: "openai-api-key", Pattern: regexp.MustCompile(`sk-[a-zA-Z0-9]{20}T3BlbkFJ[a-zA-Z0-9]{20}`)},
}

type Match struct {
	RuleID string
	Start  int
	End    int
}

// Find looks for secrets in line and returns where they are,
// one Match for each.
//
// A match counts only if it is not part of a longer word.
// The matches are sorted by Start. If nothing is found, Find returns nil.
func Find(line []byte) []Match {
	var result []Match
	for _, rule := range rules {
		locs := rule.Pattern.FindAllIndex(line, -1)
		for _, loc := range locs {
			if !okBefore(line, loc[0]) || !okAfter(line, loc[1]) {
				continue
			}
			result = append(result, Match{RuleID: rule.ID, Start: loc[0], End: loc[1]})
		}
	}
	slices.SortFunc(result, func(a, b Match) int {
		return cmp.Compare(a.Start, b.Start)
	})
	return result
}

// Preview returns a masked form of the match m in line.
//
// It shows only the first and last previewKeep bytes, with "****" between them.
// If the match is previewMinLen bytes or shorter, it returns only "****".
// It never returns the full matched text.
func Preview(line []byte, m Match) string {
	s := line[m.Start:m.End]
	if len(s) <= previewMinLen {
		return "****"
	}
	return string(s[:previewKeep]) + "****" + string(s[len(s)-previewKeep:])
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || (b == '_')
}

func okBefore(line []byte, start int) bool {
	if start == 0 {
		return true
	}
	if (start >= 2) && (line[start-2] == '\\') && (line[start-1] == 'n' || line[start-1] == 't' || line[start-1] == 'r') {
		return true
	}
	return !isWordChar(line[start-1])
}

func okAfter(line []byte, end int) bool {
	if end == len(line) {
		return true
	}
	return !isWordChar(line[end])
}
