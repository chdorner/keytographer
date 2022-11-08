package qmkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func Info(keyboard, path string) (*KeyboardInfo, error) {
	resp, err := http.Get(fmt.Sprintf(`https://keyboards.qmk.fm/v1/keyboards/%s/%s/info.json`, keyboard, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
