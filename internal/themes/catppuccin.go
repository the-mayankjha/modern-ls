package themes

// Catppuccin Mocha color reference: https://github.com/catppuccin/catppuccin
// Base palette (Mocha variant):
//   Rosewater #f5e0dc  Flamingo #f2cdcd  Pink #f5c2e7  Mauve #cba6f7
//   Red #f38ba8        Maroon  #eba0ac   Peach #fab387  Yellow #f9e2af
//   Green #a6e3a1      Teal    #94e2d5   Sky   #89dceb  Sapphire #74c7ec
//   Blue #89b4fa       Lavender #b4befe  Text  #cdd6f4  Subtext1 #bac2de
//   Overlay2 #9399b2   Surface1 #45475a  Base  #1e1e2e

func init() { Register(catppuccinTheme) }

var catppuccinTheme = &Theme{
	Name: "catppuccin",

	// Directories: Sapphire
	Dir:       Color{116, 199, 236},
	DirOpen:   Color{137, 220, 235}, // Sky – slightly lighter when expanded
	HiddenDir: Color{147, 153, 178}, // Overlay2 – muted

	// Files: Text / Subtext
	File:       Color{205, 214, 244}, // Text
	HiddenFile: Color{186, 194, 222}, // Subtext1

	// Special types
	Executable: Color{166, 227, 161}, // Green
	Symlink:    Color{137, 180, 250}, // Blue
	Pipe:       Color{250, 179, 135}, // Peach
	Socket:     Color{203, 166, 247}, // Mauve
	Special:    Color{249, 226, 175}, // Yellow

	// Git status
	GitUntracked: Color{166, 227, 161}, // Green
	GitModified:  Color{250, 179, 135}, // Peach
	GitAdded:     Color{166, 227, 161}, // Green
	GitDeleted:   Color{243, 139, 168}, // Red
	GitRenamed:   Color{116, 199, 236}, // Sapphire
	GitConflict:  Color{243, 139, 168}, // Red

	// Long listing
	Permissions: Color{148, 226, 213}, // Teal
	Owner:       Color{137, 180, 250}, // Blue
	Group:       Color{116, 199, 236}, // Sapphire
	SizeUnit:    Color{203, 166, 247}, // Mauve
	SizeNum:     Color{249, 226, 175}, // Yellow
	Date:        Color{148, 226, 213}, // Teal

	// Block size
	BlockSize: Color{137, 220, 235}, // Sky
}
