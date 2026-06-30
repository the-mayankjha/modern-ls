package themes

// Gruvbox Dark color reference: https://github.com/morhetz/gruvbox
// Hard-contrast dark palette:
//   bg0    #282828  bg1    #3c3836  bg2    #504945
//   fg     #ebdbb2  fg1    #d5c4a1  fg4    #a89984
//   red    #cc241d  red_b  #fb4934  green  #98971a  green_b #b8bb26
//   yellow #d79921  yel_b  #fabd2f  blue   #458588  blue_b  #83a598
//   purple #b16286  purp_b #d3869b  aqua   #689d6a  aqua_b  #8ec07c
//   orange #d65d0e  ora_b  #fe8019  gray   #a89984

func init() { Register(gruvboxTheme) }

var gruvboxTheme = &Theme{
	Name: "gruvbox",

	// Directories: bright yellow
	Dir:       Color{250, 189, 47},  // yellow_bright
	DirOpen:   Color{254, 128, 25},  // orange_bright
	HiddenDir: Color{168, 153, 132}, // gray

	// Files
	File:       Color{235, 219, 178}, // fg
	HiddenFile: Color{168, 153, 132}, // gray

	// Special types
	Executable: Color{184, 187, 38},  // green_bright
	Symlink:    Color{131, 165, 152}, // blue_bright
	Pipe:       Color{254, 128, 25},  // orange_bright
	Socket:     Color{211, 134, 155}, // purple_bright
	Special:    Color{142, 192, 124}, // aqua_bright

	// Git
	GitUntracked: Color{152, 151, 26},  // green
	GitModified:  Color{215, 153, 33},  // yellow
	GitAdded:     Color{184, 187, 38},  // green_bright
	GitDeleted:   Color{251, 73, 52},   // red_bright
	GitRenamed:   Color{131, 165, 152}, // blue_bright
	GitConflict:  Color{204, 36, 29},   // red

	// Long listing
	Permissions: Color{142, 192, 124}, // aqua_bright
	Owner:       Color{131, 165, 152}, // blue_bright
	Group:       Color{104, 157, 106}, // aqua
	SizeUnit:    Color{250, 189, 47},  // yellow_bright
	SizeNum:     Color{215, 153, 33},  // yellow
	Date:        Color{69, 133, 136},  // blue

	// Block size
	BlockSize: Color{104, 157, 106}, // aqua
}
