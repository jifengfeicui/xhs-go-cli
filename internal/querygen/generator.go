package querygen

import (
	"strings"
)

var sourceTypeKeywords = map[string][]string{
	"mall":           {"商场", "广场", "天地", "城", "中心", "百联", "万象", "大悦城", "iapm", "太古里", "k11"},
	"brand":          {"品牌", "美妆", "护肤", "香水", "新品", "试用", "lancome", "cpb", "kiehl", "clarins"},
	"official_event": {"官方活动", "活动ip", "巡展", "周年", "展览", "快闪ip", "联名"},
	"info_account":   {"情报", "情报号", "探店", "汇总", "攻略", "线报", "来源"},
}

var fallbackSuffixes = map[string][]string{
	"mall":           {"快闪", "打卡有礼", "到店礼"},
	"brand":          {"赠礼", "到店礼", "试用"},
	"official_event": {"上海", "快闪", "活动"},
	"info_account":   {"上海", "活动", "免费领"},
	"generic":        {"免费领", "快闪", "打卡有礼"},
}

func ClassifySourceType(sourceName, keywords string) string {
	text := strings.ToLower(sourceName + " " + keywords)
	scores := map[string]int{}
	for sourceType, markers := range sourceTypeKeywords {
		for _, marker := range markers {
			if strings.Contains(text, strings.ToLower(marker)) {
				scores[sourceType]++
			}
		}
	}
	bestType := "generic"
	bestScore := 0
	for sourceType, score := range scores {
		if score > bestScore {
			bestScore = score
			bestType = sourceType
		}
	}
	return bestType
}

func GenerateQueries(source Source, limit int) []string {
	keywords := parseKeywords(source.Keywords)
	sourceType := source.SourceType
	if sourceType == "" {
		sourceType = ClassifySourceType(source.Name, source.Keywords)
	}

	queries := make([]string, 0, limit)
	seen := map[string]bool{}
	for _, kw := range keywords {
		if isRedundantKeyword(source.Name, kw) {
			continue
		}
		query := strings.TrimSpace(source.Name + " " + kw)
		if query == "" || seen[query] {
			continue
		}
		seen[query] = true
		queries = append(queries, query)
		if len(queries) >= limit {
			return queries
		}
	}
	for _, suffix := range fallbackSuffixes[sourceType] {
		query := strings.TrimSpace(source.Name + " " + suffix)
		if query == "" || seen[query] {
			continue
		}
		seen[query] = true
		queries = append(queries, query)
		if len(queries) >= limit {
			return queries
		}
	}
	return queries
}

func parseKeywords(value string) []string {
	replacer := strings.NewReplacer("，", ",", "\n", ",")
	parts := strings.Split(replacer.Replace(value), ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func isRedundantKeyword(sourceName, keyword string) bool {
	sourceNorm := normalizeToken(sourceName)
	keywordNorm := normalizeToken(keyword)
	if keywordNorm == "" {
		return true
	}
	if keywordNorm == sourceNorm || strings.Contains(sourceNorm, keywordNorm) {
		return true
	}
	return false
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), ""))
}
