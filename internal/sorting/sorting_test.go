package sorting

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/the-mayankjha/modern-ls/internal/filesystem"
)

func generateDummyEntries(n int) []*filesystem.Entry {
	entries := make([]*filesystem.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = &filesystem.Entry{
			Name:     fmt.Sprintf("file_%d", i),
			Ext:      ".txt",
			FullName: fmt.Sprintf("file_%d.txt", i),
			Size:     rand.Int63n(1000000),
			ModTime:  time.Now().Add(-time.Duration(rand.Intn(10000)) * time.Hour),
		}
	}
	return entries
}

func BenchmarkSortAlpha(b *testing.B) {
	entries := generateDummyEntries(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clone entries to avoid sorting an already sorted slice
		clone := make([]*filesystem.Entry, len(entries))
		copy(clone, entries)
		
		sort.SliceStable(clone, func(i, j int) bool {
			return Alpha(clone[i], clone[j]) < 0
		})
	}
}

func BenchmarkSortSize(b *testing.B) {
	entries := generateDummyEntries(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := make([]*filesystem.Entry, len(entries))
		copy(clone, entries)
		
		sort.SliceStable(clone, func(i, j int) bool {
			return BySize(clone[i], clone[j]) < 0
		})
	}
}

func BenchmarkSortTime(b *testing.B) {
	entries := generateDummyEntries(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := make([]*filesystem.Entry, len(entries))
		copy(clone, entries)
		
		sort.SliceStable(clone, func(i, j int) bool {
			return ByTime(clone[i], clone[j]) < 0
		})
	}
}
