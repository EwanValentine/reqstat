package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reqstat",
	Short: "A beautiful HTTP request analyzer",
	Long: `reqstat is a CLI tool that makes HTTP requests and displays
detailed statistics including response time, size, headers,
and JSON structure analysis.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
