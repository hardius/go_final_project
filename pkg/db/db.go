package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date VARCHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "", 
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler(date);

`

func Init(dbFile string) (*sql.DB, error) {
	if len(dbFile) == 0 {
		dbFile = "scheduler.db"
	}

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
		err = nil
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return db, err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return db, err
		}
	}

	return db, nil
}
