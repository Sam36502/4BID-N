package main

import (
	"io"
	"io/ioutil"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sqweek/dialog"
)

type ControlConfig struct {
	Player1Input PlayerInputConfig `json:"player_1"`
	Player2Input PlayerInputConfig `json:"player_2"`
	ExitKey      int32             `json:"kc_exit"`
	BackKey      int32             `json:"kc_back"`
	LeftMouse    int32             `json:"mouse_left"`
}

type PlayerInputConfig struct {
	LeftKey  int32 `json:"kc_left"`
	RightKey int32 `json:"kc_right"`
	UpKey    int32 `json:"kc_up"`
	DownKey  int32 `json:"kc_down"`
	AKey     int32 `json:"kc_a"` // 'Select' in the menu
	BKey     int32 `json:"kc_b"` // 'Overlay' in the menu
	XKey     int32 `json:"kc_x"` // 'Load' in the menu
	YKey     int32 `json:"kc_y"` // 'Save' in the menu
}

func HandleMenu(f *CPU) {
	if rl.IsKeyPressed(g_options.Controls.Player1Input.UpKey) {
		g_menuState = 0
	}
	if rl.IsKeyPressed(g_options.Controls.Player1Input.DownKey) {
		g_menuState = 1
	}
	if rl.IsKeyPressed(g_options.Controls.Player1Input.AKey) {
		if g_menuState == 0 {
			g_currentScreen = SCRI_EDITOR
		} else {
			// Set up vm for running
			f.ClearMem()
			f.handleFPage()
			g_insPointer = 0

			g_currentScreen = SCRI_RUN
		}
	}
}

func HandleEditor(f *CPU) {
	if rl.IsKeyPressed(g_options.Controls.Player1Input.UpKey) {
		g_editorPage--
		if g_editorPage < 0 {
			g_editorPage = 0
		}
	}
	if rl.IsKeyPressed(g_options.Controls.Player1Input.DownKey) {
		g_editorPage++
		if g_editorPage > 15 {
			g_editorPage = 15
		}
	}
	if rl.IsKeyPressed(g_options.Controls.Player1Input.BKey) {
		g_options.EditorOverlay = !g_options.EditorOverlay
	}
	if rl.IsKeyPressed(g_options.Controls.BackKey) {
		g_currentScreen = SCRI_MENU
	}

	// Handle mouse clicks
	if rl.IsMouseButtonPressed(g_options.Controls.LeftMouse) {
		x := rl.GetMouseX() / int32(g_options.PixelSize)
		y := rl.GetMouseY() / int32(g_options.PixelSize)
		if x > 11 {
			return
		}

		ins := f.program[16*g_editorPage+int(y)]
		var insbits uint16 = (uint16(ins.ins) << 8) | (uint16(ins.arg1) << 4) | uint16(ins.arg2)
		insbits ^= 1 << (11 - x)
		f.program[16*g_editorPage+int(y)] = Instruction{
			ins:  (byte(insbits>>8) % 16),
			arg1: (byte(insbits>>4) % 16),
			arg2: byte(insbits % 16),
		}
	}
}

func HandleFileMenus(comp *CPU) {
	if rl.IsKeyPressed(g_options.Controls.Player1Input.XKey) {
		filename, err := dialog.File().Filter("4BOD Binary File", "4bb").Title("Load 4BOD Program").Load()
		if err != nil {
			ErrorPopup("Failed to get filename")
		} else {
			err = comp.LoadProgram(filename)
			if err != nil {
				ErrorPopup("Failed to load program")
			}
		}
	}

	if rl.IsKeyPressed(g_options.Controls.Player1Input.YKey) {
		filename, err := dialog.File().Filter("4BOD Binary File", "4bb").Title("Save 4BOD Program").Save()
		if err != nil {
			ErrorPopup("Failed to get filename")
		} else {
			err = comp.SaveProgram(filename)
			if err != nil {
				ErrorPopup("Failed to save program")
			}
		}
	}

}

func HandleRun(f *CPU) {
	if rl.IsKeyPressed(g_options.Controls.BackKey) {
		g_currentScreen = SCRI_MENU
		StopAllSounds()
	}
}

// Returns the state of a controller in two nybles
// D-Pad then Buttons
func GetControlNybles(in PlayerInputConfig) (byte, byte) {
	var dpad byte = 0

	if rl.IsKeyDown(in.LeftKey) {
		dpad |= 1
	}
	if rl.IsKeyDown(in.RightKey) {
		dpad |= 2
	}
	if rl.IsKeyDown(in.UpKey) {
		dpad |= 4
	}
	if rl.IsKeyDown(in.DownKey) {
		dpad |= 8
	}

	var btns byte = 0

	if rl.IsKeyDown(in.AKey) {
		btns |= 1
	}
	if rl.IsKeyDown(in.BKey) {
		btns |= 2
	}
	if rl.IsKeyDown(in.XKey) {
		btns |= 4
	}
	if rl.IsKeyDown(in.YKey) {
		btns |= 8
	}

	return dpad, btns
}

func ReadDisk(high, middle, low byte) (byte, error) {
	address := (uint16(high) << 8) | (uint16(middle) << 4) | uint16(low)
	hbit := address % 2
	offs := int64(address >> 1)
	file, err := os.Open(g_options.DiskFile)
	if err != nil {
		return 0, err
	}

	var b = []byte{0}
	_, err = file.ReadAt(b, offs)
	if err != nil && err != io.EOF {
		return 0, err
	}

	nyble := byte(0)
	if hbit == 0 {
		nyble = (b[0] >> 4) % 16
	} else {
		nyble = b[0] % 16
	}

	return nyble, file.Close()
}

func WriteDisk(high, middle, low byte, nyble byte) error {
	address := (uint16(high) << 8) | (uint16(middle) << 4) | uint16(low)
	hbit := address % 2
	offs := int64(address >> 1)
	nyble = nyble % 16

	data, err := ioutil.ReadFile(g_options.DiskFile)
	if err != nil {
		return err
	}

	for len(data) <= int(offs) {
		data = append(data, 0)
	}

	if hbit == 0 {
		data[offs] &= 0b00001111
		data[offs] |= nyble << 4
	} else {
		data[offs] &= 0b11110000
		data[offs] |= nyble
	}

	return ioutil.WriteFile(g_options.DiskFile, data, FILE_MODE)
}
