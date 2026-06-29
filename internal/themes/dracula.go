package themes

// Dracula color reference: https://draculatheme.com/contribute
// Official palette:
//   Background #282a36  Current #44475a  Foreground #f8f8f2
//   Comment    #6272a4  Cyan    #8be9fd  Green      #50fa7b
//   Orange     #ffb86c  Pink    #ff79c6  Purple     #bd93f9
//   Red        #ff5555  Yellow  #f1fa8c

func init() { Register(draculaTheme) }

var draculaTheme = &Theme{
	Name: "dracula",

	// Directories: purple
	Dir:       Color{189, 147, 249}, // purple
	DirOpen:   Color{139, 233, 253}, // cyan
	HiddenDir: Color{98, 114, 164},  // comment – muted

	// Files
	File:       Color{248, 248, 242}, // foreground
	HiddenFile: Color{98, 114, 164},  // comment

	// Special types
	Executable: Color{80, 250, 123},  // green
	Symlink:    Color{139, 233, 253}, // cyan
	Pipe:       Color{255, 184, 108}, // orange
	Socket:     Color{255, 121, 198}, // pink
	Special:    Color{241, 250, 140}, // yellow

	// Git
	GitUntracked: Color{80, 250, 123},  // green
	GitModified:  Color{255, 184, 108}, // orange
	GitAdded:     Color{80, 250, 123},  // green
	GitDeleted:   Color{255, 85, 85},   // red
	GitRenamed:   Color{139, 233, 253}, // cyan
	GitConflict:  Color{255, 85, 85},   // red

	// Long listing
	Permissions: Color{241, 250, 140}, // yellow
	Owner:       Color{139, 233, 253}, // cyan
	Group:       Color{189, 147, 249}, // purple
	SizeUnit:    Color{255, 121, 198}, // pink
	SizeNum:     Color{255, 184, 108}, // orange
	Date:        Color{139, 233, 253}, // cyan

	// Block size
	BlockSize: Color{80, 250, 123}, // green
}
