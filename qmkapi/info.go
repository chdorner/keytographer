package qmkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type KeyboardInfo struct {
	Keyboards map[string]Keyboard
}

type Keyboard struct {
	KeyboardName string `json:"keyboard_name"`
	Layouts      map[string]Layout
}

type Layout struct {
	Keys []Key `json:"layout"`
}

type Key struct {
	X float64
	Y float64
	W float64
	H float64
}

func Info(path string) (*KeyboardInfo, error) {
	url := fmt.Sprintf(`https://keyboards.qmk.fm/v1/%s`, path)
	logrus.WithField("url", url).Debug("Fetching QMK info.json")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected QMK API response with code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info KeyboardInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
