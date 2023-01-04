package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

const (
	FILE_MODE = 0650
)

// 4BID-N Instructions
const (
	ASM_BRK  = 0x0 // Halt the program
	ASM_LDAI = 0x1 // Load immediate value to acc
	ASM_LDAM = 0x2 // Load memory value to acc
	ASM_STA  = 0x3 // Store acc to memory

	ASM_IDC = 0x4 // Increment/Decrement acc
	ASM_ADD = 0x5 // Add memory value to acc
	//ASM_001 = 0x6 //
	//ASM_SHL = 0x7 //

	ASM_NOT = 0x8 // Bitwise NOT
	ASM_ORA = 0x9 // Bitwise OR memory value and acc
	ASM_AND = 0xA // Bitwise AND memory value and acc
	ASM_SHF = 0xB // Bitwise shift (l/r & rot based on high bits)

	ASM_SLP  = 0xC // Sleeps for B seconds at A scale
	ASM_BNE  = 0xD // Skips B many instructions if acc does not equal A
	ASM_JMPI = 0xE // Jump to immediate program location
	ASM_JMPM = 0xF // Jump to memory jump vector
)

// 4BID-N F-Page Addresses
const (
	FPG_P1_DPAD = 0x0 // Player 1 Direction-Pad
	FPG_P1_BTNS = 0x1 // Player 1 Buttons
	FPG_P2_DPAD = 0x2 // Player 2 Direction-Pad
	FPG_P2_BTNS = 0x3 // Player 2 Buttons

	FPG_SCR_X   = 0x4 // Screen X Coord
	FPG_SCR_Y   = 0x5 // Screen Y Coord
	FPG_SCR_VAL = 0x6 // Screen Pixel Value
	FPG_SCR_OPT = 0x7 // Screen Options

	FPG_BEEP_VOL = 0x8 // Beeper Volume
	FPG_BEEP_PTC = 0x9 // Beeper Pitch
	FPG_BEEP_OPT = 0xA // Beeper reserved

	FPG_RAND = 0xB // Pseudo-Random Number

	FPG_DSK_H   = 0xC // High-nyble of disk address   \
	FPG_DSK_M   = 0xD // Middle-nyble of disk address  } 12-bit Address
	FPG_DSK_L   = 0xE // Low-nyble of disk address    /
	FPG_DSK_VAL = 0xF // Value of the selected disk nyble
)

type Instruction struct {
	ins  byte
	arg1 byte
	arg2 byte
}

type CPU struct {
	acc       byte
	mem       [16][16]byte
	flags     [16]byte   // List of program addresses
	screen    [16]uint16 // 16 16-bit columns
	program   [256]Instruction
	isRunning bool
}

func NewCPU() *CPU {
	f := CPU{
		acc:       0,
		mem:       [16][16]byte{},
		screen:    [16]uint16{},
		program:   [256]Instruction{},
		isRunning: true,
	}

	return &f
}

func (f *CPU) ClearMem() {
	f.acc = 0
	f.mem = [16][16]byte{}
	f.screen = [16]uint16{}
	f.isRunning = true
}

func (f *CPU) ClearScreen() {
	f.screen = [16]uint16{}
}

func (f *CPU) GetPixel(x, y byte) byte {
	return byte((f.screen[y] << x) % 2)
}

func (f *CPU) SaveProgram(filename string) error {
	data := make([]byte, len(f.program)*2)
	for i := 0; i < len(f.program)*2; i += 2 {
		ins := f.program[i/2]
		data[i] = ins.ins
		data[i+1] = (ins.arg1 << 4) | ins.arg2
	}

	return ioutil.WriteFile(filename, data, FILE_MODE)
}

