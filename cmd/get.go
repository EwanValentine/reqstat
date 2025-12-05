package cmd

import (
	"github.com/ewan-valentine/reqstat/internal/client"
	"github.com/ewan-valentine/reqstat/internal/display"
	"github.com/spf13/cobra"
)

var (
	headers     []string
	showBody    bool
	prettyJSON  bool
	maxBodySize int
)

var getCmd = &cobra.Command{
	Use:   "get <url>",
	Short: "Make a GET request and analyze the response",
	Long: `Make an HTTP GET request to the specified URL and display
detailed statistics about the response including timing,
size, headers, and JSON structure.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		c := client.New()
		for _, h := range headers {
			c.AddHeader(h)
		}

		result, err := c.Get(url)
		if err != nil {
			return err
		}

		opts := display.Options{
			ShowBody:    showBody,
			PrettyJSON:  prettyJSON,
			MaxBodySize: maxBodySize,
		}

		display.Render(result, opts)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "Add custom header (can be repeated)")
	getCmd.Flags().BoolVarP(&showBody, "body", "b", false, "Show response body")
	getCmd.Flags().BoolVarP(&prettyJSON, "pretty", "p", true, "Pretty print JSON body")
	getCmd.Flags().IntVarP(&maxBodySize, "max-body", "m", 1000, "Max body characters to display")
}
