package cmd

import (
	"encoding/json"
	"os"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/repository"
	"xhs-go-cli/internal/source"

	"github.com/spf13/cobra"
)

func NewImportSourcesCmd() *cobra.Command {
	var input string

	cmd := &cobra.Command{
		Use:   "import-sources",
		Short: "导入数据源",
		RunE: func(cmd *cobra.Command, args []string) error {
			database, err := db.Open(getDBPath())
			if err != nil {
				return err
			}
			defer database.Close()

			repo := repository.NewSourceRepo(database)
			svc := source.NewRepo(repo)
			count, err := source.ImportFromJSON(cmd.Context(), svc, input)
			if err != nil {
				return err
			}
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"imported": count})
			return nil
		},
	}

	cmd.Flags().StringVar(&input, "input", "", "source json path")
	cmd.MarkFlagRequired("input")

	return cmd
}
