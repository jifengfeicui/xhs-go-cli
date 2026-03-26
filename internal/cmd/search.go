package cmd

import (
	"encoding/json"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/mcp"
	"xhs-go-cli/internal/repository"
	"xhs-go-cli/internal/search"

	"github.com/spf13/cobra"
)

func NewSearchCmd() *cobra.Command {
	var limit int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "搜索内容并存储结果",
		RunE: func(cmd *cobra.Command, args []string) error {
			database, err := db.Open(getDBPath())
			if err != nil {
				return err
			}
			defer database.Close()

			queryRepo := repository.NewQueryRepo(database)
			resultRepo := repository.NewSearchResultRepo(database)
			client := mcp.New(getMCPURL())
			service := search.NewService(queryRepo, resultRepo, client)

			queries, err := service.ListPendingQueries(cmd.Context(), limit)
			if err != nil {
				return err
			}
			out := make([]map[string]any, 0, len(queries))
			for _, q := range queries {
				count, err := service.SearchAndStore(cmd.Context(), q.ID, q.Query, pageSize)
				if err != nil {
					out = append(out, map[string]any{"query_id": q.ID, "query": q.Query, "error": err.Error()})
					continue
				}
				_ = queryRepo.UpdateStatus(cmd.Context(), q.ID, "done")
				out = append(out, map[string]any{"query_id": q.ID, "query": q.Query, "stored": count})
			}
			_ = json.NewEncoder(os.Stdout).Encode(out)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 5, "query limit")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "search page size")

	return cmd
}
