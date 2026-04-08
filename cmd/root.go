/*
Copyright © 2026 Harshwardhan Patil
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/harsh-m-patil/rivet/internal/llm"
	"github.com/openai/openai-go/v3"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rivet",
	Short: "Minimal AI coding agent",
	Long:  `AI coding agent based on minimal philosophy`,
	Run: func(cmd *cobra.Command, args []string) {
		client := &llm.Client{}
		client.Get()
		defer client.Close()

		if prompt, _ := cmd.Flags().GetString("prompt"); prompt != "" {
			processPrompt(client, prompt)
			return
		}
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					slog.Error("input error", "err", err)
				}
				break
			}
			prompt := scanner.Text()
			if prompt != "" {
				fmt.Printf("< %s\n", prompt)
				processPrompt(client, prompt)
				fmt.Println()
			}
		}
	},
}

func processPrompt(client *llm.Client, prompt string) {
	ctx := context.Background()
	messages := llm.Messages{
		openai.UserMessage(prompt),
	}

	resultChan := client.ChatCompletion(ctx, messages, true)
	for result := range resultChan {
		if result.Err != nil {
			slog.Error("error", "err", result.Err)
			continue
		}

		switch result.Event.Type {
		case llm.TEXT_DELTA:
			if result.Event.TextDelta != nil {
				fmt.Print(result.Event.TextDelta.Content)
			}
		case llm.MESSAGE_COMPLETE:
			if result.Event.FinishReason != nil {
				slog.Debug("message complete", "reason", *result.Event.FinishReason)
			}
		case llm.ERROR:
			if result.Event.Error != nil {
				slog.Error("api error", "error", *result.Event.Error)
			}
		}
	}
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
