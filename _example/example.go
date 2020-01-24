package main

import (
	"flag"
	"fmt"

	"github.com/zncoder/easycmd"
)

var db string

func defineDBFlags() {
	flag.StringVar(&db, "db", "", "path to the DB file")
}

func runDBCreate() {
	copies := flag.Int("copies", 3, "number of versions to keep")
	flag.Parse()

	fmt.Println("create db", db, *copies)
}

func runDBQuery() {
	key := flag.String("key", "", "key to query")
	last := flag.Bool("last", false, "query only the last version")
	flag.Parse()

	fmt.Println("query db", db, *key, *last)
}

func main() {
	easycmd.Handle("db", defineDBFlags, "commands to operate a DB")
	easycmd.Handle("db create", runDBCreate, "create a db")
	easycmd.Handle("db query", runDBQuery, "query a db")
	easycmd.Main()
}
