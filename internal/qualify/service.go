package qualify

import (
	"context"
	"encoding/json"
	"strings"

	"xhs-go-cli/internal/model"
	"xhs-go-cli/internal/repository"
)

type Service struct {
	detailRepo repository.DetailRepository
	qualRepo   repository.QualificationRepository
}

type DetailRow struct {
	ID         uint
	FeedID     string
	DetailJSON string
	Status     string
}

type Qualification struct {
	FeedID              string `json:"feed_id"`
	Status              string `json:"status"`
	Title               string `json:"title,omitempty"`
	SourceLink          string `json:"source_link,omitempty"`
	ClaimRule           string `json:"claim_rule,omitempty"`
	Location            string `json:"location,omitempty"`
	ParticipationMethod string `json:"participation_method,omitempty"`
	Reason              string `json:"reason,omitempty"`
}

func NewService(detailRepo repository.DetailRepository, qualRepo repository.QualificationRepository) *Service {
	return &Service{
		detailRepo: detailRepo,
		qualRepo:   qualRepo,
	}
}

func (s *Service) ListDetails(ctx context.Context, limit int) ([]DetailRow, error) {
	details, err := s.detailRepo.ListPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	rows := make([]DetailRow, len(details))
	for i, d := range details {
		rows[i] = DetailRow{
			ID:         d.ID,
			FeedID:     d.FeedID,
			DetailJSON: d.DetailJSON,
			Status:     d.FetchStatus,
		}
	}
	return rows, nil
}

func (s *Service) QualifyAndStore(ctx context.Context, rows []DetailRow) ([]Qualification, error) {
	out := make([]Qualification, 0, len(rows))
	for _, row := range rows {
		q := qualifyOne(row)
		err := s.qualRepo.Create(ctx, &model.Qualification{
			FeedID:              q.FeedID,
			Status:              q.Status,
			Title:               q.Title,
			SourceLink:          q.SourceLink,
			ClaimRule:           q.ClaimRule,
			Location:            q.Location,
			ParticipationMethod: q.ParticipationMethod,
			Reason:              q.Reason,
		})
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
	if q.Title == "" {
		missing = append(missing, "title")
	}
	if q.SourceLink == "" {
		missing = append(missing, "source_link")
	}
	if q.ClaimRule == "" {
		missing = append(missing, "claim_rule")
	}
	if q.Location == "" {
		missing = append(missing, "location")
	}
	if q.ParticipationMethod == "" {
		missing = append(missing, "participation_method")
	}
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
	if s, ok := v.(string); ok {
		return s
	}
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
