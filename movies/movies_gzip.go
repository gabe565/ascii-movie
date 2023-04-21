//go:build gzip

package movies

import "embed"

const Default = "sw1.txt.gz"

//go:embed *.txt.gz
var Movies embed.FS
