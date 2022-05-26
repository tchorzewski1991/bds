package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/business/sys/database"
	"io"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// private

func run() error {

	fmt.Println("Loading books")
	start := time.Now()

	var source string
	flag.StringVar(&source, "source", "", "Source file with books to be parsed and inserted into db")

	flag.Parse()

	if source == "" {
		fmt.Println("ERR: source cannot be empty")
		os.Exit(1)
	}

	f, err := os.Open(source)
	if err != nil {
		fmt.Println("ERR: source file does not exist")
		os.Exit(1)
	}

	db, err := database.Open(database.Config{
		User: "postgres",
		Pass: "password",
		Host: "localhost",
		Name: "bds",
	})
	if err != nil {
		fmt.Printf("ERR: cannot connect to db: %v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = database.Check(ctx, db)
	if err != nil {
		fmt.Printf("ERR: checking db status failes: %v\n", err)
		os.Exit(1)
	}

	r := csv.NewReader(f)

	_, err = r.Read()
	if err != nil {
		fmt.Println("ERR: cannot read source file")
		os.Exit(1)
	}

	var entry []string

	stats := struct {
		total   int
		success int
		fail    int
	}{0, 0, 0}

	for {
		stats.total += 1

		entry, err = r.Read()
		if err == io.EOF {
			stats.total -= 1
			break
		}
		if err != nil {
			stats.fail += 1
			continue
		}

		err = save(db, entry)
		if err != nil {
			stats.fail += 1
			continue
		}

		stats.success += 1
	}

	end := time.Since(start)

	fmt.Printf("Books loaded. Stats: %+v | Took: %v\n", stats, end)

	return nil
}

func save(db *sqlx.DB, entry []string) error {
	const q = `
		insert into books
     		(isbn, title, author, publication_year, publisher)
		values
		    (:isbn, :title, :author, :publication_year, :publisher)
	`

	data := struct {
		Isbn            string `db:"isbn"`
		Title           string `db:"title"`
		Author          string `db:"author"`
		PublicationYear string `db:"publication_year"`
		Publisher       string `db:"publisher"`
	}{
		Isbn:            entry[0],
		Title:           entry[1],
		Author:          entry[2],
		PublicationYear: entry[3],
		Publisher:       entry[4],
	}

	_, err := db.NamedExec(q, data)
	if err != nil {
		return err
	}

	return nil
}
