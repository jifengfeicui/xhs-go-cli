package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/querygen"
	"xhs-go-cli/internal/source"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("xhs-go-cli: use subcommands import-sources/query-gen/search/fetch-detail/qualify")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "import-sources":
		runImportSources(os.Args[2:])
	case "query-gen":
		runQueryGen(os.Args[2:])
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runImportSources(args []string) {
	fs := flag.NewFlagSet("import-sources", flag.ExitOnError)
	dbPath := fs.String("db", "xhs.db", "sqlite db path")
	input := fs.String("input", "", "source json path")
	_ = fs.Parse(args)
	if *input == "" {
		fmt.Println("--input is required")
		os.Exit(1)
	}
	database, err := db.Open(*dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	repo := source.NewRepo(database)
	count, err := source.ImportFromJSON(repo, *input)
	if err != nil {
		panic(err)
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"imported": count})
}

func runQueryGen(args []string) {
	fs := flag.NewFlagSet("query-gen", flag.ExitOnError)
	dbPath := fs.String("db", "xhs.db", "sqlite db path")
	limit := fs.Int("limit", 10, "source limit")
	perSource := fs.Int("per-source", 3, "queries per source")
	_ = fs.Parse(args)
	_ = perSource

	database, err := db.Open(*dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	repo := source.NewRepo(database)
	sources, err := repo.List(*limit)
	if err != nil {
		panic(err)
	}

	result := make([]map[string]any, 0, len(sources))
	for _, src := range sources {
		qsrc := querygen.Source{
			ID:         src.ID,
			Name:       src.Name,
			Keywords:   src.Keywords,
			SourceType: querygen.ClassifySourceType(src.Name, src.Keywords),
		}
		queries := querygen.GenerateQueries(qsrc, *perSource)
		result = append(result, map[string]any{
			"source_id":   src.ID,
			"source_name": src.Name,
			"source_type": qsrc.SourceType,
			"queries":     queries,
		})
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
}
