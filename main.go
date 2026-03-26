package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/detail"
	"xhs-go-cli/internal/mcp"
	"xhs-go-cli/internal/qualify"
	"xhs-go-cli/internal/querygen"
	"xhs-go-cli/internal/search"
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
	case "search":
		runSearch(os.Args[2:])
	case "fetch-detail":
		runFetchDetail(os.Args[2:])
	case "qualify":
		runQualify(os.Args[2:])
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

	database, err := db.Open(*dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	repo := source.NewRepo(database)
	searchSvc := search.NewService(database, nil)
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
		for _, query := range queries {
			_ = searchSvc.SaveGeneratedQuery(src.ID, query, qsrc.SourceType)
		}
		result = append(result, map[string]any{
			"source_id":   src.ID,
			"source_name": src.Name,
			"source_type": qsrc.SourceType,
			"queries":     queries,
		})
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
}

func runSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	dbPath := fs.String("db", "xhs.db", "sqlite db path")
	limit := fs.Int("limit", 5, "query limit")
	pageSize := fs.Int("page-size", 10, "search page size")
	baseURL := fs.String("base-url", "http://127.0.0.1:18060", "mcp base url")
	_ = fs.Parse(args)

	database, err := db.Open(*dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	client := mcp.New(*baseURL)
	service := search.NewService(database, client)
	queries, err := service.ListQueries(*limit)
	if err != nil {
		panic(err)
	}
	out := make([]map[string]any, 0, len(queries))
	for _, q := range queries {
		count, err := service.SearchAndStore(q.ID, q.Query, *pageSize)
		if err != nil {
			out = append(out, map[string]any{"query_id": q.ID, "query": q.Query, "error": err.Error()})
			continue
		}
		out = append(out, map[string]any{"query_id": q.ID, "query": q.Query, "stored": count})
	}
	_ = json.NewEncoder(os.Stdout).Encode(out)
}

func runFetchDetail(args []string) {
	fs := flag.NewFlagSet("fetch-detail", flag.ExitOnError)
	dbPath := fs.String("db", "xhs.db", "sqlite db path")
	limit := fs.Int("limit", 20, "detail row limit")
	concurrency := fs.Int("concurrency", 3, "detail fetch concurrency")
	baseURL := fs.String("base-url", "http://127.0.0.1:18060", "mcp base url")
	_ = fs.Parse(args)

	database, err := db.Open(*dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	client := mcp.New(*baseURL)
	service := detail.NewService(database, client)
	rows, err := service.ListPending(*limit)
	if err != nil {
		panic(err)
	}
	result, err := service.FetchAndStore(rows, *concurrency)
	if err != nil {
		panic(err)
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
}

func runQualify(args []string) {
	fs := flag.NewFlagSet("qualify", flag.ExitOnError)
	dbPath := fs.String("db", "xhs.db", "sqlite db path")
	limit := fs.Int("limit", 20, "qualification row limit")
	_ = fs.Parse(args)

	database, err := db.Open(*dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	service := qualify.NewService(database)
	rows, err := service.ListDetails(*limit)
	if err != nil {
		panic(err)
	}
	result, err := service.QualifyAndStore(rows)
	if err != nil {
		panic(err)
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
}
