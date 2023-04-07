package config

import "path/filepath"

var MovieDir = "config/movies"
var MovieFile = "sw1.txt"

const FrameHeight = 14
const Width = 67

const PadTop = 3
const PadLeft = 6
const PadBottom = 3

var OutputDir = filepath.Join("internal", "movie")
