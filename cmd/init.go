package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/chdorner/keytographer/config"
	"github.com/chdorner/keytographer/qmkapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewInitCommand() *cobra.Command {
	var ctx context.Context
	var keyboardFlag string
	var pathFlag string
	var layoutFlag string
	var outFile string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a starting configuration",

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

			layoutConfig := config.LayoutConfig{}
			for _, qmkKey := range qmkLayout.Keys {
				w, h := 1.0, 1.0
				if qmkKey.W > 0 {
					w = qmkKey.W
				}
				if qmkKey.H > 0 {
					h = qmkKey.H
				}
				layoutConfig.Keys = append(layoutConfig.Keys, config.LayoutKeyConfig{
					X: qmkKey.X,
					Y: qmkKey.Y,
					W: w,
					H: h,
				})
			}
			config := config.Config{
				Name:     "My awesome layout",
				Keyboard: keyboard.KeyboardName,
				Canvas: config.CanvasConfig{
					Width:  800,
					Height: 600,
				},
				Layers: []config.Layer{
					{Name: "Base"},
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
	fl.StringP("keyboard", "k", "", "name of the keyboard to fetch")
	fl.StringP("path", "p", "", "path to the remote directory containing info.json")
	fl.StringP("layout", "l", "", "name of the layout macro function")
	fl.StringP("out", "o", "", "path to the keytographer config output file")

	return cmd
}
