package database

import (
	"bufio"
	"database/sql"
	"os"
	"strings"
)

// Takes a path to a sql file and divides all the file in strings with ";" as the separator
func CreateTables(path string, db *sql.DB) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	sql := string(file)
	queries := strings.SplitSeq(sql, ";")

	for query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}

// Take a path to a sql file and reads each line of the file (ignoring lines that start with "--") and execute each line
func Populate(path string, db *sql.DB) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "--") {
			continue
		}

		db.Exec(line)
	}
}
