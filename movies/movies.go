package movies

import (
	"embed"
)

const Default = "sw1.txt"

//go:embed sw1.txt rick_roll.txt
var Movies embed.FS
