package format

// Fit trims a string to the specified length adding ... at the end when required.
func Fit(s string, n int) string {
	if len(s) > n {
		return s[0:n] + "..."
	}
	return s
}
