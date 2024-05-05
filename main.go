package main

import (
	. "phabricator/bugz"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
db, err := sql.Open("sqlite3", "bugs.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS bugs (
	id INTEGER PRIMARY KEY,
	CreationTime        TEXT,
	Creator             TEXT,
	Summary             TEXT
);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	err = importBugsFromJSON(db, "bugsJson")
	if err != nil {
		log.Fatalf("Error importing bugs from JSON: %v", err)
	}

	fmt.Println("Database and schema are ready. Bugs imported successfully.")
}

func importBugsFromJSON(db *sql.DB, directory string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		bug := Bug{}
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", filePath, err)
		}

		err = json.Unmarshal(fileData, &bug)
		if err != nil {
			return fmt.Errorf("error decoding JSON from file %s: %v", filePath, err)
		}

		_, err = db.Exec("INSERT INTO bugs (ID, CreationTime, Creator, Summary) VALUES (?, ?, ?, ?)", bug.ID, bug.CreationTime, bug.Creator, bug.Summary)
		if err != nil {
			return fmt.Errorf("error inserting bug into database: %v", err)
		}
	}

	return nil
}
