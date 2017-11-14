package mods

import (
	"github.com/c9s/goprocinfo/linux"
)

// DiskStat defines global cpu usage. Idle, system and user cpu time are supported.
type DiskStat struct {
	rio, wio store
	lastrio  int64
	lastwio  int64
	first    bool
}

// Len returns the amount of datapoints of the data available.
func (c *DiskStat) Len() int {
	return len(c.rio.values)
}

// Data returns a slice of all Datasets available.
func (c *DiskStat) Data() []Dataset {
	return []Dataset{
		{"Read Bytes", c.rio.values},
		{"Write Bytes", c.wio.values},
	}
}

// Update retreives new values from /proc/stat
func (c *DiskStat) Update() error {
	stat, err := linux.ReadDiskStats("/proc/diskstats")
	if err != nil {
		return err
	}
	rio, wio := int64(0), int64(0)
	for _, n := range stat {
		rio += n.GetReadBytes()
		wio += n.GetWriteBytes()
	}
	if !c.first {
		c.first = true
	} else {
		c.rio.set(float64(rio - c.lastrio))
		c.wio.set(float64(wio - c.lastwio))
	}
	c.lastrio = rio
	c.lastwio = wio
	return nil
}
