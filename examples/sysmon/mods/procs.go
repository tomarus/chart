package mods

import (
	"github.com/c9s/goprocinfo/linux"
)

// Procs defines global cpu usage. Idle, system and user cpu time are supported.
type Procs struct {
	tot, run store
}

// Len returns the amount of datapoints of the data available.
func (c *Procs) Len() int {
	return len(c.tot.values)
}

// Data returns a slice of all Datasets available.
func (c *Procs) Data() []Dataset {
	return []Dataset{
		{"Total Procs", c.tot.values},
		{"Running Procs", c.run.values},
	}
}

// Update retreives new values from /proc/stat
func (c *Procs) Update() error {
	stat, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		return err
	}
	c.tot.set(float64(stat.ProcessTotal))
	c.run.set(float64(stat.ProcessRunning))
	return nil
}
