package cmd

import (
	"encoding/json"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/detail"
	"xhs-go-cli/internal/mcp"
	"xhs-go-cli/internal/repository"

	"github.com/spf13/cobra"
)

func NewFetchDetailCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "fetch-detail",
		Short: "获取内容详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			database, err := db.Open(getDBPath())
			if err != nil {
				return err
			}
			defer database.Close()

			searchResultRepo := repository.NewSearchResultRepo(database)
			detailRepo := repository.NewDetailRepo(database)
			client := mcp.New(getMCPURL())
			service := detail.NewService(searchResultRepo, detailRepo, client)

			rows, err := service.ListPending(cmd.Context(), limit)
			if err != nil {
				return err
			}
			result, err := service.FetchAndStore(cmd.Context(), rows)
			if err != nil {
				return err
			}
			_ = json.NewEncoder(os.Stdout).Encode(result)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 20, "detail row limit")

	return cmd
}
