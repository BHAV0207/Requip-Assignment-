// Package version holds build metadata injected at build time with -ldflags.
// It lets every running pod report exactly which commit it was built from.
package version

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)
