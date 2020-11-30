package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	// Default arguments for viper
	defaults = map[string]interface{}{
		"file":      ".svfile",
		"magnitude": "ms",
		"api-url":   "https://mnemosyne.dkendal.com/api/v1/pull_requests/{{.Id}}",
	}
)

type Root struct {
	PullRequest PullRequest `json:"pull_request"`
}

type PullRequest struct {
	Metrics []Metric `json:"data"`
}

type Metric struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

func checkErr(err error) {
	if err != nil {
		printError(err)
	}
}

func printError(msg ...interface{}) {
	msg = append([]interface{}{"Error:"}, msg...)
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

func appendMetric(path string, metric Metric) {
	var root Root

	file, err := os.OpenFile(path, os.O_RDWR, 0600)

	if err == nil {
		// Decode JSON if the file already exists
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&root)

		if err != nil {
			log.Fatalln("Failed to parse JSON:", err)
		}
	}

	if os.IsNotExist(err) {
		// If the file doesn't exist create and open it
		file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
		defer file.Close()

		if err != nil {
			log.Fatalln("Failed to create file:", err)
		}
	}

	// Truncate the file and rewind to the beginning of the file before writing
	file.Truncate(0)
	file.Seek(0, 0)

	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatalln(err)
	}

	root.PullRequest.Metrics = append(root.PullRequest.Metrics, metric)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(root)

	if err != nil {
		log.Fatalln("Failed to write results:", err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sv",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.squidviz.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".squidviz" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".squidviz")
		viper.SetEnvPrefix("sv")
		viper.AutomaticEnv() // read in environment variables that match
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
