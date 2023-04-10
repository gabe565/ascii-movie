package movies

import (
	"embed"
)

const Default = "sw1.txt"

//go:embed sw1.txt
var Movies embed.FS
