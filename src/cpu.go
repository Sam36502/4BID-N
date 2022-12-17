package main

import (
	"io/ioutil"
	"math/rand"
)

// 4BID-N Instructions
const (
	ASM_NOP  = 0x0 // Do nothing
	ASM_LDAI = 0x1 // Load immediate value to acc
	ASM_LDAM = 0x2 // Load memory value to acc
	ASM_STA  = 0x3 // Store acc to memory

	ASM_INC = 0x4 // Increment/Decrement acc
	ASM_ADD = 0x5 // Add memory value to acc
	//ASM_001 = 0x6 //
	//ASM_SHL = 0x7 //

	ASM_NOT = 0x8 // Bitwise NOT
	ASM_ORA = 0x9 // Bitwise OR memory value and acc
	ASM_AND = 0xA // Bitwise AND memory value and acc
	ASM_SHF = 0xB // Bitwise shift (l/r & rot based on high bits)

	//ASM_JMP = 0xC //
	//ASM_CEQ = 0xD //
	ASM_BNE = 0xE // Skips B many instructions if acc does not equal A
	ASM_JMP = 0xF // Jump to specific point in program
)

// 4BID-N F-Page Addresses
const (
	FPG_INPUT = 0x0 // Arrow keys state (D, U, R, L)

	FPG_RAND = 0x1 // Random Number
	FPG_DECF = 0x2 // Decremented every frame
	FPG_SAME = 0x3 // Same as 0011 (?)

	FPG_REDIR    = 0x4 // Memory Redirect
	FPG_REDIR_AD = 0x5 // Memory Redirect Address
	FPG_REDIR_PG = 0x6 // Memory Redirect Page

	FPG_BEEP_P  = 0x7 // Beeper Pitch
	FPG_BEEP_V  = 0x8 // Beeper Volume
	FPG_BEEP_DP = 0x9 // Beeper Delta Pitch
	FPG_BEEP_DV = 0xA // Beeper Delta Volume

	FPG_CLR_FG = 0xB // Screen Colour Foreground
	FPG_CLR_BG = 0xC // Screen Colour Background

	FPG_SCORE_1     = 0xD // Score 1
	FPG_SCORE_2     = 0xE // Score 2
	FPG_SCORE_STATE = 0xF // See documentation
)

type Instruction struct {
	ins  byte
	arg1 byte
	arg2 byte
}

type FBOD struct {
	acc     byte
	mem     [16][16]byte
	flags   [16]byte   // List of program addresses
	screen  [16]uint16 // 16 16-bit columns
	program [256]Instruction
}

func NewFBOD() *FBOD {
	f := FBOD{
		acc:     0,
		mem:     [16][16]byte{},
		flags:   [16]byte{},
		screen:  [16]uint16{},
		program: [256]Instruction{},
	}

	// Set default F-Page Values
	f.mem[0xF][FPG_CLR_FG] = 0b0000 // Black
	f.mem[0xF][FPG_CLR_BG] = 0b0111 // Grey

	return &f
}

func (f *FBOD) ClearMem() {
	f.acc = 0
	f.mem = [16][16]byte{}
	f.screen = [16]uint16{}

	// Set default F-Page Values
	f.mem[0xF][FPG_CLR_FG] = 0b0000 // Black
	f.mem[0xF][FPG_CLR_BG] = 0b0111 // Grey
}

func (f *FBOD) ClearScreen() {
	f.screen = [16]uint16{}
}

func (f *FBOD) FlipPixel(x, y byte) {
	f.screen[y] ^= 1 << (15 - x)
}

func (f *FBOD) GetPixel(x, y byte) byte {
	return byte((f.screen[y] << x) % 2)
}

func (f *FBOD) SaveProgram(filename string) error {
	data := make([]byte, len(f.program)*2)
	for i := 0; i < len(f.program)*2; i += 2 {
		ins := f.program[i/2]
		data[i] = ins.ins
		data[i+1] = (ins.arg1 << 4) | ins.arg2
	}

	return ioutil.WriteFile(filename, data, FILE_MODE)
}

func (f *FBOD) LoadProgram(filename string) error {
	f.ClearMem()
	f.program = [256]Instruction{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i += 2 {
		f.program[i/2] = Instruction{
			ins:  data[i],
			arg1: (data[i+1] >> 4) % 16,
			arg2: data[i+1] % 16,
		}
	}
	return nil
}

// Reads through the program and indexes all flags
func (f *FBOD) ReadFlags() {
	for i, ins := range f.program {
		if ins.ins == ASM_FLG {
			f.flags[ins.arg1] = byte(i)
		}
	}
}

// Returns the index of the next instruction to perform
func (f *FBOD) PerformInstruction(progIndex byte) byte {
	nextIndex := progIndex + 1
	ins := f.program[progIndex]
	resArg := f.mem[ins.arg2][ins.arg1] // resolved memory argument

	switch ins.ins {

	case ASM_NOP:
		// Does Nothing

	case ASM_MVA:
		f.acc = resArg

	case ASM_MVM:
		f.mem[ins.arg2][ins.arg1] = f.acc

	case ASM_STA:
		f.acc = ins.arg1

	case ASM_DEC:
		f.acc--
		if f.acc < 0 {
			f.acc = 15
		}

	case ASM_INC:
		f.acc++
		if f.acc > 15 {
			f.acc = 0
		}

	case ASM_CLS:
		f.ClearScreen()

	case ASM_SHL:
		f.acc <<= 1
		f.acc %= 16 // Chop off shifted bits outside of nybl

	case ASM_SHR:
		f.acc >>= 1

	case ASM_RDP:
		f.acc = f.GetPixel(f.mem[0][ins.arg1], f.mem[0][ins.arg2])

	case ASM_FLP:
		f.FlipPixel(f.mem[0][ins.arg1], f.mem[0][ins.arg2])

	case ASM_FLG:
		// Does nothing; flags are read before execution

	case ASM_JMP:
		nextIndex = f.flags[resArg]

	case ASM_CEQ:
		if !(resArg == f.acc) {
			nextIndex++
		}

	case ASM_CGT:
		if !(resArg > f.acc) {
			nextIndex++
		}

	case ASM_CLT:
		if !(resArg < f.acc) {
			nextIndex++
		}

	}

	return nextIndex
}

func (f *FBOD) handleFPage() {
	f.mem[0xF][FPG_INPUT] = GetArrowsNybl()

	f.mem[0xF][FPG_RAND] = byte(rand.Intn(0xF))

	f.mem[0xF][FPG_DECF] = f.mem[0xF][FPG_DECF] - 1
	if f.mem[0xF][FPG_DECF] < 0 {
		f.mem[0xF][FPG_DECF] = 0xF
	}

	// f.mem[0xF][FPG_SAME] = f.mem[0xF][FPG_SAME] // Not sure what this location is supposed to do

	redirAddr := f.mem[0xF][FPG_REDIR_AD]
	redirPage := f.mem[0xF][FPG_REDIR_PG]
	f.mem[0xF][FPG_REDIR] = f.mem[redirPage][redirAddr]

	beepPitch := f.mem[0xF][FPG_BEEP_P]
	beepVol := f.mem[0xF][FPG_BEEP_V]
	beep(Pitch(beepPitch), beepVol)

	// Pitch and volume delta left out for now

	g_colourFG = fourBitColour(f.mem[0xF][FPG_CLR_FG])
	g_colourBG = fourBitColour(f.mem[0xF][FPG_CLR_BG])

	// Score is WIP
}
