package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"xhs-go-cli/internal/logger"
)

var (
	cfgFile  string
	dbPath   string
	baseURL  string
	logLevel string
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "xhs-go-cli",
		Short: "xhs-go-cli: 小红书数据采集工具",
		Long:  `import-sources/query-gen/search/fetch-detail/qualify`,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "sqlite db path (default from config)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "mcp-url", "", "mcp base url (default from config)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level (default from config)")

	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("mcp.base-url", rootCmd.PersistentFlags().Lookup("mcp-url"))

	return rootCmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	viper.ReadInConfig()

	level := logLevel
	if level == "" {
		level = viper.GetString("log.level")
	}
	if level == "" {
		level = "info"
	}
	if err := logger.Init(level); err != nil {
		logger.Fatal("Failed to init logger", "error", err)
	}
}

func getDBPath() string {
	if dbPath != "" {
		return dbPath
	}
	return viper.GetString("db")
}

func getMCPURL() string {
	if baseURL != "" {
		return baseURL
	}
	return viper.GetString("mcp.base-url")
}
