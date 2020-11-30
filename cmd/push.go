package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push --pr-node-id <node_id>",
	Short: "Publish metrics for this pull request.",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		verbose := viper.GetBool("verbose")
		dry := viper.GetBool("dry-run")

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

		prId := viper.GetString("pr-id")

		if prId == "" {
			return errors.New(`required flag "pr-id" not set`)
		}

		type Params struct {
			Id string
		}

		params := Params{prId}

		urlTemplate, err := template.New("apiUrl").Parse(apiUrl)

		if err != nil {
			printError(err)
		}

		buff := bytes.NewBufferString("")

		urlTemplate.Execute(buff, params)

		dataFile, err := os.Open(file)

		if err != nil {
			msg := fmt.Sprintf("Couldn't read report file \"%s\"", file)
			printError(msg, err)
		}

		defer dataFile.Close()

		url := buff.String()

		req, err := http.NewRequest("POST", url, dataFile)
		checkErr(err)

		token := fmt.Sprintf("Bearer %s", apiKey)

		req.Method = "PUT"
		req.Header.Set("Accepts", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		if verbose {
			fmt.Printf("%+v\n", req)
		}

		if dry {
			return nil
		}

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Fprintf(cmd.ErrOrStderr(), "Non 200 HTTP status: %+v\n", resp)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Fprintf(cmd.ErrOrStderr(), "Body: %+v\n", string(body))
			os.Exit(1)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(body))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolP("verbose", "v", false, "Print more information")
	pushCmd.Flags().BoolP("dry-run", "n", false, "Stop before publishing results")
	pushCmd.Flags().String("api-key", "", "Your SquidViz api-key (default is $SV_API_KEY)")
	pushCmd.Flags().String("api-url", "", "API endpoint to push to")

	pushCmd.Flags().StringP("pr-id", "p", "", `id of the pull request to publish metrics to.
If you're using GitHub Actions this should be
${{ github.event.pull_request.id }}. This the GitHub v3 API ID field, not the GraphQL
GID, and not SquidViz's internal ID.

See https://developer.github.com/v3/pulls/#list-pull-requests

If you have the GitHub CLI installed you can view this with:
	gh api repos/:user/:repo/pulls | jq 'values | .[] | {title:.title, id:.id}'
`)

	pushCmd.Flags().StringP("file", "f", "", "Input file containing all metrics to publish")
}
