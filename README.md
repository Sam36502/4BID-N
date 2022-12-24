# 4BID-N
A 4-Bit Fantasy Console heavily inspired by
the [4BoD](https://puarsliburf.itch.io/4bod-fantaly-console) console by Puarsliburf.

It's pronounced like "forbidden" and stands for *"4-Bit I DunNo"*
(i.e. it doesn't really stand for anything)

## Download
 - ### *None yet...*

## Documentation
 - ### *None yet...*

### Instructions
| Binary | Hex | Opcode | Arguments | Description                                                  |
|--------|-----|--------|-----------|--------------------------------------------------------------|
| 0000   | 0   | BRK    | - -       | Halts the program                                            |
| 0001   | 1   | LDA    | # -       | Loads an immediate value into acc                            |
| 0010   | 2   | LDA    | $ $       | Loads the value at address AB into acc                       |
| 0011   | 3   | STA    | $ $       | Stores the value of acc to address AB                        |
| 0100   | 4   | IDC    | # #       | Increments the acc A many times; Decrements B many times     |
| 0101   | 5   | ADD    | $ $       | Adds the value at AB to acc                                  |
| 0110   | 6   |        | - -       | Not assigned                                                 |
| 0111   | 7   |        | - -       | Not assigned                                                 |
| 1000   | 8   | NOT    | - -       | Bitwise NOT (inverts) the acc                                |
| 1001   | 9   | ORA    | $ $       | Bitwise OR the acc with the value at address AB              |
| 1010   | A   | AND    | $ $       | Bitwise AND the acc with the value at address AB             |
| 1011   | B   | SHF    | # -       | Bitwise shifts the acc left or right (see Shifting below)    |
| 1100   | C   | SLP    | # #       | Waits for the B many seconds at scale A (see Sleeping below) |
| 1101   | D   | BNE    | # #       | Skips the next B many instructions if acc doesn't equal A    |
| 1110   | E   | JMP    | # #       | Jumps to the position at program instruction A page B        |
| 1111   | F   | JMP    | $ $       | Jumps to the instruction pointed to at memory location A & B |

#### Arguments:
The arguments column shows how arguments A & B are interpreted according to following symbols
 - `-`: No argument required
 - `#`: Immediate (literal) value
 - `$`: Memory Value; the value at the given memory address

Addresses are all 8-bit and use both arguments; A is the address within the page and B is
the page number. The `JMP` instruction follows the same convention but in program memory.

#### Bit-Shifting:
The `SHF` instruction shifts the bits of the accumulator by a certain amount. Whether it shifts
left or right and whether it wraps around (circular shift) are determined by the high 2 bits
and the amount to shift by the lower 2 bits:

    8   0: No wrap; 1: circular shift
    4   0: Left shift; 1: Right shift
    2   \
    1   / amount to shift: 0-3
    
    E.g.:
    
    LDA     #b1000  ; Load 1000 into acc
    SHF     #b0110  ; Non-circular, right-shift by 2
                    ; acc = b0010
                    
    LDA     #b1000  ; Load 1000 into acc
    SHF     #b1011  ; Circular, left-shift by 3
                    ; acc = b0100

#### Sleeping
To maximise the useful time intervals that can be waited with this instruction,
the arguments have been split into a 5-bit number and a 3-bit scale. The scale
sets what order of magnitude we're using and consists of the upper 3 bits of argument B:

| Bits | Nr. | Multiplier | Digits |
|------|-----|------------|--------|
| 000  | 0   | 0.0001     | 0.00XX |
| 001  | 1   | 0.001      | 0.0XX  |
| 010  | 2   | 0.01       | 0.XX   |
| 011  | 3   | 0.1        | X.X    |
| 100  | 4   | 1          | XX     |
| 101  | 5   | 10         | XX0    |
| 110  | 6   | 100        | XX00   |
| 111  | 7   | 1000       | XX000  |

Then argument A plus the lowest bit of B (as the 16s place) combine to say how many
seconds to wait. E.g.:

    SLP     #b1000  #b0101  ; Waits 5 seconds
    SLP     #b1001  #b1111  ; Waits 31 seconds
    SLP     #b0110  #b0111  ; Waits 500 milliseconds

#### Jumping
There are two variants of the `JMP` instruction, the first (`0xE`) jumps to a direct program address
with instruction number A and program page B.

The second (`0xF`) jumps to the program address same as the first jump but the program address is
taken from the zero page. E.g.:

    LDA     #$A         ; Store 0xA at memory location 1 on the zero-page
    STA     $1      $0
    
    LDA     #$B         ; Store 0xB at memory location 2 on the zero-page
    STA     $2      $0
    
    JMP     $2      $1    ; Resolves jump vector from addresses 0x02 -> 0xB and 0x01 -> 0xA
                        ; Jumps to instruction 0xB on program page 0xA
    
### F-Page
The F-Page (memory page 15) is reserved for special hardware access such as the input states and
screen access. These special addresses are as follows:

| Binary | Hex | Name                    | Description                                     |
|--------|-----|-------------------------|-------------------------------------------------|
| 0000   | 0   | Player 1 D-Pad          | The state of Player 1's D-Pad (DURL)            |
| 0001   | 1   | Player 1 Buttons        | The state of Player 1's Buttons                 |
| 0010   | 2   | Player 2 D-Pad          | The state of Player 2's D-Pad (DURL)            |
| 0011   | 3   | Player 2 Buttons        | The state of Player 2's Buttons                 |
| 0100   | 4   | Screen X                | Horizontal position on the screen to access     |
| 0101   | 5   | Screen Y                | Horizontal position on the screen to access     |
| 0110   | 6   | Pixel Value             | The value of the pixel at the selected position |
| 0111   | 7   | Screen Options          | Options for how to operate the screen           |
| 1000   | 8   | Sound Volume            | Sets the volume of the beeper                   |
| 1001   | 9   | Sound Pitch             | Sets the pitch of the beeper (See table below)  |
| 1010   | A   | Sound Options           | Options for how to operate the sound            |
| 1011   | B   | Random Value            | A pseudo-random number                          |
| 1100   | C   | Disk Address N2         | The high-nyble of the disk address              |
| 1101   | D   | Disk Address N1         | The middle-nyble of the disk address            |
| 1110   | E   | Disk Address N0         | The low-nyble of the disk address               |
| 1111   | F   | Disk Data               | The data at the provided disk address           |

#### Controllers
The console accepts input in the form of 4-directional buttons and 4 regular buttons, the
state of which is stored in a pair of nybles for each player. (Player 1 and 2's controllers
are simply mapped to different keys on a keyboard)

Directional Pad (D-Pad)
| Place | Digit | Key   |
|-------|-------|-------|
| 1s    | 000X  | Left  |
| 2s    | 00X0  | Right |
| 4s    | 0X00  | Up    |
| 8s    | X000  | Down  |

#### Screen Control
The X & Y position variables select a position on the screen which can then be read or written
from/to the Pixel Value memory address.
Extra options/operations can be activated with address `0x7`:
The highest bit (8s-place) sets whether to clear the screen and the
next higest bit (4s-place) sets whether to invert the screen.

**Lower 2 bits:**
| Binary | Function                   |
|--------|----------------------------|
| 00     | Pixel-Mode addressing      |
| 01     | Horizontal-Mode addressing |
| 10     | Vertical-Mode addressing   |
| 11     | Square-Mode addressing     |

These addressing modes change what section of the screen is selected by the X & Y coordinates.

By default, it's in pixel-mode, which only selects a single pixel at the given position and
sets the value to either 1 or 0 depending if it's on or off.

The other modes select a whole nyble (4-bits) of pixels in different arrangements.
Horizontal and Vertical select a line of 4 pixels starting at the given position and moving
right or down respectively:

    X=3, Y=2
    
    -------...
    -------...
    ---XHHH...
    ---V---...
    ---V---...
    ---V---...
    ..........

The bits of the nyble selected start with the most-significant-bit at the selected
position.

Square-Mode selects a 2x2 square with its top-left corner at the given position:
    
    X=3, Y=2
    
    -------...
    -------...  Bits:
    ---XS--...      84
    ---SS--...      21
    -------...
    -------...
    ..........

The order of these operations is that the screen is first cleared and/or inverted,
then the pixel value address (`0x6`) is updated with the current value on the screen,
before finally - after the next program instruction has run - the value at the pixel value address (`0x6`) is written to the screen

### Sound
The console has basic audio output which can be set to a certain pitch and volume and
plays continuously. The volume is simply a range from zero/off up to a maximum of the
program's master volume (as defined in `options.json`). The pitch is one according to the following table:
(The octave is set by sound options below)

| Binary | Hex | Note |
|--------|-----|------|
| 0000   | 0   | A    |
| 0001   | 1   | Bb   |
| 0010   | 2   | B    |
| 0011   | 3   | C    |
| 0100   | 4   | Cs   |
| 0101   | 5   | D    |
| 0110   | 6   | Ds   |
| 0111   | 7   | E    |
| 1000   | 8   | F    |
| 1001   | 9   | Fs   |
| 1010   | A   | G    |
| 1011   | B   | Gs   |
| 1100   | C   | A    |
| 1101   | D   | Bb   |
| 1110   | E   | B    |
| 1111   | F   | C    |

#### Sound Options
The sound options let you select which waveform and octave you're working in.

The upper 2 bits of this nyble select the waveform as follows
| bits | nr | Wave     |
|------|----|----------|
| 00   | 0  | Square   |
| 01   | 1  | Triangle |
| 10   | 2  | Sawtooth |
| 11   | 3  | Noise    |

The lower 2 bits select the octaves as follows:
| bits | nr | Octave | Lowest Note | Highest Note |
|------|----|--------|-------------|--------------|
| 00   | 0  | 2      | A1          | C3           |
| 01   | 1  | 3      | A2          | C4           |
| 10   | 2  | 4      | A3          | C5           |
| 11   | 3  | 5      | A4          | C6           |
    
## Options
If you want to customise the interface, you can do so with the included `options.json` file.
Options include:
| JSON Key         | Description                                                                                             |
|------------------|---------------------------------------------------------------------------------------------------------|
| `splash_millis`  | How many milliseconds to spend on the splash screen                                                     |
| `pixel_size`     | Side length of the 4BOD pixels in real pixels                                                           |
| `target_fps`     | What framerate to limit the program to. Helps to see what's actually happening (set to -1 for no limit) |
| `color_fg`       | The colour of foreground pixels                                                                         |
| `color_bg`       | The colour of background                                                                                |
| `color_overlay`  | The colour of the editor overlay                                                                        |
| `old_menu`       | Whether to use the old menu images                                                                      |
| `editor_overlay` | Whether to have the editor overlay on by default (Should save if turned off with it on)                 |
| `debug_keycodes` | Whether to display the last pressed keycode on the screen (helpful for changing controls)               |
| `controls`       | A list of various keys and their keycodes (see `debug_keycodes`)                                        |

## Changing Keyboard Inputs
The easiest way to change which keys do what is to set the `debug_keycodes` option by setting it to `true`
and starting the machine. Then, you can press the keys you want each thing to do and write down what the
keycode is. After that, you can change the respective `kc_...` settings in the controls part of `options.json`.
You can also use `0` to unbind the key.

