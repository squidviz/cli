package cmd

import (
	"errors"
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

		file := viper.GetString("file")

		if file == "" {
			return errors.New(`required flag "file" not set`)
		}

		apiKey := viper.GetString("api-key")

		if apiKey == "" {
			return errors.New(`required flag "api-key" not set`)
		}

		apiUrl := viper.GetString("api-url")

		if apiUrl == "" {
			return errors.New(`required flag "api-url" not set`)
		}

		dataFile, err := os.Open(file)

		if err != nil {
			msg := fmt.Sprintf("Couldn't read report file \"%s\"", file)
			printError(msg, err)
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
