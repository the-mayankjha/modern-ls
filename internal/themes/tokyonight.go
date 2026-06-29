package themes

// Tokyo Night color reference: https://github.com/folke/tokyonight.nvim
// Storm variant palette (hex → RGB):
//   bg      #24283b   bg_dark  #1f2335   bg_float #1f2335
//   blue    #7aa2f7   blue1    #2ac3de   blue2    #0db9d7
//   cyan    #7dcfff   green    #9ece6a   green1   #73daca
//   magenta #bb9af7   orange   #ff9e64   red      #f7768e
//   red1    #db4b4b   teal     #1abc9c   yellow   #e0af68
//   fg      #c0caf5   comment  #565f89

func init() { Register(tokyoNightTheme) }

var tokyoNightTheme = &Theme{
	Name: "tokyonight",

	// Directories: blue
	Dir:       Color{122, 162, 247}, // blue
	DirOpen:   Color{125, 207, 255}, // cyan
	HiddenDir: Color{86, 95, 137},   // comment – muted

	// Files
	File:       Color{192, 202, 245}, // fg
	HiddenFile: Color{86, 95, 137},   // comment

	// Special types
	Executable: Color{158, 206, 106}, // green
	Symlink:    Color{42, 195, 222},  // blue1
	Pipe:       Color{255, 158, 100}, // orange
	Socket:     Color{187, 154, 247}, // magenta
	Special:    Color{224, 175, 104}, // yellow

	// Git
	GitUntracked: Color{115, 218, 202}, // green1
	GitModified:  Color{224, 175, 104}, // yellow
	GitAdded:     Color{158, 206, 106}, // green
	GitDeleted:   Color{247, 118, 142}, // red
	GitRenamed:   Color{122, 162, 247}, // blue
	GitConflict:  Color{219, 75, 75},   // red1

	// Long listing
	Permissions: Color{115, 218, 202}, // green1
	Owner:       Color{122, 162, 247}, // blue
	Group:       Color{42, 195, 222},  // blue1
	SizeUnit:    Color{187, 154, 247}, // magenta
	SizeNum:     Color{224, 175, 104}, // yellow
	Date:        Color{42, 195, 222},  // blue1

	// Block size
	BlockSize: Color{115, 218, 202}, // green1
}
