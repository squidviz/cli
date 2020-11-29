package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add --label <label>",
	Short: "Add a new metric to the output file from stdin",
	Long: `"add" reads a decimal value from stdin and adds it as a new entry to
the output file specified by --file. The command prints the value of --label,
along with the input to stdout.
  `,
	Example: `
Supplying the metric value through a pipe.

  sh$ echo "1.3" | sv add --label "some label"
  some label 1.3

You can also capture more detailed timing values with GNU time(1).

  bash$ command time -f "%U + %S" go test 2> >(bc | sv add --label "go test")
  PASS
  ok  	sv	0.005s
  go test 1.64

The format string "%U + %S" reports the user time and system time. Because
time(1) outputs the timing output over stderr, and cannot eval math
expressions, stderr is fed through process substitution into bc(1) and finally
piped to "add".

If you can't use bash you can swap stderr and stdout like below.

  sh$ command time -f "%U + %S" 3>&2 2>&1 1>&3 go test |
  ..$   bc | sv add --label "go test"
  ok  	sv	0.005s
  go test 1.94
  `,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		var value float64

		typevalue := viper.GetString("type")

		label := viper.GetString("label")

		if label == "" {
			return errors.New(`required flag "label" not set`)
		}

		file := viper.GetString("file")

		if file == "" {
			return errors.New(`required flag "file" not set`)
		}

		_, err := fmt.Fscanln(cmd.InOrStdin(), &value)
		checkErr(err)

		metric := Metric{
			Label: label,
			Value: value,
			Type:  typevalue,
		}

		appendMetric(file, metric)

		fmt.Println(label, value)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("file", "f", "", `Output report file (default $SV_FILE, ".svfile")`)
	addCmd.Flags().StringP("label", "l", "", "Label for this data point (default $SV_LABEL)")

	addCmd.Flags().StringP("type", "t", "", "Set the datatype (default $SV_TYPE)")
	addCmd.Flags().MarkHidden("type")
}
