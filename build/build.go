package build

import "runtime/debug"

// Version is dynamically set during build by the toolchain or overridden in the Makefile.
var Version = "dev"

// Date is dynamically set during build time in the Makefile.
var Date = "unknown"

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
