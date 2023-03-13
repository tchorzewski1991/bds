package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/olekukonko/tablewriter"
	"github.com/tchorzewski1991/bds/business/sys/database"
)

var build = "develop"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	cfg := struct {
		conf.Version
		DB struct {
			User string `conf:"default:postgres"`
			Pass string `conf:"default:password,mask"`
			Host string `conf:"default:db"`
			Name string `conf:"default:bds"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Current build version",
		},
	}

	const prefix = "BOOKS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config failed: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config failed: %w", err)
	}
	fmt.Println(out)

	db, err := database.Open(database.Config{
		User: cfg.DB.User,
		Pass: cfg.DB.Pass,
		Host: cfg.DB.Host,
		Name: cfg.DB.Name,
	})
	if err != nil {
		return fmt.Errorf("opening database failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	schema, err := database.GenerateSchema(ctx, db)
	if err != nil {
		return fmt.Errorf("generating database schema failed: %w", err)
	}

	fmt.Println("Database schema generated successfully.")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"table_name", "column_name", "data_type"})
	table.SetAutoMergeCellsByColumnIndex([]int{0})

	for _, entry := range schema.Entries {
		table.Append([]string{entry.TableName, entry.ColumnName, entry.DataType})
	}

	table.Render()
	return nil
}
