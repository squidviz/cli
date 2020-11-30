package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
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

		type Flags struct {
			File  string `validate:"required"`
			Label string `validate:"required"`
			Magnitude  string `validate:"required,oneof=s seconds ms milliseconds ns nanosecond"`
			Help  string ``
		}

		var flags Flags

		err := viper.Unmarshal(&flags)

		if err != nil {
			panic(err)
		}

		validate := validator.New()
		err = validate.Struct(flags)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Fprintf(cmd.ErrOrStderr(), "%+v\nOptions: %+v\nValue: %+v\n", err, err.Value(), err.Param())
			}

			return errors.New("Argument error")
		}

		// Buffer all output from STDIN and print to output.
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

		var value float64
		var mag string

		switch flags.Magnitude {
		case "s", "seconds":
			value = duration.Seconds()
			mag = "seconds"
		case "ms", "milliseconds":
			value = float64(duration.Milliseconds())
			mag = "milliseconds"
		case "ns", "nanoseconds":
			value = float64(duration.Nanoseconds())
			mag = "nanoseconds"
		}

		metric := Metric{
			Label: flags.Label,
			Value: value,
			Type:  mag,
		}

		appendMetric(flags.File, metric)

		fmt.Fprintln(os.Stderr, flags.Label, value, mag)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().StringP("file", "f", "", `Report output file (default $SV_FILE, ".svfile")`)
	timeCmd.Flags().StringP("label", "l", "", `Set the label for this data point (default $SV_LABEL)`)
	timeCmd.Flags().StringP("magnitude", "m", "", `Unit of time to record with. (default $SV_MAGNITUDE, "ms")
	One of "seconds", "milliseconds", or "nanoseconds", or their SI symbol
	equivalents ("s", "ms", "ns"). Note that any unit other than "seconds" is
	reported as an integer and will have a loss in precision!`)
}
