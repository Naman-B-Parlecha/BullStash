package cmd

import (
	"fmt"
	"log"

	"github.com/Naman-B-Parlecha/BullStash/internal/backup"
	"github.com/Naman-B-Parlecha/BullStash/internal/config"
	"github.com/spf13/cobra"
)

// incrementalCmd represents the incremental command
var incrementalCmd = &cobra.Command{
	Use:   "incremental",
	Short: "Perform an incremental PostgreSQL backup",
	Long: `Perform an incremental backup of PostgreSQL database using WAL archiving.
	
This command requires:
1. PostgreSQL configured with wal_level = replica
2. archive_mode = on
3. Proper permissions to access the data directory

Example:
  bullstash incremental --config /path/to/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := config.LoadConfig(cfgPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// Verify PostgreSQL configuration first
		if err := backup.VerifyPGConfig(cfg); err != nil {
			log.Fatalf("PostgreSQL configuration error: %v", err)
		}

		// Take WAL backup
		if err := backup.PerformBackup(cfg); err != nil {
			log.Fatalf("Incremental backup failed: %v", err)
		}

		fmt.Println("✅ Incremental backup completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(incrementalCmd)
	incrementalCmd.Flags().StringP("config", "c", "", "Path to config file (required)")
	incrementalCmd.MarkFlagRequired("config")
}
