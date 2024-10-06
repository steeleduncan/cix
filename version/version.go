package version

// set during build
var BuildRevision string

func Version() string {
	if BuildRevision == "" {
		return "(dev)"
	} else {
		return BuildRevision
	}
}
