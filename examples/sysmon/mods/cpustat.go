package mods

import (
	"math"
	"time"

	"github.com/c9s/goprocinfo/linux"
)

// CPUStat defines global cpu usage. Idle, system and user cpu time are supported.
type CPUStat struct {
	idle      store
	user      store
	sys       store
	lastIdle  float64
	lastUser  float64
	lastSys   float64
	lasttime  time.Time
	firsttime bool
}

// Len returns the amount of datapoints of the data available.
func (c *CPUStat) Len() int {
	return len(c.idle.values)
}

// Data returns a slice of all Datasets available.
func (c *CPUStat) Data() []Dataset {
	return []Dataset{
		{"CPU Idle", c.idle.values},
		{"CPU User", c.user.values},
		{"CPU System", c.sys.values},
	}
}

// Update retreives new values from /proc/stat
func (c *CPUStat) Update() error {
	stat, err := linux.ReadStat("/proc/stat")
	if err != nil {
		return err
	}
	now := time.Now()

	tdelta := float64(now.Sub(c.lasttime)) / 1e9
	numcpus := float64(len(stat.CPUStats))

	idle := float64(stat.CPUStatAll.Idle)
	fidle := ((idle - c.lastIdle) / numcpus) / tdelta
	c.lastIdle = idle
	user := float64(stat.CPUStatAll.User + stat.CPUStatAll.Nice)
	fuser := ((user - c.lastUser) / numcpus) / tdelta
	c.lastUser = user
	sys := float64(stat.CPUStatAll.System)
	fsys := ((sys - c.lastSys) / numcpus) / tdelta
	c.lastSys = sys

	if !c.firsttime {
		// skip first val, we calculate diffs between 2 points
		c.firsttime = true
	} else {
		c.idle.set(math.Min(fidle, 100.))
		c.user.set(math.Min(fuser, 100.))
		c.sys.set(math.Min(fsys, 100.))
	}
	c.lasttime = now
	return nil
}
