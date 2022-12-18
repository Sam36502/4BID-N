package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ControlConfig struct {
	LeftKey    int32 `json:"kc_left"`
	RightKey   int32 `json:"kc_right"`
	UpKey      int32 `json:"kc_up"`
	DownKey    int32 `json:"kc_down"`
	ExitKey    int32 `json:"kc_exit"`
	SelectKey  int32 `json:"kc_enter"`
	BackKey    int32 `json:"kc_back"`
	SaveKey    int32 `json:"kc_save"`
	LoadKey    int32 `json:"kc_load"`
	OptionsKey int32 `json:"kc_opts"`
	OverlayKey int32 `json:"kc_overlay"`
	LeftMouse  int32 `json:"mouse_left"`
}

func HandleMenu(f *CPU) {
	if rl.IsKeyPressed(g_options.Controls.UpKey) {
		g_menuState = 0
	}
	if rl.IsKeyPressed(g_options.Controls.DownKey) {
		g_menuState = 1

		// Set up vm for running
		f.ClearMem()
		g_insPointer = 0
	}
	if rl.IsKeyPressed(g_options.Controls.SelectKey) {
		if g_menuState == 0 {
			g_currentScreen = SCRI_EDITOR
		} else {
			g_currentScreen = SCRI_RUN
		}
	}
}

func HandleEditor(f *CPU) {
	if rl.IsKeyPressed(g_options.Controls.UpKey) {
		g_editorPage--
		if g_editorPage < 0 {
			g_editorPage = 0
		}
	}
	if rl.IsKeyPressed(g_options.Controls.DownKey) {
		g_editorPage++
		if g_editorPage > 15 {
			g_editorPage = 15
		}
	}
	if rl.IsKeyPressed(g_options.Controls.OverlayKey) {
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

func HandleRun() {
	if rl.IsKeyPressed(g_options.Controls.BackKey) {
		g_currentScreen = SCRI_MENU
	}
}

func (c *CPU) HandleInputs() {
	// Player 1
	c.mem[0xF][FPG_P1_DPAD] = GetArrowsNybl()

	// TODO: Player 2
}

// Returns the state of the arrow keys as a nybl:
// 8s: Down
// 4s: Up
// 2s: Right
// 1s: Left
func GetArrowsNybl() byte {
	var nybl byte = 0

	if rl.IsKeyDown(g_options.Controls.LeftKey) {
		nybl |= 1
	}
	if rl.IsKeyDown(g_options.Controls.RightKey) {
		nybl |= 2
	}
	if rl.IsKeyDown(g_options.Controls.DownKey) {
		nybl |= 4
	}
	if rl.IsKeyDown(g_options.Controls.UpKey) {
		nybl |= 8
	}

	return nybl
}
