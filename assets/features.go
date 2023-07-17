package assets

import "embed"

//go:embed features.json
var Features []byte

func dummy() {
	embed.FS.ReadFile(embed.FS{}, "")
}
