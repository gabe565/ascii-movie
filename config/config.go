package config

import (
	_ "embed"
	"path/filepath"
)

//go:embed movies/sw1.txt
var DefaultMovie []byte

const FrameHeight = 14
const Width = 67

const PadTop = 3
const PadLeft = 6
const PadBottom = 3

var OutputDir = filepath.Join("internal", "movie")
