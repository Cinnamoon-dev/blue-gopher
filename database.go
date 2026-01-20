package main

import (
	"database/sql"
	"fmt"
	"log"
)

func createTables(db *sql.DB) {
	userTable := "CREATE TABLE IF NOT EXISTS usuarios (id INTEGER PRIMARY KEY AUTOINCREMENT, nome TEXT, idade INTEGER)"
	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllUsers(db *sql.DB) {
	rows, err := db.Query("SELECT * from usuarios")
	if err != nil {
		log.Fatal(err)
	}

	var user struct {
		ID    int
		Nome  string
		Idade int
	}

	for rows.Next() {
		rows.Scan(&user.ID, &user.Nome, &user.Idade)
		fmt.Printf("%+v\n", user)
	}
}

func database() {
	db, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables(db)
	stmt, err := db.Prepare("INSERT INTO usuarios(nome, idade) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := stmt.Exec("pedro", 15); err != nil {
		log.Fatal(err)
	}

	if _, err := stmt.Exec("daniel", 25); err != nil {
		log.Fatal(err)
	}

	if _, err := stmt.Exec("guilherme", 21); err != nil {
		log.Fatal(err)
	}

	getAllUsers(db)
}
