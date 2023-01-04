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

const (
	OCTAVE_2 = 0
	OCTAVE_3 = 1
	OCTAVE_4 = 2
	OCTAVE_5 = 3
)

const (
	WAVE_SQUARE   = 0
	WAVE_TRIANGLE = 1
	WAVE_SAWTOOTH = 2
	WAVE_NOISE    = 3
)

const SEMITONE_INTERVAL = 1.059463

var g_samples [4]rl.Sound
var g_prevWave byte
var g_prevOct byte
var g_prevPitch Pitch
var g_prevVol byte

func InitAudio() {
	rl.InitAudioDevice()
	rl.SetMasterVolume(g_options.MasterVol)
	g_samples = [4]rl.Sound{
		WAVE_SQUARE:   rl.LoadSound(g_options.SquareSample),
		WAVE_TRIANGLE: rl.LoadSound(g_options.TriangleSample),
		WAVE_SAWTOOTH: rl.LoadSound(g_options.SawtoothSample),
		WAVE_NOISE:    rl.LoadSound(g_options.NoiseSample),
	}
}

// Plays a beep with a given pitch, octave, waveform and volume
func beep(ptc Pitch, oct byte, wave byte, vol byte) {
	vol %= 0xF
	if wave != g_prevWave || oct != g_prevOct || ptc != g_prevPitch || vol != g_prevVol {
		if vol > 0 {
			rl.SetSoundVolume(g_samples[wave], (1/float32(0xF))*float32(vol))
			pitchMul := float32(math.Pow(SEMITONE_INTERVAL, float64(ptc))) * float32(oct+1)
			rl.SetSoundPitch(g_samples[wave], pitchMul)
			rl.PlaySound(g_samples[wave])
		} else {
			rl.StopSound(g_samples[wave])
		}
		g_prevWave = wave
		g_prevOct = oct
		g_prevPitch = ptc
		g_prevVol = vol
	}
}

func StopAllSounds() {
	for _, s := range g_samples {
		rl.StopSound(s)
	}
}

func CloseAudio() {
	for _, s := range g_samples {
		rl.UnloadSound(s)
	}
	rl.CloseAudioDevice()
}
