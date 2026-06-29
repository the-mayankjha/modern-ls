// Package sorting provides sort strategies for directory entries.
// Each Strategy is a standalone comparison function with the same signature
// as cmp.Compare, so strategies can be composed with standard library helpers
// such as slices.SortFunc.
package sorting

import (
	"strings"
	"time"
)

// Entry is the minimal interface required by all sort strategies.
// Implementations are provided by the filesystem package.
type Entry interface {
	GetName() string
	GetExt() string
	GetSize() int64
	GetModTime() time.Time
}

// Strategy is a comparison function returning a negative integer when a < b,
// zero when a == b, and a positive integer when a > b.
type Strategy func(a, b Entry) int

// mainCmp is the canonical alphabetic tiebreaker used across strategies.
// It strips leading dots (hidden files display beside their non-hidden
// counterparts) and compares case-insensitively. The special entries "." and
// ".." are compared verbatim so they always sort at the top.
func mainCmp(a, b string) int {
	clean := func(s string) string {
		switch s {
		case ".", "..":
			return s
		default:
			return strings.TrimPrefix(s, ".")
		}
	}
	la := strings.ToLower(clean(a))
	lb := strings.ToLower(clean(b))
	if la < lb {
		return -1
	}
	if la > lb {
		return 1
	}
	return 0
}

// Alpha sorts entries alphabetically by concatenated name+ext, ignoring
// leading dots and case.
func Alpha(a, b Entry) int {
	return mainCmp(a.GetName()+a.GetExt(), b.GetName()+b.GetExt())
}

// BySize sorts entries largest-first; ties are broken alphabetically.
func BySize(a, b Entry) int {
	if a.GetSize() > b.GetSize() {
		return -1
	}
	if a.GetSize() < b.GetSize() {
		return 1
	}
	return mainCmp(a.GetName()+a.GetExt(), b.GetName()+b.GetExt())
}

// ByTime sorts entries newest-first (most recently modified first).
// Ties are broken alphabetically.
func ByTime(a, b Entry) int {
	at, bt := a.GetModTime(), b.GetModTime()
	if at.After(bt) {
		return -1
	}
	if bt.After(at) {
		return 1
	}
	return mainCmp(a.GetName()+a.GetExt(), b.GetName()+b.GetExt())
}

// ByExt sorts alphabetically by file extension (case-insensitive);
// ties are broken by name.
func ByExt(a, b Entry) int {
	ae := strings.ToLower(a.GetExt())
	be := strings.ToLower(b.GetExt())
	if ae < be {
		return -1
	}
	if ae > be {
		return 1
	}
	return mainCmp(a.GetName()+a.GetExt(), b.GetName()+b.GetExt())
}

// Natural sorts by name+ext using byte-order string comparison, which
// approximates "natural" sort for ASCII-majority filenames without requiring
// a full numeric-aware implementation.
func Natural(a, b Entry) int {
	an := a.GetName() + a.GetExt()
	bn := b.GetName() + b.GetExt()
	if an < bn {
		return -1
	}
	if an > bn {
		return 1
	}
	return 0
}

// Unsorted preserves the directory-read order returned by the OS.
func Unsorted(_, _ Entry) int {
	return 0
}

// Reversed wraps any Strategy and inverts its result, turning an
// ascending sort into a descending one.
func Reversed(s Strategy) Strategy {
	return func(a, b Entry) int {
		return s(b, a)
	}
}
