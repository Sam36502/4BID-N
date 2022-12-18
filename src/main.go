package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sqweek/dialog"
)

const (
	WINDOW_TITLE = "FBID-N"
)

var g_currentScreen int
var g_insPointer byte

func main() {
	err := LoadOptions()
	if err != nil {
		fmt.Printf(
			"\n[ERROR]: Failed to load options file '%s'\n  Please make sure a valid options file is in the same directory as the executable\n",
			OPT_FILE,
		)
		return
	}
	InitAudio()

	rand.Seed(time.Now().UnixMicro())
	rand.Int()

	ps := int32(g_options.PixelSize)
	rl.InitWindow(16*ps, 16*ps, WINDOW_TITLE)
	if g_options.TargetFPS > 0 {
		rl.SetTargetFPS(int32(g_options.TargetFPS))
	}

	rl.SetExitKey(g_options.Controls.ExitKey)

	comp := NewCPU()

	lastKey := int32(0)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		switch g_currentScreen {

		case SCRI_SPLASH:
			DrawBitmap(BMP_SPLASH)
			rl.EndDrawing()
			wait := g_options.SplashMillis * int(time.Millisecond)
			time.Sleep(time.Duration(wait))
			g_currentScreen = SCRI_MENU

		case SCRI_MENU:
			HandleMenu(comp)
			DrawMenu()

		case SCRI_EDITOR:
			HandleEditor(comp)
			DrawEditor(comp)

		case SCRI_RUN:
			HandleRun()
			comp.HandleScreen() // Updates F-Page and screen
			g_insPointer = comp.PerformInstruction(g_insPointer)
			comp.handleFPage()
			DrawBitmap(comp.screen)
		}

		// Universal Keys
		if rl.IsKeyPressed(g_options.Controls.LoadKey) {
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

		if rl.IsKeyPressed(g_options.Controls.SaveKey) {
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

		if g_options.DebugKeycodes {
			key := rl.GetKeyPressed()
			if key != int32(lastKey) && key != 0 {
				lastKey = key
			}
			rl.DrawText(fmt.Sprintf("Last Key Pressed: %d", lastKey), 5, 5, 20, g_options.ColourFG)
		}

		rl.EndDrawing()
	}

	CloseAudio()
	rl.CloseWindow()
	SaveOptions()
}
