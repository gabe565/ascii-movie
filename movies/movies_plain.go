//go:build !gzip

package movies

import "embed"

const Default = "sw1.txt"

//go:embed *.txt
var Movies embed.FS
