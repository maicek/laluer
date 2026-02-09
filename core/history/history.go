package history

import (
	"database/sql"
	"os"
	"path/filepath"

	"time"

	_ "modernc.org/sqlite"
)

type EntryTypes string

const (
	ENTRY_TYPE_APP EntryTypes = "app"
)

type HistoryEntry struct {
	Type      EntryTypes
	EntryName string
	LastUsed  int64
	UseCount  int
}

type HistoryService struct {
	DB *sql.DB
}

var Service *HistoryService

func Init() (*sql.DB, error) {

	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		stateDir = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	historyPath := filepath.Join(stateDir, "laluer", "history.db")

	if err := os.MkdirAll(filepath.Dir(historyPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", historyPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			name TEXT NOT NULL,
			last_used INTEGER NOT NULL,
			use_count INTEGER NOT NULL DEFAULT 1,
			UNIQUE(type, name)
		);

		CREATE INDEX IF NOT EXISTS idx_history_name ON history(name);
	`)

	Service = &HistoryService{DB: db}

	return db, err
}

func (s *HistoryService) Push(t EntryTypes, name string) error {
	now := time.Now().Unix()

	_, err := s.DB.Exec(`
		INSERT INTO history (type, name, last_used, use_count)
		VALUES (?, ?, ?, 1)
		ON CONFLICT(type, name)
		DO UPDATE SET
			last_used = ?,
			use_count = use_count + 1
	`, t, name, now, now)

	return err
}

func (s *HistoryService) GetLast() ([]HistoryEntry, error) {
	var entries []HistoryEntry

	rows, err := s.DB.Query(`
		SELECT type, name, last_used, use_count FROM history
		ORDER BY last_used DESC
		LIMIT 10
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry HistoryEntry

		if err := rows.Scan(&entry.Type, &entry.EntryName, &entry.LastUsed, &entry.UseCount); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// func (s *HistoryService) Get(t EntryTypes, name string) (*HistoryEntry, error) {
// 	var entry HistoryEntry

// 	return &entry, nil
// }
