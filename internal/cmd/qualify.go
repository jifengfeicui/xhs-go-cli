package cmd

import (
	"encoding/json"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/qualify"
	"xhs-go-cli/internal/repository"

	"github.com/spf13/cobra"
)

func NewQualifyCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "qualify",
		Short: "内容资质校验",
		RunE: func(cmd *cobra.Command, args []string) error {
			database, err := db.Open(getDBPath())
			if err != nil {
				return err
			}
			defer database.Close()

			detailRepo := repository.NewDetailRepo(database)
			qualRepo := repository.NewQualificationRepo(database)
			service := qualify.NewService(detailRepo, qualRepo)

			rows, err := service.ListDetails(cmd.Context(), limit)
			if err != nil {
				return err
			}
			result, err := service.QualifyAndStore(cmd.Context(), rows)
			if err != nil {
				return err
			}
			_ = json.NewEncoder(os.Stdout).Encode(result)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 20, "qualification row limit")

	return cmd
}
