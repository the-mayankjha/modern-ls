package themes

func init() { Register(defaultTheme) }

// defaultTheme reproduces the original logo-ls color palette.
// Colors are sampled from the Material Design palette used by the
// original project.
var defaultTheme = &Theme{
	Name: "default",

	// Directories: amber-ish yellow
	Dir:       Color{224, 177, 77},
	DirOpen:   Color{224, 177, 77},
	HiddenDir: Color{224, 177, 77},

	// Files: steel blue
	File:       Color{65, 129, 190},
	HiddenFile: Color{65, 129, 190},

	// Special file types
	Executable: Color{76, 175, 80},  // green
	Symlink:    Color{66, 165, 245}, // light blue
	Pipe:       Color{250, 111, 66}, // orange
	Socket:     Color{66, 165, 245}, // light blue
	Special:    Color{175, 180, 43}, // yellow-green

	// Git status
	GitUntracked: Color{55, 183, 21},
	GitModified:  Color{192, 154, 107},
	GitAdded:     Color{55, 183, 21},
	GitDeleted:   Color{229, 77, 58},
	GitRenamed:   Color{66, 165, 245},
	GitConflict:  Color{229, 77, 58},

	// Long listing
	Permissions: Color{175, 180, 43},
	Owner:       Color{66, 165, 245},
	Group:       Color{66, 165, 245},
	SizeUnit:    Color{255, 213, 79},
	SizeNum:     Color{175, 180, 43},
	Date:        Color{66, 165, 245},

	// Block size
	BlockSize: Color{175, 180, 43},
}

// Default returns the built-in default theme. It is always valid and never nil.
func Default() *Theme { return defaultTheme }
