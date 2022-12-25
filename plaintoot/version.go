package plaintoot

import "fmt"

// ldflags will be set by goreleaser
var version = "vDEV"
var commit = "NONE"
var date = "UNKNOWN"

func VersionString() string {
	return fmt.Sprintf("plaintoot %s (%s), built on %s", version, commit, date)
}

func VersionStringShort() string {
	return fmt.Sprintf("plaintoot-%s; commit=%s", version, commit)
}
