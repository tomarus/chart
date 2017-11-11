package data

// Collection defines an array of datasets.
type Collection []Data

// Normalize normalizes all values.
func (c Collection) Normalize(limit int) {
	for n := range c {
		c[n].normalize(limit)
	}

	max := c.max()
	for i := range c {
		c[i].normalizeMax(limit, max)
	}
}

// max returns the max value of a Collection.
func (c Collection) max() float64 {
	max := 0.
	for _, cl := range c {
		if max < cl.Max {
			max = cl.Max
		}
	}
	return max
}

// Implement Sort interface

func (c Collection) Len() int {
	return len(c)
}

func (c Collection) Less(i, j int) bool {
	return c[i].Max > c[j].Max
}

func (c Collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
