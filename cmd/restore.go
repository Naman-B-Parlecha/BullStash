/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		input, _ := cmd.Flags().GetString("input")

		// iscompressed := strings.Contains(input, ".gz")

		if dbtype != "postgres" {
			util.CallWebHook("Unsupported database type: "+dbtype, true)
			fmt.Printf("Unsupported database type: %s\n", dbtype)
			return
		}

		fmt.Printf("Restoring database %s from %s\n", dbname, input)
		restoreCmd := exec.Command("psql",
			"-h", host,
			"-p", strconv.Itoa(port),
			"-U", user,
			"-d", dbname,
			"-f", input)

		restoreCmd.Env = append(restoreCmd.Env, "PGPASSWORD="+password)

		if err := restoreCmd.Run(); err != nil {
			util.CallWebHook("Failed to restore database: "+err.Error(), true)
			fmt.Printf("Failed to restore database: %v\n", err)
			return
		}

		util.CallWebHook("Database restored successfully", false)
		fmt.Println("Database restored successfully")
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().String("dbtype", "postgres", "Database type")
	restoreCmd.Flags().String("host", "localhost", "Database host")
	restoreCmd.Flags().Int("port", 5432, "Database port")
	restoreCmd.Flags().String("user", "postgres", "Database user")
	restoreCmd.Flags().String("password", "", "Database password")
	restoreCmd.Flags().String("dbname", "", "Database name")
	restoreCmd.Flags().String("input", "", "Input file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
