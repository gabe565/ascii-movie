package config

import _ "embed"

//go:embed movies/sw1.txt
var Movie []byte

const FrameHeight = 14
const Width = 67

const PadTop = 3
const PadLeft = 6
const PadBottom = 3

const OutputDir = "generated_frames"
