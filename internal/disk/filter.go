package disk

import (
	"sort"
	"strings"

	"df-go/internal/config"
)

// FilterOptions specifies rules to filter and sort mount points
type FilterOptions struct {
	ShowAll       bool     // Show all mounts (including pseudo and zero size)
	ExplicitPaths []string // Specific paths requested via CLI arguments or -m
	SortBy        string   // Column to sort by: mount, size, used, avail, perc, device
	ReverseSort   bool     // Sort descending
}

// FilterMounts filters and sorts the list of mounts according to config and options
func FilterMounts(mounts []MountInfo, cfg *config.Config, opts FilterOptions) []MountInfo {
	// If explicit CLI paths are specified, prioritize them
	filterList := opts.ExplicitPaths
	if len(filterList) == 0 && len(cfg.DefaultMounts) > 0 {
		filterList = cfg.DefaultMounts
	}

	excludeMountSet := make(map[string]bool)
	for _, m := range cfg.ExcludeMounts {
		excludeMountSet[m] = true
	}

	excludeFSTypeSet := make(map[string]bool)
	for _, t := range cfg.ExcludeFSTypes {
		excludeFSTypeSet[strings.ToLower(t)] = true
	}

	var filtered []MountInfo

	for _, m := range mounts {
		// 1. If explicit paths or default_mounts are set, check if this mount matches
		if len(filterList) > 0 {
			matched := false
			for _, target := range filterList {
				target = strings.TrimSuffix(target, "/")
				mountClean := strings.TrimSuffix(m.MountPoint, "/")
				if target == "" {
					target = "/"
				}
				if mountClean == "" {
					mountClean = "/"
				}

				if mountClean == target || m.Device == target {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// 2. If not in ShowAll mode, check exclude lists and real disk criteria
		if !opts.ShowAll {
			// Check excluded mounts
			if excludeMountSet[m.MountPoint] {
				continue
			}

			// Check excluded prefix
			isExcludedPrefix := false
			for _, ex := range cfg.ExcludeMounts {
				if strings.HasPrefix(m.MountPoint, ex+"/") {
					isExcludedPrefix = true
					break
				}
			}
			if isExcludedPrefix {
				continue
			}

			// Check excluded filesystem types
			if excludeFSTypeSet[strings.ToLower(m.FSType)] {
				continue
			}

			// If OnlyRealDisks is true, filter out non-physical devices and zero-size mounts
			if cfg.OnlyRealDisks {
				if m.Total == 0 {
					continue
				}
				if m.IsSpecial {
					continue
				}
				if !strings.HasPrefix(m.Device, "/dev/") &&
					!strings.HasPrefix(m.FSType, "fuse") &&
					!strings.Contains(m.Options, "local") {
					continue
				}
			}
		}

		filtered = append(filtered, m)
	}

	// Sorting
	sortMounts(filtered, opts.SortBy, opts.ReverseSort)

	return filtered
}

func sortMounts(mounts []MountInfo, sortBy string, reverse bool) {
	if sortBy == "" {
		sortBy = "mount"
	}

	sort.SliceStable(mounts, func(i, j int) bool {
		var less bool
		switch strings.ToLower(sortBy) {
		case "device", "fs", "filesystem":
			less = mounts[i].Device < mounts[j].Device
		case "size", "total":
			less = mounts[i].Total < mounts[j].Total
		case "used":
			less = mounts[i].Used < mounts[j].Used
		case "avail", "free":
			less = mounts[i].Avail < mounts[j].Avail
		case "perc", "use%", "usage":
			less = mounts[i].UsagePerc < mounts[j].UsagePerc
		case "type", "fstype":
			less = mounts[i].FSType < mounts[j].FSType
		case "mount", "on", "mounted":
			fallthrough
		default:
			// Root "/" always first, then alphabetical
			if mounts[i].MountPoint == "/" {
				less = true
			} else if mounts[j].MountPoint == "/" {
				less = false
			} else {
				less = mounts[i].MountPoint < mounts[j].MountPoint
			}
		}

		if reverse {
			return !less
		}
		return less
	})
}
