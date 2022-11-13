package qmkapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var moonlanderInfo = &KeyboardInfo{
	Keyboards: map[string]Keyboard{
		"crkbd": {
			KeyboardName: "Corne",
			Layouts: map[string]Layout{
				"LAYOUT_split_3x6_3": {Keys: []Key{{X: 3, Y: 6}}},
				"LAYOUT_split_3x5_3": {Keys: []Key{{X: 3, Y: 5}}},
			},
		},
	},
}

func TestInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/keyboards/crkbd/rev1/info.json" {
			resp, err := json.Marshal(moonlanderInfo)
			require.NoError(t, err)
			_, err = w.Write(resp)
			require.NoError(t, err)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	actual, err := client.Info("keyboards/crkbd/rev1/info.json")
	require.NoError(t, err)
	require.Equal(t, moonlanderInfo, actual)

	_, err = client.Info("keyboard/missing/info.json")
	require.ErrorContains(t, err, "unexpected QMK API response")
}
