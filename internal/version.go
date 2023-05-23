package internal

import (
	"terra9.it/vadovia/version"
)

// Version contains version and Git commit information.
//
// The placeholders are replaced on `git archive` using the `export-subst` attribute.
var Version = version.Version("1.0.0-rc2", "$Format:%(describe)$", "$Format:%H$")
