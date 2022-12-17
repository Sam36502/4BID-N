package main

import (
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	FBOD_SCRI_SPLASH = 0
	FBOD_SCRI_MENU   = 1
	FBOD_SCRI_EDITOR = 2
	FBOD_SCRI_RUN    = 3
)

var BMP_SPLASH = [16]uint16{
	0b1111111100000000,
	0b1001100101111100,
	0b1001100101100110,
	0b1001100101111100,
	0b1000000101100110,
	0b1111100101100110,
	0b1111100101111100,
	0b1111111100000000,
	0b0000000011111111,
	0b0011110010000011,
	0b0110011010011001,
	0b0110011010011001,
	0b0110011010011001,
	0b0110011010011001,
	0b0011110010000011,
	0b0000000011111111,
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

var BMP_OLD_MENU_EDIT = [16]uint16{
	0b0000000000000000,
	0b0011111111111000,
	0b0100010011000100,
	0b0100010101101100,
	0b0101110101101100,
	0b0100010011101100,
	0b0011111111111000,
	0b0000000000000000,
	0b0000000000000000,
	0b0000000000000000,
	0b0011101010110000,
	0b0010101010101000,
	0b0011001010101000,
	0b0010101110101000,
	0b0000000000000000,
	0b0000000000000000,
}

var BMP_OLD_MENU_RUN = [16]uint16{
	0b0000000000000000,
	0b0000000000000000,
	0b0011101100111000,
	0b0011101010010000,
	0b0010001010010000,
	0b0011101100010000,
	0b0000000000000000,
	0b0000000000000000,
	0b0000000000000000,
	0b0011111111111000,
	0b0100010101001100,
	0b0101010101010100,
	0b0100110101010100,
	0b0101010001010100,
	0b0011111111111000,
	0b0000000000000000,
}

var g_menuState = 0
var g_editorPage = 0

var g_colourFG = fourBitColour(0b0000) // Black
var g_colourBG = fourBitColour(0b0111) // Grey

func DrawBitmap(bmp [16]uint16) {
	rl.ClearBackground(g_colourBG)
	var x, y int32
	for y = 0; y < 16; y++ {
		for x = 0; x < 16; x++ {
			if (bmp[y]>>(15-x))%2 == 1 {
				ps := int32(g_options.PixelSize)
				rl.DrawRectangle(
					x*ps, y*ps,
					ps, ps,
					g_colourFG,
				)
			}
		}
	}
}

func DrawMenu() {
	if g_menuState == 0 {
		// 'EDT' Selected
		if g_options.OldMenu {
			DrawBitmap(BMP_OLD_MENU_EDIT)
		} else {
			DrawBitmap(BMP_MENU_EDIT)
		}
	} else {
		// 'RUN' Selected
		if g_options.OldMenu {
			DrawBitmap(BMP_OLD_MENU_RUN)
		} else {
			DrawBitmap(BMP_MENU_RUN)
		}
	}
}

func DrawEditor(f *FBOD) {
	rl.ClearBackground(g_colourBG)

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
	rl.DrawRectangle(13*ps, int32(g_editorPage)*ps, 2*ps, ps, g_colourFG)

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
					g_colourFG,
				)
			}
		}
	}
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

func fourBitColour(clr byte) rl.Color {
	clr %= 0xF
	factor := 128 * ((clr>>3)%2 + 1)
	return color.RGBA{
		R: factor * (clr >> 2 % 2),
		G: factor * (clr >> 1 % 2),
		B: factor * (clr % 2),
		A: 255,
	}
}
