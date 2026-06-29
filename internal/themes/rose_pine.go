package themes

// Rosé Pine color reference: https://rosepinetheme.com/palette/
// Main (dark) variant:
//   base    #191724  surface #1f1d2e  overlay #26233a
//   muted   #6e6a86  subtle  #908caa  text    #e0def4
//   love    #eb6f92  gold    #f6c177  rose    #ebbcba
//   pine    #31748f  foam    #9ccfd8  iris    #c4a7e7
//   hl_low  #21202e  hl_med  #403d52  hl_high #524f67

func init() { Register(rosePineTheme) }

var rosePineTheme = &Theme{
	Name: "rosepine",

	// Directories: pine (teal-blue)
	Dir:       Color{49, 116, 143},  // pine  #31748f
	DirOpen:   Color{156, 207, 216}, // foam  #9ccfd8
	HiddenDir: Color{110, 106, 134}, // muted #6e6a86

	// Files
	File:       Color{224, 222, 244}, // text  #e0def4
	HiddenFile: Color{144, 140, 170}, // subtle #908caa

	// Special types
	Executable: Color{156, 207, 216}, // foam
	Symlink:    Color{196, 167, 231}, // iris  #c4a7e7
	Pipe:       Color{246, 193, 119}, // gold  #f6c177
	Socket:     Color{235, 111, 146}, // love  #eb6f92
	Special:    Color{235, 188, 186}, // rose  #ebbcba

	// Git
	GitUntracked: Color{156, 207, 216}, // foam
	GitModified:  Color{246, 193, 119}, // gold
	GitAdded:     Color{156, 207, 216}, // foam
	GitDeleted:   Color{235, 111, 146}, // love
	GitRenamed:   Color{49, 116, 143},  // pine
	GitConflict:  Color{235, 111, 146}, // love

	// Long listing
	Permissions: Color{196, 167, 231}, // iris
	Owner:       Color{156, 207, 216}, // foam
	Group:       Color{49, 116, 143},  // pine
	SizeUnit:    Color{235, 188, 186}, // rose
	SizeNum:     Color{246, 193, 119}, // gold
	Date:        Color{156, 207, 216}, // foam

	// Block size
	BlockSize: Color{196, 167, 231}, // iris
}
