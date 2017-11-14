package mods

import (
	"github.com/c9s/goprocinfo/linux"
)

// LoadAvg defines global cpu usage. Idle, system and user cpu time are supported.
type LoadAvg struct {
	m1, m5, m15 store
}

// Len returns the amount of datapoints of the data available.
func (c *LoadAvg) Len() int {
	return len(c.m1.values)
}

// Data returns a slice of all Datasets available.
func (c *LoadAvg) Data() []Dataset {
	return []Dataset{
		{"1 Minute", c.m1.values},
		{"5 Minute", c.m5.values},
		{"15 Minute", c.m15.values},
	}
}

// Update retreives new values from /proc/stat
func (c *LoadAvg) Update() error {
	stat, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		return err
	}
	c.m1.set(stat.Last1Min)
	c.m5.set(stat.Last5Min)
	c.m15.set(stat.Last15Min)
	return nil
}
