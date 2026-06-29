//go:build windows

package filesystem

import "io/fs"

func ownerGroup(_ fs.FileInfo) (owner, group string) { return "", "" }
func fileBlocks(_ fs.FileInfo) int64                 { return 0 }
func dirBlocks(_ *Entry, _ fs.FileInfo)               {}
