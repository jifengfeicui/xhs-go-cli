package source

import (
	"database/sql"
)

type Source struct {
	ID         int64
	Name       string
	Keywords   string
	SourceType string
	Priority   string
	Level      string
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Insert(src Source) (int64, error) {
	result, err := r.db.Exec(
		`INSERT INTO sources(name, keywords, source_type, priority, level) VALUES(?, ?, ?, ?, ?)`,
		src.Name, src.Keywords, src.SourceType, src.Priority, src.Level,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repo) List(limit int) ([]Source, error) {
	query := `SELECT id, name, keywords, source_type, priority, level FROM sources ORDER BY id ASC`
	var rows *sql.Rows
	var err error
	if limit > 0 {
		query += ` LIMIT ?`
		rows, err = r.db.Query(query, limit)
	} else {
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Source
	for rows.Next() {
		var src Source
		if err := rows.Scan(&src.ID, &src.Name, &src.Keywords, &src.SourceType, &src.Priority, &src.Level); err != nil {
			return nil, err
		}
		out = append(out, src)
	}
	return out, rows.Err()
}
