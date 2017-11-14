package mods

import (
	"time"

	"github.com/c9s/goprocinfo/linux"
)

// NetDev defines global cpu usage. Idle, system and user cpu time are supported.
type NetDev struct {
	rx, tx   store
	lastrx   uint64
	lasttx   uint64
	first    bool
	lasttime time.Time
}

// Len returns the amount of datapoints of the data available.
func (c *NetDev) Len() int {
	return len(c.rx.values)
}

// Data returns a slice of all Datasets available.
func (c *NetDev) Data() []Dataset {
	return []Dataset{
		{"RX Bytes", c.rx.values},
		{"TX Bytes", c.tx.values},
	}
}

// Update retreives new values from /proc/stat
func (c *NetDev) Update() error {
	stat, err := linux.ReadNetworkStat("/proc/net/dev")
	if err != nil {
		return err
	}
	now := time.Now()
	tdelta := float64(now.Sub(c.lasttime)) / 1e9
	c.lasttime = now

	tx, rx := uint64(0), uint64(0)
	for _, n := range stat {
		tx += n.TxBytes
		rx += n.RxBytes
	}
	if !c.first {
		c.first = true
	} else {
		c.tx.set(float64(tx-c.lasttx) / tdelta)
		c.rx.set(float64(rx-c.lastrx) / tdelta)
	}
	c.lastrx = rx
	c.lasttx = tx
	return nil
}
