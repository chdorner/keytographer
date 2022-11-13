package cmd

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "keytographer",
	Short: "Render beautiful keymap visualizations",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func RootCommand() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")

	rootCmd.AddCommand(NewRenderCommand())
	rootCmd.AddCommand(NewValidateCommand())
	rootCmd.AddCommand(NewExportCommand())

	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewLiveCommand())

	rootCmd.AddCommand(NewVersionCommand())
}

func createContext(flags *pflag.FlagSet) context.Context {
	debug, _ := flags.GetBool("debug")
	return context.WithValue(context.Background(), "debug", debug) //nolint:staticcheck
}

func configureLogging(ctx context.Context) {
	if ctx.Value("debug").(bool) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode turned on")
	}
}
