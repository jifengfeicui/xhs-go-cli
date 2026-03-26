package cmd

import (
	"encoding/json"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/querygen"
	"xhs-go-cli/internal/repository"
	"xhs-go-cli/internal/search"

	"github.com/spf13/cobra"
)

func NewQueryGenCmd() *cobra.Command {
	var limit int
	var perSource int

	cmd := &cobra.Command{
		Use:   "query-gen",
		Short: "生成搜索关键词",
		RunE: func(cmd *cobra.Command, args []string) error {
			database, err := db.Open(getDBPath())
			if err != nil {
				return err
			}
			defer database.Close()

			sourceRepo := repository.NewSourceRepo(database)
			queryRepo := repository.NewQueryRepo(database)

			sources, err := sourceRepo.List(cmd.Context(), limit)
			if err != nil {
				return err
			}

			searchSvc := search.NewService(queryRepo, nil, nil)

			result := make([]map[string]any, 0, len(sources))
			for _, src := range sources {
				qsrc := querygen.Source{
					ID:         int64(src.ID),
					Name:       src.Name,
					Keywords:   src.Keywords,
					SourceType: querygen.ClassifySourceType(src.Name, src.Keywords),
				}
				queries := querygen.GenerateQueries(qsrc, perSource)
				newCount := 0
				for _, query := range queries {
					exists, _ := queryRepo.Exists(cmd.Context(), src.ID, query)
					if exists {
						continue
					}
					_ = searchSvc.SaveGeneratedQuery(cmd.Context(), src.ID, query, qsrc.SourceType)
					newCount++
				}
				if newCount > 0 {
					_ = sourceRepo.IncQueryCount(cmd.Context(), src.ID)
				}
				result = append(result, map[string]any{
					"source_id":   src.ID,
					"source_name": src.Name,
					"source_type": qsrc.SourceType,
					"queries":     queries,
				})
			}
			_ = json.NewEncoder(os.Stdout).Encode(result)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "source limit")
	cmd.Flags().IntVar(&perSource, "per-source", 3, "queries per source")

	return cmd
}
