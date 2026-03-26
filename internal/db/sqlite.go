package db

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS sources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			keywords TEXT NOT NULL DEFAULT '',
			source_type TEXT NOT NULL DEFAULT '',
			priority TEXT NOT NULL DEFAULT '',
			level TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS generated_queries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_id INTEGER NOT NULL,
			query TEXT NOT NULL,
			query_type TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS search_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			query_id INTEGER,
			feed_id TEXT NOT NULL,
			xsec_token TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL DEFAULT '',
			author TEXT NOT NULL DEFAULT '',
			raw_json TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS details (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			feed_id TEXT NOT NULL,
			xsec_token TEXT NOT NULL DEFAULT '',
			detail_json TEXT NOT NULL DEFAULT '',
			fetch_status TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS qualifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			feed_id TEXT NOT NULL,
			status TEXT NOT NULL,
			title TEXT NOT NULL DEFAULT '',
			source_link TEXT NOT NULL DEFAULT '',
			claim_rule TEXT NOT NULL DEFAULT '',
			location TEXT NOT NULL DEFAULT '',
			participation_method TEXT NOT NULL DEFAULT '',
			reason TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
