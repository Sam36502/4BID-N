package main

import (
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCRI_SPLASH = 0
	SCRI_MENU   = 1
	SCRI_EDITOR = 2
	SCRI_RUN    = 3
)

var BMP_SPLASH = [16]uint16{
	0b0000000000000000,
	0b0101011001011100,
	0b0101010100010010,
	0b0111011001010010,
	0b0001010101010010,
	0b0001010101010010,
	0b0001011101011110,
	0b0000000000000000,
	0b0111111111111110,
	0b0000000000000000,
	0b0000001000001100,
	0b0000010100010000,
	0b0000100010100000,
	0b0011000001000000,
	0b0000000000000000,
	0b0000000000000000,
}

var BMP_MENU_EDIT = [16]uint16{
	0b0000000000000000,
	0b0000000000000000,
	0b0111011001110010,
	0b0111010100100110,
	0b0100010100100110,
	0b0111011000100010,
	0b0000000000000000,
	0b0000000000000000,
	0b0000000000000000,
	0b0000000000000000,
	0b0111010101100000,
	0b0101010101010000,
	0b0110010101010000,
	0b0101011101010000,
	0b0000000000000000,
	0b0000000000000000,
}

var BMP_MENU_RUN = [16]uint16{
	0b0000000000000000,
	0b0000000000000000,
	0b0111011001110000,
	0b0111010100100000,
	0b0100010100100000,
	0b0111011000100000,
	0b0000000000000000,
	0b0000000000000000,
	0b0000000000000000,
	0b0000000000000000,
	0b0111010101100010,
	0b0101010101010110,
	0b0110010101010110,
	0b0101011101010010,
	0b0000000000000000,
	0b0000000000000000,
}

var g_menuState = 0
var g_editorPage = 0

func DrawBitmap(bmp [16]uint16) {
	rl.ClearBackground(g_options.ColourBG)
	var x, y int32
	for y = 0; y < 16; y++ {
		for x = 0; x < 16; x++ {
			if (bmp[y]>>(15-x))%2 == 1 {
				ps := int32(g_options.PixelSize)
				rl.DrawRectangle(
					x*ps, y*ps,
					ps, ps,
					g_options.ColourFG,
				)
			}
		}
	}
}

func DrawMenu() {
	if g_menuState == 0 {
		// 'EDT' Selected
		DrawBitmap(BMP_MENU_EDIT)
	} else {
		// 'RUN' Selected
		DrawBitmap(BMP_MENU_RUN)
	}
}

func DrawEditor(f *CPU) {
	rl.ClearBackground(g_options.ColourBG)

	// Draw Overlay
	ps := int32(g_options.PixelSize)
	if g_options.EditorOverlay {
		for i := int32(0); i < 16; i += 2 {
			rl.DrawRectangle(0, i*ps, 4*ps, ps, g_options.ColourOverlay)
			rl.DrawRectangle(4*ps, i*ps+ps, 4*ps, ps, g_options.ColourOverlay)
			rl.DrawRectangle(8*ps, i*ps, 4*ps, ps, g_options.ColourOverlay)
		}
	}

	// Draw 'Scrollbar'
	rl.DrawRectangle(13*ps, int32(g_editorPage)*ps, 2*ps, ps, g_options.ColourFG)

	// Draw Program pixels
	var x, y int32
	for y = 0; y < 16; y++ {
		ins := f.program[16*g_editorPage+int(y)]
		var insbits uint16 = (uint16(ins.ins) << 8) | (uint16(ins.arg1) << 4) | uint16(ins.arg2)
		for x = 0; x < 12; x++ {
			if (insbits>>(11-x))%2 == 1 {
				ps := int32(g_options.PixelSize)
				rl.DrawRectangle(
					x*ps, y*ps,
					ps, ps,
					g_options.ColourFG,
				)
			}
		}
	}
}

// Handles any changes to the screen based on the F-Page addresses
func (f *CPU) HandleScreen() {

	// Handle special Options
	opt := f.mem[0xF][FPG_SCR_OPT]

	// Clear Screen
	if (opt>>3)%2 == 1 {
		f.ClearScreen()
		f.mem[0xF][FPG_SCR_OPT] ^= 1 << 3
	}

	// Invert Screen
	if (opt>>2)%2 == 1 {
		for i, row := range f.screen {
			f.screen[i] = ^row
		}
		f.mem[0xF][FPG_SCR_OPT] ^= 1 << 2
	}

	x := f.mem[0xF][FPG_SCR_X]
	y := f.mem[0xF][FPG_SCR_Y]

	// Read pixel value
	var val byte
	switch opt % 4 {

	case 0b00:
		val = byte((f.screen[y] >> (15 - x)) % 2)

	case 0b01:
		val = byte((f.screen[y] >> (15 - x)) % 0xF)

	case 0b10:
		for i := byte(0); i < 4; i++ {
			val |= byte((f.screen[y+i]>>(15-x))%2) << (3 - i)
		}

	case 0b11:
		val |= byte((f.screen[y]>>(15-x))%4) << 2
		if y < 15 {
			val |= byte((f.screen[y+1] >> (15 - x)) % 4)
		}

	}

	f.mem[0xF][FPG_SCR_VAL] = val
}

func ErrorPopup(msg string) {
	// Draw Box
	red := color.RGBA{255, 64, 64, 255}
	darkRed := color.RGBA{200, 32, 32, 255}
	width := 300
	height := 150
	x := (g_options.PixelSize * 16 / 2) - width/2
	y := (g_options.PixelSize * 16 / 2) - height/2
	rec := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: float32(height)}
	rl.DrawRectangleRec(rec, red)
	rl.DrawRectangleLinesEx(rec, 5, darkRed)

	// Draw Text
	rl.DrawText(msg, int32(x+25), int32(y+25), 20, darkRed)
	rl.EndDrawing()
	time.Sleep(3 * time.Second)
}
