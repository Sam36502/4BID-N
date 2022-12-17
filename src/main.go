package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sqweek/dialog"
)

var g_currentScreen int = FBOD_SCRI_SPLASH
var g_insPointer byte = 0

func main() {
	err := LoadOptions()
	if err != nil {
		fmt.Printf("\n[ERROR]: Failed to load options file '%s'\n  Please make sure a valid options file is in the same directory as the 4BOD executable\n", FBOD_OPT_FILE)
		return
	}
	InitAudio()

	rand.Seed(time.Now().UnixMicro())
	rand.Int()

	ps := int32(g_options.PixelSize)
	rl.InitWindow(16*ps, 16*ps, "4BOD")
	if g_options.TargetFPS > 0 {
		rl.SetTargetFPS(int32(g_options.TargetFPS))
	}

	rl.SetExitKey(g_options.Controls.ExitKey)

	fvm := NewFBOD()

	lastKey := int32(0)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		switch g_currentScreen {

		case FBOD_SCRI_SPLASH:
			DrawBitmap(BMP_SPLASH)
			rl.EndDrawing()
			wait := g_options.SplashMillis * int(time.Millisecond)
			time.Sleep(time.Duration(wait))
			g_currentScreen = FBOD_SCRI_MENU

		case FBOD_SCRI_MENU:
			HandleMenu(fvm)
			DrawMenu()

		case FBOD_SCRI_EDITOR:
			HandleEditor(fvm)
			DrawEditor(fvm)

		case FBOD_SCRI_RUN:
			HandleRun()
			g_insPointer = fvm.PerformInstruction(g_insPointer)
			fvm.handleFPage()
			DrawBitmap(fvm.screen)
		}

		// Universal Keys
		if rl.IsKeyPressed(g_options.Controls.LoadKey) {
			filename, err := dialog.File().Filter("4BOD Binary File", "4bb").Title("Load 4BOD Program").Load()
			if err != nil {
				ErrorPopup("Failed to get filename")
			} else {
				err = fvm.LoadProgram(filename)
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
				err = fvm.SaveProgram(filename)
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
			rl.DrawText(fmt.Sprintf("Last Key Pressed: %d", lastKey), 5, 5, 20, g_colourFG)
		}

		rl.EndDrawing()
	}

	CloseAudio()
	rl.CloseWindow()
	SaveOptions()
}
