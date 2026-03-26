package qualify

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type Service struct {
	db *sql.DB
}

type DetailRow struct {
	ID         int64
	FeedID     string
	DetailJSON string
	Status     string
}

type Qualification struct {
	FeedID             string `json:"feed_id"`
	Status             string `json:"status"`
	Title              string `json:"title,omitempty"`
	SourceLink         string `json:"source_link,omitempty"`
	ClaimRule          string `json:"claim_rule,omitempty"`
	Location           string `json:"location,omitempty"`
	ParticipationMethod string `json:"participation_method,omitempty"`
	Reason             string `json:"reason,omitempty"`
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListDetails(limit int) ([]DetailRow, error) {
	rows, err := s.db.Query(`SELECT id, feed_id, detail_json, fetch_status FROM details ORDER BY id ASC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DetailRow
	for rows.Next() {
		var row DetailRow
		if err := rows.Scan(&row.ID, &row.FeedID, &row.DetailJSON, &row.Status); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *Service) QualifyAndStore(rows []DetailRow) ([]Qualification, error) {
	out := make([]Qualification, 0, len(rows))
	for _, row := range rows {
		q := qualifyOne(row)
		_, err := s.db.Exec(`INSERT INTO qualifications(feed_id, status, title, source_link, claim_rule, location, participation_method, reason) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
			q.FeedID, q.Status, q.Title, q.SourceLink, q.ClaimRule, q.Location, q.ParticipationMethod, q.Reason,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, nil
}

func qualifyOne(row DetailRow) Qualification {
	payload := map[string]any{}
	_ = json.Unmarshal([]byte(row.DetailJSON), &payload)
	data := nestedMap(payload, "data", "data")
	desc := asString(data["desc"])
	title := asString(data["title"])
	sourceLink := asString(data["note_url"])
	location := firstMatched(desc, []string{"上海", "静安", "徐汇", "新天地", "前滩", "环贸", "南京东路"})
	claimRule := firstSentenceContaining(desc, []string{"领", "赠", "打卡", "到店", "预约", "免费"})
	participation := firstSentenceContaining(desc, []string{"到店", "打卡", "预约", "参与", "领取"})

	q := Qualification{FeedID: row.FeedID, Title: title, SourceLink: sourceLink, ClaimRule: claimRule, Location: location, ParticipationMethod: participation}
	missing := []string{}
	if q.Title == "" { missing = append(missing, "title") }
	if q.SourceLink == "" { missing = append(missing, "source_link") }
	if q.ClaimRule == "" { missing = append(missing, "claim_rule") }
	if q.Location == "" { missing = append(missing, "location") }
	if q.ParticipationMethod == "" { missing = append(missing, "participation_method") }
	if len(missing) == 0 {
		q.Status = "accepted"
		return q
	}
	q.Status = "rejected"
	q.Reason = "missing: " + strings.Join(missing, ",")
	return q
}

func nestedMap(value map[string]any, keys ...string) map[string]any {
	current := value
	for _, key := range keys {
		next, _ := current[key].(map[string]any)
		if next == nil {
			return map[string]any{}
		}
		current = next
	}
	return current
}

func asString(v any) string {
	if s, ok := v.(string); ok { return s }
	return ""
}

func firstMatched(text string, options []string) string {
	for _, option := range options {
		if strings.Contains(text, option) {
			return option
		}
	}
	return ""
}

func firstSentenceContaining(text string, markers []string) string {
	parts := strings.FieldsFunc(text, func(r rune) bool { return r == '\n' || r == '。' || r == '！' || r == '!' })
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		for _, marker := range markers {
			if strings.Contains(part, marker) {
				return part
			}
		}
	}
	return ""
}
