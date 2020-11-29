package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		file, err := cmd.Flags().GetString("file")
		checkErr(err)

		apiKey, err := cmd.Flags().GetString("api-key")
		checkErr(err)

		apiUrl, err := cmd.Flags().GetString("api-url")
		checkErr(err)

		dataFile, err := os.Open(file)

		if err != nil {
			printError("couldn't read report file", err)
		}

		defer dataFile.Close()

		req, err := http.NewRequest("POST", apiUrl, dataFile)
		checkErr(err)

		token := fmt.Sprintf("Bearer %s", apiKey)

		req.Header.Set("Accepts", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			printError("couldn't publish results:", err)
		}

		defer resp.Body.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	// https://mnemosyne.dkendal.com/api/v1/pull_requests
	pushCmd.Flags().String("api-key", "", "Your SquidViz api-key (default is $SV_API_KEY)")
	pushCmd.Flags().String("api-url", "", "API endpoint to push to")
	pushCmd.Flags().StringP("file", "f", "", "Input file containing all metrics to publish")
}
