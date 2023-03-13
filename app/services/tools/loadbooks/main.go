package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/business/sys/database"
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
	flag.StringVar(&source, "source", "", "The source file with books to be inserted into db.")

	var bufferSize int
	flag.IntVar(&bufferSize, "buffer", 10_000, "The size of books buffer used within single db tx.")

	flag.Parse()

	if source == "" {
		return errors.New("source cannot be empty")
	}

	if bufferSize < 2 || bufferSize > 50_000 {
		return errors.New("buffer size is not valid")
	}

	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("source file does not exist: %w", err)
	}

	db, err := database.Open(database.Config{
		User: "postgres",
		Pass: "password",
		Host: "db",
		Name: "bds",
	})
	if err != nil {
		return fmt.Errorf("cannot connect to db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = database.Check(ctx, db)
	if err != nil {
		return fmt.Errorf("cannot check db status: %w", err)
	}

	r := csv.NewReader(f)

	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("cannot read source: %w", err)
	}

	stats := struct {
		total   int
		success int
		failure int
		retries int
	}{0, 0, 0, 0}

	bufferPos := 0
	buffer := make([][]string, 0, bufferSize)

	var tx *sqlx.Tx

	releaseBuffer := func() (success, failure int) {
		tx, err = db.Beginx()
		if err != nil {
			return success, failure
		}

		for idx := range buffer {
			err = save(tx, buffer[idx])
			if err != nil {
				failure += 1
				continue
			}
			success += 1
		}

		// TODO: revisit releaseBuffer impl. and handle err properly
		_ = tx.Commit()

		return success, failure
	}

	var row []string

	for {
		if len(buffer) < bufferSize {
			row, err = r.Read()
			if errors.Is(err, io.EOF) {
				success, failure := releaseBuffer()
				stats.success += success
				stats.failure += failure
				stats.total += success + failure
				break
			}
			if err != nil {
				stats.retries += 1
				continue
			}

			buffer = append(buffer, row)
			bufferPos += 1
			continue
		}

		success, failure := releaseBuffer()

		stats.success += success
		stats.failure += failure
		stats.total += success + failure

		buffer = make([][]string, 0, bufferSize)
		bufferPos = 0

		if stats.total%10_000 == 0 {
			fmt.Printf("Books loaded. Stats: %+v\n", stats)
		}
	}

	end := time.Since(start)
	fmt.Printf("Books loaded. Stats: %+v | Took: %v\n", stats, end)

	return nil
}

type Book struct {
	Isbn            string `db:"isbn"`
	Title           string `db:"title"`
	Author          string `db:"author"`
	PublicationYear string `db:"publication_year"`
	Publisher       string `db:"publisher"`
}

func save(tx *sqlx.Tx, entry []string) error {
	const q = `
		insert into books
     		(isbn, title, author, publication_year, publisher)
		values
		    (:isbn, :title, :author, :publication_year, :publisher)
	`

	data := Book{
		Isbn:            entry[0],
		Title:           entry[1],
		Author:          entry[2],
		PublicationYear: entry[3],
		Publisher:       entry[4],
	}

	_, err := tx.NamedExec(q, data)
	if err != nil {
		return err
	}

	return nil
}