func (f *CPU) LoadProgram(filename string) error {
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

// Returns the index of the next instruction to perform
func (f *CPU) PerformInstruction(progIndex byte) byte {
	if !f.isRunning {
		return progIndex
	}

	nextIndex := progIndex + 1
	ins := f.program[progIndex]
	resMem := f.mem[ins.arg2][ins.arg1] // resolved memory argument

	switch ins.ins {

	case ASM_BRK:
		f.isRunning = false

	case ASM_LDAI:
		f.acc = ins.arg1

	case ASM_LDAM:
		f.acc = resMem

	case ASM_STA:
		f.mem[ins.arg2][ins.arg1] = f.acc

		// Update Screen if screen-value changed
		if ins.arg2 == 0xF && ins.arg1 == FPG_SCR_VAL {
			f.updateScreenValue()
		}

		// Update Disk if disk-value changed
		if ins.arg2 == 0xF && ins.arg1 == FPG_DSK_VAL {
			err := WriteDisk(
				f.mem[0xF][FPG_DSK_H],
				f.mem[0xF][FPG_DSK_M],
				f.mem[0xF][FPG_DSK_L],
				f.acc,
			)
			if err != nil {
				ErrorPopup("Failed to read from disk")
				fmt.Printf("Failed to read from disk: %v\n", err)
			}
		}

	case ASM_IDC:
		f.acc += ins.arg1
		f.acc %= 16
		f.acc -= ins.arg2
		f.acc %= 16

	case ASM_ADD:
		f.acc += resMem
		f.acc %= 0xF

	case ASM_NOT:
		f.acc = ^f.acc

	case ASM_ORA:
		f.acc |= resMem

	case ASM_AND:
		f.acc &= resMem

	case ASM_SHF:
		if (ins.arg1>>3)%2 == 0 {
			if (ins.arg1>>2)%2 == 0 {
				f.acc <<= ins.arg2 % 4
				f.acc %= 0xF
			} else {
				f.acc >>= ins.arg2 % 4
			}
		} else {
			if (ins.arg1>>2)%2 == 0 {
				f.acc <<= ins.arg2 % 4
				f.acc |= (f.acc >> 4) % 0xF
			} else {
				f.acc <<= 4
				f.acc >>= ins.arg2 % 4
				f.acc |= (f.acc >> 4) % 0xF
			}
		}

	case ASM_BNE:
		if f.acc != ins.arg1 {
			nextIndex += ins.arg2
		}

	case ASM_SLP:
		scale := (ins.arg1 >> 1) % 8
		mul := math.Pow10(int(scale) - 4)
		length := (ins.arg1%2)<<4 | ins.arg2
		dur := time.Duration(mul * float64(length) * float64(time.Second.Nanoseconds()))
		//time.Sleep(dur)

		f.isRunning = false
		time.AfterFunc(dur, func() {
			f.isRunning = true
		})

	case ASM_JMPI:
		nextIndex = (ins.arg2 * 0xF) + ins.arg1

	case ASM_JMPM:
		addr := f.mem[0x0][ins.arg1]
		page := f.mem[0x0][ins.arg2]
		nextIndex = (page * 0xF) + addr

	}

	return nextIndex
}

func (f *CPU) updateScreenValue() {
	// Screen Updating
	x := f.mem[0xF][FPG_SCR_X]
	y := f.mem[0xF][FPG_SCR_Y]

	val := f.mem[0xF][FPG_SCR_VAL]
	opt := f.mem[0xF][FPG_SCR_OPT]
	switch opt % 4 {

	case 0b00:
		f.screen[y] ^= uint16((val % 2)) << (15 - x)

	case 0b01:
		f.screen[y] ^= uint16(val) << (15 - 3 - x)

	case 0b10:
		f.screen[y] ^= uint16(((val >> 3) % 2)) << (15 - x)
		f.screen[y+1] ^= uint16(((val >> 2) % 2)) << (15 - x)
		f.screen[y+2] ^= uint16(((val >> 1) % 2)) << (15 - x)
		f.screen[y+3] ^= uint16((val % 2)) << (15 - x)

	case 0b11:
		f.screen[y] ^= uint16(((val >> 2) % 4)) << ((15 - 1 - x) % 4)
		f.screen[y+1] ^= uint16((val % 2)) << ((15 - 1 - x) % 4)

	}
}

func (f *CPU) handleFPage() {
	dpad, btns := GetControlNybles(g_options.Controls.Player1Input)
	f.mem[0xF][FPG_P1_DPAD] = dpad
	f.mem[0xF][FPG_P1_BTNS] = btns

	dpad, btns = GetControlNybles(g_options.Controls.Player2Input)
	f.mem[0xF][FPG_P2_DPAD] = dpad
	f.mem[0xF][FPG_P2_BTNS] = btns

	soundVol := f.mem[0xF][FPG_BEEP_VOL]
	soundPitch := f.mem[0xF][FPG_BEEP_PTC]
	soundOpt := f.mem[0xF][FPG_BEEP_OPT]
	wav := (soundOpt >> 2) % 4
	oct := soundOpt % 4
	beep(Pitch(soundPitch), oct, wav, soundVol)

	f.mem[0xF][FPG_RAND] = byte(rand.Intn(0xF))

	dskVal, err := ReadDisk(
		f.mem[0xF][FPG_DSK_H],
		f.mem[0xF][FPG_DSK_M],
		f.mem[0xF][FPG_DSK_L],
	)
	if err != nil {
		ErrorPopup("Failed to read from disk")
		fmt.Printf("Failed to read from disk: %v\n", err)
	}
	f.mem[0xF][FPG_DSK_VAL] = dskVal
}

var OPCODE_STR = map[byte]string{
	0x0: "BRK",
	0x1: "LDA imm.",
	0x2: "LDA mem.",
	0x3: "STA",
	0x4: "IDC",
	0x5: "ADD",
	0x8: "NOT",
	0x9: "ORA",
	0xA: "AND",
	0xB: "SHF",
	0xC: "SLP",
	0xD: "BNE",
	0xE: "JMP imm.",
	0xF: "JMP mem.",
}
