package themes

// Nord color reference: https://www.nordtheme.com/docs/colors-and-palettes
// Palette groups:
//   Polar Night: #2e3440 #3b4252 #434c5e #4c566a
//   Snow Storm:  #d8dee9 #e5e9f0 #eceff4
//   Frost:       #8fbcbb #88c0d0 #81a1c1 #5e81ac
//   Aurora:      #bf616a #d08770 #ebcb8b #a3be8c #b48ead

func init() { Register(nordTheme) }

var nordTheme = &Theme{
	Name: "nord",

	// Directories: Frost blue
	Dir:       Color{129, 161, 193}, // nord9  #81a1c1
	DirOpen:   Color{136, 192, 208}, // nord8  #88c0d0
	HiddenDir: Color{76, 86, 106},   // nord3  #4c566a – muted

	// Files
	File:       Color{216, 222, 233}, // nord4  #d8dee9
	HiddenFile: Color{76, 86, 106},   // nord3

	// Special types
	Executable: Color{163, 190, 140}, // aurora green  #a3be8c
	Symlink:    Color{136, 192, 208}, // nord8
	Pipe:       Color{208, 135, 112}, // aurora orange #d08770
	Socket:     Color{180, 142, 173}, // aurora purple #b48ead
	Special:    Color{235, 203, 139}, // aurora yellow #ebcb8b

	// Git
	GitUntracked: Color{163, 190, 140}, // green
	GitModified:  Color{235, 203, 139}, // yellow
	GitAdded:     Color{163, 190, 140}, // green
	GitDeleted:   Color{191, 97, 106},  // aurora red   #bf616a
	GitRenamed:   Color{129, 161, 193}, // nord9
	GitConflict:  Color{191, 97, 106},  // red

	// Long listing
	Permissions: Color{143, 188, 187}, // nord7  #8fbcbb – teal
	Owner:       Color{136, 192, 208}, // nord8
	Group:       Color{94, 129, 172},  // nord10 #5e81ac
	SizeUnit:    Color{235, 203, 139}, // yellow
	SizeNum:     Color{163, 190, 140}, // green
	Date:        Color{136, 192, 208}, // nord8

	// Block size
	BlockSize: Color{143, 188, 187}, // nord7
}
