package cmd

import (
	"testing"

	"github.com/chdorner/keytographer/qmkapi"
	"github.com/stretchr/testify/require"
)

func TestFirstKeyboard(t *testing.T) {
	keyboard := qmkapi.Keyboard{KeyboardName: "First"}
	info := &qmkapi.KeyboardInfo{
		Keyboards: map[string]qmkapi.Keyboard{"first": keyboard},
	}

	name, kb, ok := firstKeyboard(info)
	require.True(t, ok)
	require.Equal(t, "first", name)
	require.Equal(t, &keyboard, kb)

	info = &qmkapi.KeyboardInfo{Keyboards: map[string]qmkapi.Keyboard{}}
	name, kb, ok = firstKeyboard(info)
	require.False(t, ok)
	require.Empty(t, name)
	require.Nil(t, kb)
}

func TestFindLayout(t *testing.T) {
	layout1 := qmkapi.Layout{Keys: []qmkapi.Key{{X: 1, Y: 1}}}
	layout2 := qmkapi.Layout{Keys: []qmkapi.Key{{X: 2, Y: 2}}}
	keyboard := &qmkapi.Keyboard{
		Layouts: map[string]qmkapi.Layout{"l_1": layout1, "l_2": layout2},
	}

	name, layout, err := findLayout("l_2", keyboard)
	require.NoError(t, err)
	require.Equal(t, "l_2", name)
	require.Equal(t, &layout2, layout)

	name, layout, err = findLayout("", keyboard)
	require.NoError(t, err)
	require.NotEmpty(t, name)
	require.NotNil(t, layout)

	name, layout, err = findLayout("l_3", keyboard)
	require.ErrorContains(t, err, "could not find layout with given name")
	require.Empty(t, name)
	require.Nil(t, layout)

	keyboard = &qmkapi.Keyboard{}
	name, layout, err = findLayout("", keyboard)
	require.ErrorContains(t, err, "could not find any layout")
	require.Empty(t, name)
	require.Nil(t, layout)
}

func TestInitConfig(t *testing.T) {
	layout := &qmkapi.Layout{Keys: []qmkapi.Key{
		{X: 1, Y: 1, W: 2, H: 2.5},
		{X: 3, Y: 1},
	}}
	config := initConfig("keyboard-name", "layout-name", layout)

	require.Equal(t, "keyboard-name", config.Keyboard)
	require.Equal(t, "layout-name", config.Layout.Macro)
	require.Len(t, config.Layout.Keys, 2)

	key := config.Layout.Keys[0]
	require.Equal(t, 1.0, key.X)
	require.Equal(t, 1.0, key.Y)
	require.Equal(t, 2.0, key.W)
	require.Equal(t, 2.5, key.H)

	key = config.Layout.Keys[1]
	require.Equal(t, 3.0, key.X)
	require.Equal(t, 1.0, key.Y)
	require.Equal(t, 1.0, key.W)
	require.Equal(t, 1.0, key.H)
}
