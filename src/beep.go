package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Pitch uint32

const (
	PITCH_A3  = Pitch(0)
	PITCH_Bb3 = Pitch(1)
	PITCH_B3  = Pitch(2)
	PITCH_C4  = Pitch(3)
	PITCH_Cs4 = Pitch(4)
	PITCH_D4  = Pitch(5)
	PITCH_Ds4 = Pitch(6)
	PITCH_E4  = Pitch(7)
	PITCH_F4  = Pitch(8)
	PITCH_Fs4 = Pitch(9)
	PITCH_G4  = Pitch(10)
	PITCH_Gs4 = Pitch(11)
	PITCH_A4  = Pitch(12)
	PITCH_Bb4 = Pitch(13)
	PITCH_B4  = Pitch(14)
	PITCH_C5  = Pitch(15)
)

var g_beeperSample rl.Sound
var g_prevPitch Pitch
var g_prevVol byte

func InitAudio() {
	rl.InitAudioDevice()
	rl.SetMasterVolume(g_options.BeeperVol)
	g_beeperSample = rl.LoadSound(g_options.BeeperSample)
}

// Plays a beep with a given pitch and volume
func beep(p Pitch, v byte) {
	v %= 0xF
	if p != g_prevPitch || v != g_prevVol {
		if v > 0 {
			rl.SetSoundVolume(g_beeperSample, float32(0xF)/float32(v))
			rl.SetSoundPitch(g_beeperSample, 0.5*float32(math.Pow(1.0594, float64(p))))
			rl.PlaySound(g_beeperSample)
		} else {
			rl.StopSound(g_beeperSample)
		}
		g_prevPitch = p
		g_prevVol = v
	}
}

func CloseAudio() {
	rl.UnloadSound(g_beeperSample)
	rl.CloseAudioDevice()
}
