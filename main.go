package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func getConn(name string) *sql.DB {
	os.Mkdir("dbs", 0755)
	db, err := sql.Open("sqlite3", fmt.Sprintf("dbs/%s", name))
	if err != nil {
		log.Fatal(err)
	}

	query1 := "CREATE TABLE IF NOT EXISTS main (x TEXT PRIMARY KEY, y TEXT)"
	_, err = db.Exec(query1)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func set(db *sql.DB, key string, val string) {
	queryCheck := "SELECT x FROM main WHERE x = ?"
	var existingKey string
	err := db.QueryRow(queryCheck, key).Scan(&existingKey)

	if err == sql.ErrNoRows {
		queryInsert := "INSERT INTO main (x, y) VALUES (?, ?)"
		_, err = db.Exec(queryInsert, key, val)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		queryUpdate := "UPDATE main SET y = ? WHERE x = ?"
		_, err = db.Exec(queryUpdate, val, key)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func get(db *sql.DB, key string) (string, bool) {
	query := "SELECT y FROM main WHERE x = ?"
	var value string
	err := db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", false
	} else if err != nil {
		log.Fatal(err)
	}
	return value, true
}

func main() {
	db := getConn("test.db")
	defer db.Close()

	set(db, "key1", "value1")
	value, exists := get(db, "key1")
	if exists {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key not found")
	}
}
