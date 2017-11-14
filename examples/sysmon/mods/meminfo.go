package mods

import (
	"github.com/c9s/goprocinfo/linux"
)

// MemInfo defines global cpu usage. Idle, system and user cpu time are supported.
type MemInfo struct {
	free store
}

// Len returns the amount of datapoints of the data available.
func (c *MemInfo) Len() int {
	return len(c.free.values)
}

// Data returns a slice of all Datasets available.
func (c *MemInfo) Data() []Dataset {
	return []Dataset{
		{"MEM Free", c.free.values},
	}
}

// Update retreives new values from /proc/stat
func (c *MemInfo) Update() error {
	stat, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return err
	}
	c.free.set(float64(stat.MemFree * 1024))
	return nil
}
