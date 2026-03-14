package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
	"urlshort/internal/config"
	"urlshort/internal/storage/dragonfly"
)

func main() {
	cfg := config.MustLoad()
	storage := dragonfly.New(cfg.Dragonfly)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var aliasesForDelete []string
	flag.Func("D", "URL alias to delete (can be specified multiple times)", func(s string) error {
		aliasesForDelete = append(aliasesForDelete, s)
		return nil
	})

	flag.Parse()

	if len(aliasesForDelete) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no URL aliases specified")
		fmt.Fprintln(os.Stderr, "Usage: delete-url -D <alias> [-D <alias> ...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var deleted, failed []string

	for _, alias := range aliasesForDelete {
		err := storage.DeleteURL(ctx, alias)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete %q: %v\n", alias, err)
			failed = append(failed, alias)
		} else {
			fmt.Printf("Deleted: %s\n", alias)
			deleted = append(deleted, alias)
		}
	}

	fmt.Printf("\nSummary: %d deleted, %d failed\n", len(deleted), len(failed))

	if len(failed) > 0 {
		os.Exit(1)
	}
}
