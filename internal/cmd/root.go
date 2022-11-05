package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "keymap-render",
	Short: "Render keymaps as SVG.",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode.")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to the keymap configuration file to watch for changes.")

	rootCmd.AddCommand(NewRenderCommand())
	rootCmd.AddCommand(NewLiveCommand())
}

func createContext(flags *pflag.FlagSet) context.Context {
	debug, _ := flags.GetBool("debug")
	return context.WithValue(context.Background(), "debug", debug)
}

func configureLogging(ctx context.Context) {
	if ctx.Value("debug").(bool) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode turned on")
	}
}
