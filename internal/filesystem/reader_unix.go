//go:build !windows

package filesystem

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// ownerGroup returns the owner username and group name for the file.
// Results are cached to avoid repeated lookups for the same UIDs/GIDs.
var userCache = map[uint32]string{}
var groupCache = map[uint32]string{}

func ownerGroup(fi fs.FileInfo) (owner, group string) {
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", ""
	}

	uid := stat.Uid
	if name, cached := userCache[uid]; cached {
		owner = name
	} else {
		if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
			owner = u.Username
		}
		userCache[uid] = owner
	}

	gid := stat.Gid
	if name, cached := groupCache[gid]; cached {
		group = name
	} else {
		if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
			group = g.Name
		}
		groupCache[gid] = group
	}
	return
}

// fileBlocks returns the number of 512-byte blocks allocated to the file.
func fileBlocks(fi fs.FileInfo) int64 {
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return stat.Blocks
	}
	return 0
}

// dirBlocks fills in the block count for a directory entry.
func dirBlocks(e *Entry, fi fs.FileInfo) {
	e.Blocks = fileBlocks(fi)
}

// getEntryStat fills in the syscall-specific fields for the given path.
func getEntryStat(path string) (int64, string, string) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, "", ""
	}
	blocks := fileBlocks(fi)
	owner, group := ownerGroup(fi)
	return blocks, owner, group
}
