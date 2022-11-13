package cmd

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/chdorner/keytographer/config"
	"github.com/chdorner/keytographer/qmkapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewInitCommand() *cobra.Command {
	var ctx context.Context
	var infoPath string
	var outFile string

	var keyboardFlag string
	var layoutFlag string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a starting configuration",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx = createContext(cmd.Flags())
			configureLogging(ctx)

			infoPath, _ = cmd.Flags().GetString("info")
			if infoPath == "" {
				return errors.New("missing path to info.json to fetch layout")
			}
			if !strings.HasPrefix(infoPath, "keyboards/") {
				infoPath = "keyboards/" + infoPath
			}

			outFile, _ = cmd.Flags().GetString("out")
			if outFile == "" {
				return errors.New("missing path to the keytographer output file")
			}

			keyboardFlag, _ = cmd.Flags().GetString("keyboard")
			layoutFlag, _ = cmd.Flags().GetString("layout")

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			logrus.WithFields(logrus.Fields{
				"info":     infoPath,
				"out":      outFile,
				"keyboard": keyboardFlag,
				"layout":   layoutFlag,
			}).Debug("init")

			client := qmkapi.NewDefaultClient()
			info, err := client.Info(infoPath)
			if err != nil {
				logrus.WithField("error", err).Error("failed to fetch info.json from QMK API")
				os.Exit(1)
			}

			kbName, keyboard, ok := firstKeyboard(info)
			if !ok {
				logrus.Error("could not find keyboard with given name and path")
				os.Exit(1)
			}
			logrus.WithField("keyboard", kbName).Debug("found keyboard")

			layoutName, layout, err := findLayout(layoutFlag, keyboard)
			if err != nil {
				logrus.Error(err)
				os.Exit(1)
			}

			config := initConfig(kbName, layoutName, layout)
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
	fl.StringP("info", "i", "", "path to the info.json in QMK's repository")
	fl.StringP("out", "o", "", "path to the keytographer config output file")
	fl.StringP("keyboard", "k", "", "name of the keyboards")
	fl.StringP("layout", "l", "", "name of the layout macro function")

	return cmd
}

func firstKeyboard(info *qmkapi.KeyboardInfo) (string, *qmkapi.Keyboard, bool) {
	for key, kb := range info.Keyboards {
		return key, &kb, true
	}

	return "", nil, false
}

func findLayout(layoutFlag string, keyboard *qmkapi.Keyboard) (string, *qmkapi.Layout, error) {
	if layoutFlag != "" {
		l, ok := keyboard.Layouts[layoutFlag]
		if !ok {
			return "", nil, errors.New("could not find layout with given name")
		}
		return layoutFlag, &l, nil
	}

	var name string
	var layout *qmkapi.Layout

	for key, l := range keyboard.Layouts {
		layout = &l
		name = key
		break
	}
	if layout == nil {
		return "", nil, errors.New("could not find any layout")
	}
	logrus.WithField("layout", name).Debug("found first layout")
	return name, layout, nil
}

func initConfig(keyboardName, layoutName string, layout *qmkapi.Layout) *config.Config {
	layoutConfig := config.LayoutConfig{
		Macro: layoutName,
	}
	for _, qmkKey := range layout.Keys {
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
	return &config.Config{
		Name:     "My awesome layout",
		Keyboard: keyboardName,
		Canvas: config.CanvasConfig{
			Width:  800,
			Height: 600,
		},
		Layers: []config.Layer{
			{Name: "Base"},
		},
		Layout: layoutConfig,
	}
}
