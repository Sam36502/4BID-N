package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
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
			HandleFileMenus(comp)
			HandleMenu(comp)
			DrawMenu()

		case SCRI_EDITOR:
			HandleFileMenus(comp)
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
