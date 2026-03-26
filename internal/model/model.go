package model

import "time"

type Source struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Keywords       string    `gorm:"size:255;not null;default:''" json:"keywords"`
	SourceType     string    `gorm:"size:50;not null;default:''" json:"source_type"`
	Priority       int       `gorm:"not null;default:0" json:"priority"`
	Level          string    `gorm:"size:50;not null;default:''" json:"level"`
	City           string    `gorm:"size:50;not null;default:''" json:"city"`
	Confidence     string    `gorm:"size:20;not null;default:''" json:"confidence"`
	Reason         string    `gorm:"size:500;not null;default:''" json:"reason"`
	PreferenceType string    `gorm:"size:255;not null;default:''" json:"preference_type"`
	Remark         string    `gorm:"type:text;not null;default:''" json:"remark"`
	LastStatus     string    `gorm:"size:50;not null;default:''" json:"last_status"`
	QueryCount     int       `gorm:"not null;default:0" json:"query_count"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Source) TableName() string {
	return "sources"
}

type GeneratedQuery struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SourceID  uint      `gorm:"not null;index" json:"source_id"`
	Query     string    `gorm:"type:text;not null" json:"query"`
	QueryType string    `gorm:"size:50;not null;default:''" json:"query_type"`
	Status    string    `gorm:"size:20;not null;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`

	Source Source `gorm:"foreignKey:SourceID"`
}

func (GeneratedQuery) TableName() string {
	return "generated_queries"
}

type SearchResult struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	QueryID   uint      `gorm:"not null;index" json:"query_id"`
	FeedID    string    `gorm:"size:100;not null" json:"feed_id"`
	XsecToken string    `gorm:"size:255;not null;default:''" json:"xsec_token"`
	Title     string    `gorm:"size:500;not null;default:''" json:"title"`
	Author    string    `gorm:"size:255;not null;default:''" json:"author"`
	RawJSON   string    `gorm:"type:text;not null" json:"raw_json"`
	Status    string    `gorm:"size:20;not null;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (SearchResult) TableName() string {
	return "search_results"
}

type Detail struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FeedID      string    `gorm:"size:100;not null" json:"feed_id"`
	XsecToken   string    `gorm:"size:255;not null;default:''" json:"xsec_token"`
	DetailJSON  string    `gorm:"type:text;not null" json:"detail_json"`
	FetchStatus string    `gorm:"size:50;not null;default:''" json:"fetch_status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Detail) TableName() string {
	return "details"
}

type Qualification struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	FeedID              string    `gorm:"size:100;not null" json:"feed_id"`
	Status              string    `gorm:"size:50;not null" json:"status"`
	Title               string    `gorm:"size:500;not null;default:''" json:"title"`
	SourceLink          string    `gorm:"size:500;not null;default:''" json:"source_link"`
	ClaimRule           string    `gorm:"size:1000;not null;default:''" json:"claim_rule"`
	Location            string    `gorm:"size:255;not null;default:''" json:"location"`
	ParticipationMethod string    `gorm:"size:255;not null;default:''" json:"participation_method"`
	Reason              string    `gorm:"size:1000;not null;default:''" json:"reason"`
	CreatedAt           time.Time `json:"created_at"`
}

func (Qualification) TableName() string {
	return "qualifications"
}
