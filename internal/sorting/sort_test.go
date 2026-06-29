package sorting

import (
	"testing"
	"time"
)

// testEntry is a simple Entry implementation for use in tests.
type testEntry struct {
	name    string
	ext     string
	size    int64
	modTime time.Time
}

func (e testEntry) GetName() string      { return e.name }
func (e testEntry) GetExt() string       { return e.ext }
func (e testEntry) GetSize() int64       { return e.size }
func (e testEntry) GetModTime() time.Time { return e.modTime }

// makeEntry is a convenience constructor.
func makeEntry(name, ext string, size int64, t time.Time) testEntry {
	return testEntry{name: name, ext: ext, size: size, modTime: t}
}

var (
	t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	t2 = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

// sign returns -1, 0, or 1 for negative, zero, or positive n.
func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	default:
		return 0
	}
}

func TestAlpha(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Entry
		wantSign int // expected sign of Alpha(a,b)
	}{
		{
			name:     "a before b",
			a:        makeEntry("apple", ".go", 0, t0),
			b:        makeEntry("banana", ".go", 0, t0),
			wantSign: -1,
		},
		{
			name:     "same name",
			a:        makeEntry("foo", ".go", 0, t0),
			b:        makeEntry("foo", ".go", 0, t0),
			wantSign: 0,
		},
		{
			name:     "hidden dot stripped",
			a:        makeEntry(".bar", "", 0, t0),
			b:        makeEntry("baz", "", 0, t0),
			wantSign: -1,
		},
		{
			name:     "case insensitive",
			a:        makeEntry("Zebra", "", 0, t0),
			b:        makeEntry("alpha", "", 0, t0),
			wantSign: 1,
		},
		{
			name:     "z after a",
			a:        makeEntry("zoo", ".txt", 0, t0),
			b:        makeEntry("ant", ".txt", 0, t0),
			wantSign: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sign(Alpha(tc.a, tc.b))
			if got != tc.wantSign {
				t.Errorf("Alpha(%q, %q) sign = %d, want %d",
					tc.a.GetName(), tc.b.GetName(), got, tc.wantSign)
			}
		})
	}
}

func TestBySize(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Entry
		wantSign int
	}{
		{
			name:     "larger first",
			a:        makeEntry("big", ".bin", 1000, t0),
			b:        makeEntry("small", ".bin", 10, t0),
			wantSign: -1,
		},
		{
			name:     "smaller second",
			a:        makeEntry("small", ".bin", 10, t0),
			b:        makeEntry("big", ".bin", 1000, t0),
			wantSign: 1,
		},
		{
			name:     "tie broken alphabetically",
			a:        makeEntry("apple", ".go", 100, t0),
			b:        makeEntry("banana", ".go", 100, t0),
			wantSign: -1,
		},
		{
			name:     "same size same name",
			a:        makeEntry("foo", ".go", 50, t0),
			b:        makeEntry("foo", ".go", 50, t0),
			wantSign: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sign(BySize(tc.a, tc.b))
			if got != tc.wantSign {
				t.Errorf("BySize(%q, %q) sign = %d, want %d",
					tc.a.GetName(), tc.b.GetName(), got, tc.wantSign)
			}
		})
	}
}

func TestByTime(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Entry
		wantSign int
	}{
		{
			name:     "newer first",
			a:        makeEntry("new", "", 0, t2),
			b:        makeEntry("old", "", 0, t0),
			wantSign: -1,
		},
		{
			name:     "older second",
			a:        makeEntry("old", "", 0, t0),
			b:        makeEntry("new", "", 0, t2),
			wantSign: 1,
		},
		{
			name:     "same time same name",
			a:        makeEntry("file", ".go", 0, t1),
			b:        makeEntry("file", ".go", 0, t1),
			wantSign: 0,
		},
		{
			name:     "same time tie broken alphabetically",
			a:        makeEntry("alpha", ".go", 0, t1),
			b:        makeEntry("zeta", ".go", 0, t1),
			wantSign: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sign(ByTime(tc.a, tc.b))
			if got != tc.wantSign {
				t.Errorf("ByTime(%q, %q) sign = %d, want %d",
					tc.a.GetName(), tc.b.GetName(), got, tc.wantSign)
			}
		})
	}
}

func TestByExt(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Entry
		wantSign int
	}{
		{
			name:     ".go before .txt",
			a:        makeEntry("file", ".go", 0, t0),
			b:        makeEntry("file", ".txt", 0, t0),
			wantSign: -1,
		},
		{
			name:     ".txt after .go",
			a:        makeEntry("file", ".txt", 0, t0),
			b:        makeEntry("file", ".go", 0, t0),
			wantSign: 1,
		},
		{
			name:     "same ext tie broken by name",
			a:        makeEntry("alpha", ".go", 0, t0),
			b:        makeEntry("zeta", ".go", 0, t0),
			wantSign: -1,
		},
		{
			name:     "case insensitive ext",
			a:        makeEntry("a", ".GO", 0, t0),
			b:        makeEntry("b", ".go", 0, t0),
			wantSign: -1, // same ext -> alpha tiebreak: a < b
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sign(ByExt(tc.a, tc.b))
			if got != tc.wantSign {
				t.Errorf("ByExt(%q%q, %q%q) sign = %d, want %d",
					tc.a.GetName(), tc.a.GetExt(),
					tc.b.GetName(), tc.b.GetExt(),
					got, tc.wantSign)
			}
		})
	}
}

func TestReversed(t *testing.T) {
	small := makeEntry("a", ".go", 1, t0)
	large := makeEntry("z", ".go", 100, t0)

	// Alpha normally: small < large → negative
	if sign(Alpha(small, large)) != -1 {
		t.Fatal("precondition: Alpha(small,large) should be negative")
	}

	rev := Reversed(Alpha)

	// Reversed should flip: small > large → positive
	if got := sign(rev(small, large)); got != 1 {
		t.Errorf("Reversed(Alpha)(small, large) sign = %d, want 1", got)
	}
	// And the inverse pair
	if got := sign(rev(large, small)); got != -1 {
		t.Errorf("Reversed(Alpha)(large, small) sign = %d, want -1", got)
	}
	// Equal pair should still be zero
	eq := makeEntry("x", ".go", 10, t0)
	if got := sign(rev(eq, eq)); got != 0 {
		t.Errorf("Reversed(Alpha)(eq, eq) sign = %d, want 0", got)
	}
}

func TestUnsorted(t *testing.T) {
	a := makeEntry("z", ".go", 100, t2)
	b := makeEntry("a", ".txt", 1, t0)
	if got := Unsorted(a, b); got != 0 {
		t.Errorf("Unsorted() = %d, want 0", got)
	}
}
