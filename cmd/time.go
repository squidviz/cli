package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time --label <label>",
	Short: "Measure how long a command takes",
	Example: `
  $ ls | time --label "ls"
  file.go
  ls : 59.743µs
  `,
	Long: `"time" is a convience function to record how long the command piped into it
from stdin takes to complete. This should be treated as rough estimate. For
more advanced usages use the "add" command with GNU time(1).

  command time ls 2> >(sv add --label ls)
`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		viper.BindPFlags(cmd.Flags())

		file := viper.GetString("file")

		if file == "" {
			return errors.New(`required flag "file" not set`)
		}

		label := viper.GetString("label")

		if label == "" {
			return errors.New(`required flag "label" not set`)
		}

		// Buffer all output from stdin and print to output.
		// Just acts as a pass through so we can capture the "real time" of
		// the command execution
		scanner := bufio.NewScanner(cmd.InOrStdin())

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				printError("failed to read input,", err)
			}
			fmt.Println(scanner.Text())
		}

		end := time.Now()
		duration := end.Sub(start)

		metric := Metric{
			Label: label,
			Value: duration.Seconds(),
			Type:  "seconds",
		}

		appendMetric(file, metric)

		fmt.Fprintln(os.Stderr, label, duration)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().StringP("file", "f", "", `Report output file (default $SV_FILE, ".svfile")`)
	timeCmd.Flags().StringP("label", "l", "", `Set the label for this data point (default $SV_LABEL)`)
}
