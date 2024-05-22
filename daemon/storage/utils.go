package storage

import "database/sql"

//storage/utils.go

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		size INTEGER,
		upload_at DATETIME,
		file_path TEXT
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
