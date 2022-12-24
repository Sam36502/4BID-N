package main

import (
	"encoding/json"
	"io/ioutil"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	OPT_FILE   = "options.json"
	OPT_INDENT = "    "
)

var g_options Options = Options{}

type Options struct {
	DiskFile      string        `json:"disk_file"`
	SplashMillis  int           `json:"splash_millis"`
	PixelSize     int           `json:"pixel_size"`
	TargetFPS     int           `json:"target_fps"`
	BeeperVol     float32       `json:"beeper_vol"`
	BeeperSample  string        `json:"beeper_sample"`
	ColourFG      rl.Color      `json:"color_fg"`
	ColourBG      rl.Color      `json:"color_bg"`
	ColourOverlay rl.Color      `json:"color_overlay"`
	EditorOverlay bool          `json:"editor_overlay"`
	DebugKeycodes bool          `json:"debug_keycodes"`
	DebugMode     bool          `json:"debug_mode"`
	Controls      ControlConfig `json:"controls"`
}

func LoadOptions() error {
	data, err := ioutil.ReadFile(OPT_FILE)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &g_options)
}

func SaveOptions() error {
	data, err := json.MarshalIndent(g_options, "", OPT_INDENT)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(OPT_FILE, data, FILE_MODE)
}
