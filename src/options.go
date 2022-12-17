package main

import (
	"encoding/json"
	"io/ioutil"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	FBOD_OPT_FILE   = "options.json"
	FBOD_OPT_INDENT = "    "
)

var g_options Options = Options{}

type Options struct {
	SplashMillis  int           `json:"splash_millis"`
	PixelSize     int           `json:"pixel_size"`
	TargetFPS     int           `json:"target_fps"`
	BeeperVol     float32       `json:"beeper_vol"`
	BeeperSample  string        `json:"beeper_sample"`
	ColourOverlay rl.Color      `json:"color_overlay"`
	OldMenu       bool          `json:"old_menu"`
	EditorOverlay bool          `json:"editor_overlay"`
	DebugKeycodes bool          `json:"debug_keycodes"`
	Controls      ControlConfig `json:"controls"`
}

func LoadOptions() error {
	data, err := ioutil.ReadFile(FBOD_OPT_FILE)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &g_options)
}

func SaveOptions() error {
	data, err := json.MarshalIndent(g_options, "", FBOD_OPT_INDENT)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(FBOD_OPT_FILE, data, FBOD_FILE_MODE)
}
