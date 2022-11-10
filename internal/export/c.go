package export

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/chdorner/keytographer/internal/keytographer"
	"github.com/sirupsen/logrus"
)

var (
	//go:embed c.tpl
	cTpl      string
	markStart = "//keytographer:generated:start"
	markStop  = "//keytographer:generated:end"
)

type CExporter struct {
	ctx context.Context
}

func NewCExporter(ctx context.Context) *CExporter {
	return &CExporter{ctx}
}

func (e *CExporter) Export(config *keytographer.Config, outFile string) error {
	logrus.Debug("export config to c keymap")

	tpl, err := template.New("live").Parse(cTpl)
	if err != nil {
		return err
	}

	w := bytes.NewBufferString("")
	err = tpl.Execute(w, map[string]interface{}{
		"markStart": markStart,
		"markStop":  markStop,
		"config":    config,
	})
	if err != nil {
		return err
	}
	logrus.Debug("generated keymap")

	keymapb, err := os.ReadFile(outFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	keymap := string(keymapb)

	startIdx := strings.Index(keymap, markStart)
	stopIdx := strings.Index(keymap, markStop) + len(markStop) + 1

	var keymapOut string
	if startIdx == -1 || stopIdx == -1 {
		logrus.Debug("Could not find markers, appending keymap at the end")
		keymapOut = keymap + "\n" + w.String()
	} else {
		logrus.Debug("Found markers, replacing existing keymap")
		keymapOut = keymap[:startIdx] + w.String() + keymap[stopIdx:]
	}

	return os.WriteFile(outFile, []byte(keymapOut), 0644)
}
