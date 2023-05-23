package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

type VersionInfo struct {
	Version string
	Commit  string
}

// Version determines version and commit information based on multiple data sources:
//   - Version information dynamically added by `git archive` in the remaining to parameters.
//   - A hardcoded version number passed as first parameter.
//   - Commit information added to the binary by `go build`.
//
// It's supposed to be called like this in combination with setting the `export-subst` attribute for the corresponding
// file in .gitattributes:
//
//	var Version = version.Version("1.0.0-rc2", "$Format:%(describe)$", "$Format:%H$")
//
// When exported using `git archive`, the placeholders are replaced in the file and this version information is
// preferred. Otherwise the hardcoded version is used and augmented with commit information from the build metadata.
func Version(version, gitDescribe, gitHash string) *VersionInfo {
	if !strings.HasPrefix(gitDescribe, "$") && !strings.HasPrefix(gitHash, "$") {
		return &VersionInfo{
			Version: gitDescribe,
			Commit:  gitHash,
		}
	} else {
		commit := ""

		if info, ok := debug.ReadBuildInfo(); ok {
			modified := false

			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					commit = setting.Value
				case "vcs.modified":
					modified, _ = strconv.ParseBool(setting.Value)
				}
			}

			const hashLen = 7 // Same truncation length for the commit hash as used by git describe.

			if len(commit) >= hashLen {
				version += "-g" + commit[:hashLen]

				if modified {
					version += "-dirty"
					commit += " (modified)"
				}
			}
		}

		return &VersionInfo{
			Version: version,
			Commit:  commit,
		}
	}
}

// Print writes verbose version output to stdout.
func (v *VersionInfo) Info() string {
	ret := fmt.Sprintf(`
### vadoVIA version: 

[%v](https://github.com/gterranova/vadovia) <br />

<br />

### Build information:

Go version: %s (%s, %s)

`, v.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	if v.Commit != "" {
		ret += "Git commit: " + v.Commit + "\n\n"
	}

	ret += `

---

### This software is licensed under the terms of the BSD 3-Clause License to: 

**Sani Zangrando Avvocati**

---

**DISCLAIMER:** The information provided by this application does not, and is not intended to, constitute legal advice; instead, they are meant for general informational purposes only.  Information provided by this application  may not constitute the most up-to-date legal or other information. 
Users of this application should contact their attorney to obtain advice with respect to any particular legal matter. No user of this application should act or refrain from acting on the basis of information provided by the application without first seeking legal advice from counsel in the relevant jurisdiction. Only your individual attorney can provide assurances that the information contained herein – and your interpretation of it – is applicable or appropriate to your particular situation. Use of, and access to, this application do not create an attorney-client relationship between the user, author, contributing law firm, its respective professional. 

---

*Copyright (c) 2023 [Gianpaolo Terranova](mailto:g.terranova@sazalex.com)* <br /><br />

`

	return ret
}

// Print writes verbose version output to stdout.
func (v *VersionInfo) Print() {
	fmt.Println("vadoVIA version:", v.Version)
	fmt.Println()

	fmt.Println("Build information:")
	fmt.Printf("  Go version: %s (%s, %s)\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if v.Commit != "" {
		fmt.Println("  Git commit:", v.Commit)
	}
}
