package disk

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// MountInfo contains metrics for a single filesystem mount
type MountInfo struct {
	Device     string  // Device name (e.g., /dev/nvme0n1p4)
	MountPoint string  // Mount point path (e.g., /home)
	FSType     string  // Filesystem type (e.g., ext4, btrfs, vfat)
	Options    string  // Mount options
	Total      uint64  // Total capacity in bytes
	Used       uint64  // Used bytes
	Avail      uint64  // Available bytes for unprivileged user
	Free       uint64  // Total free bytes
	UsagePerc  float64 // Usage percentage (0.0 - 100.0)
	IsReadOnly bool    // Read-only filesystem
	IsSpecial  bool    // Pseudo or virtual filesystem
	Inodes     uint64  // Total inodes
	InodesUsed uint64  // Used inodes
	InodesFree uint64  // Free inodes
}

// GetMounts reads mount entries from /proc/mounts or /etc/mtab and queries Statfs
func GetMounts() ([]MountInfo, error) {
	mountFilePath := "/proc/mounts"
	if _, err := os.Stat(mountFilePath); err != nil {
		mountFilePath = "/etc/mtab"
	}

	file, err := os.Open(mountFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mounts []MountInfo
	scanner := bufio.NewScanner(file)

	// Keep track of visited mount points to avoid duplicates
	seenMounts := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		device := unescapeMount(fields[0])
		mountPoint := unescapeMount(fields[1])
		fsType := fields[2]
		options := fields[3]

		if seenMounts[mountPoint] {
			continue
		}
		seenMounts[mountPoint] = true

		info, err := queryMount(device, mountPoint, fsType, options)
		if err != nil {
			// If statfs fails (e.g., permission denied or unavailable mount), still record basic info
			info = MountInfo{
				Device:     device,
				MountPoint: mountPoint,
				FSType:     fsType,
				Options:    options,
				IsSpecial:  isSpecialFSType(fsType) || !strings.HasPrefix(device, "/"),
			}
		}
		mounts = append(mounts, info)
	}

	return mounts, scanner.Err()
}

func queryMount(device, mountPoint, fsType, options string) (MountInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(mountPoint, &stat)
	if err != nil {
		return MountInfo{}, err
	}

	// Fundamental block size (use Frsize if available, else Bsize)
	bsize := uint64(stat.Bsize)
	if stat.Frsize > 0 {
		bsize = uint64(stat.Frsize)
	}

	total := stat.Blocks * bsize
	free := stat.Bfree * bsize
	avail := stat.Bavail * bsize

	var used uint64
	if stat.Blocks >= stat.Bfree {
		used = (stat.Blocks - stat.Bfree) * bsize
	} else {
		used = 0
	}

	// Calculate percentage: standard df formula based on available space for users
	var usagePerc float64
	if used+avail > 0 {
		usagePerc = (float64(used) / float64(used+avail)) * 100.0
	} else if total > 0 {
		usagePerc = (float64(used) / float64(total)) * 100.0
	}

	isRO := strings.Contains(options, "ro") || stat.Flags&syscall.MS_RDONLY != 0
	isSpecial := isSpecialFSType(fsType) || !strings.HasPrefix(device, "/") || total == 0

	var inodesTotal, inodesFree, inodesUsed uint64
	inodesTotal = stat.Files
	inodesFree = stat.Ffree
	if stat.Files >= stat.Ffree {
		inodesUsed = stat.Files - stat.Ffree
	}

	// Resolve symlink for LVM / mapper if possible
	resolvedDevice := device
	if strings.HasPrefix(device, "/dev/mapper/") {
		if realPath, err := filepath.EvalSymlinks(device); err == nil {
			resolvedDevice = realPath
		}
	}

	return MountInfo{
		Device:      resolvedDevice,
		MountPoint:  mountPoint,
		FSType:      fsType,
		Options:     options,
		Total:       total,
		Used:        used,
		Avail:       avail,
		Free:        free,
		UsagePerc:   usagePerc,
		IsReadOnly:  isRO,
		IsSpecial:   isSpecial,
		Inodes:      inodesTotal,
		InodesUsed:  inodesUsed,
		InodesFree:  inodesFree,
	}, nil
}

// Unescape octal escape sequences like \040 for spaces in mount entries
func unescapeMount(s string) string {
	res := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) && isOctal(s[i+1]) && isOctal(s[i+2]) && isOctal(s[i+3]) {
			val := (s[i+1]-'0')*64 + (s[i+2]-'0')*8 + (s[i+3] - '0')
			res = append(res, val)
			i += 3
		} else {
			res = append(res, s[i])
		}
	}
	return string(res)
}

func isOctal(b byte) bool {
	return b >= '0' && b <= '7'
}

func isSpecialFSType(fs string) bool {
	switch strings.ToLower(fs) {
	case "sysfs", "proc", "devtmpfs", "devpts", "tmpfs", "cgroup", "cgroup2",
		"pstore", "bpf", "securityfs", "selinuxfs", "fusectl", "mqueue",
		"hugetlbfs", "tracefs", "debugfs", "configfs", "autofs", "nsfs",
		"overlay", "efivarfs", "binfmt_misc", "squashfs", "ramfs":
		return true
	default:
		return false
	}
}
