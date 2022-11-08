package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/chdorner/keytographer/internal/keytographer"
	"github.com/chdorner/keytographer/internal/qmkapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewGenerateCommand() *cobra.Command {
	var ctx context.Context
	var keyboardFlag string
	var pathFlag string
	var layoutFlag string
	var outFile string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate starting configuration.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx = createContext(cmd.Flags())
			configureLogging(ctx)

			keyboardFlag, _ = cmd.Flags().GetString("keyboard")
			if keyboardFlag == "" {
				return errors.New("missing keyboard name to fetch layou")
			}

			pathFlag, _ = cmd.Flags().GetString("path")
			layoutFlag, _ = cmd.Flags().GetString("layout")
			outFile, _ = cmd.Flags().GetString("out")

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			logrus.WithFields(logrus.Fields{
				"keyboard": keyboardFlag,
				"path":     pathFlag,
				"layout":   layoutFlag,
				"out":      outFile,
			})

			info, err := qmkapi.Info(keyboardFlag, pathFlag)
			if err != nil {
				logrus.Error(err)
				os.Exit(1)
			}

			keyboard, ok := info.Keyboards[fmt.Sprintf(`%s/%s`, keyboardFlag, pathFlag)]
			if !ok {
				logrus.Error("could not find keyboard with given name and path")
				os.Exit(1)
			}

			qmkLayout, ok := keyboard.Layouts[layoutFlag]
			if !ok {
				logrus.Error("could not find layout with given name")
				os.Exit(1)
			}

			layoutConfig := keytographer.LayoutConfig{}
			for _, qmkKey := range qmkLayout.Keys {
				layoutConfig.Keys = append(layoutConfig.Keys, keytographer.LayoutKeyConfig{
					X: qmkKey.X,
					Y: qmkKey.Y,
					W: qmkKey.W,
					H: qmkKey.H,
				})
			}
			config := keytographer.Config{
				Name:     "My awesome layout",
				Keyboard: keyboard.KeyboardName,
				Canvas: keytographer.CanvasConfig{
					Width:  800,
					Height: 600,
				},
				Layout: layoutConfig,
			}

			configYAML, err := yaml.Marshal(config)
			if err != nil {
				logrus.WithField("error", err).Error("failed to render YAML")
				os.Exit(1)
			}

			err = os.WriteFile(outFile, configYAML, 0644)
			if err != nil {
				logrus.WithField("error", err).Error("failed to write YAML to file")
				os.Exit(1)
			}
		},
	}

	fl := cmd.Flags()
	fl.StringP("keyboard", "k", "", "Name of the keyboard to fetch.")
	fl.StringP("path", "p", "", "Path to the directory containing info.json")
	fl.StringP("layout", "l", "", "Name of the layout to choose.")
	fl.StringP("out", "o", "", "Path to the output file.")

	return cmd
}
