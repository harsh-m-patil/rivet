/*
Copyright © 2026 Harshwardhan Patil
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/harsh-m-patil/rivet/internal/llm"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rivet",
	Short: "Minimal AI coding agent",
	Long:  `AI coding agent based on minimal philosophy`,
	Run: func(cmd *cobra.Command, args []string) {
		if prompt, _ := cmd.Flags().GetString("prompt"); prompt != "" {
			llm.GetAnswer(prompt)
			return
		}
		for {
			var prompt string
			print("> ")
			_, err := fmt.Scan(&prompt)
			if err != nil {
				slog.Error(err.Error())
			}
			if prompt != "" {
				llm.GetAnswer(prompt)
				println()
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rivet.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("prompt", "p", "", "Prompt to print")
}
