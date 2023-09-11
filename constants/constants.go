package constants

import (
	"lib/golog"
)

var Log *golog.TypLog

const LogLevel = golog.LogLvlDebug

const WorkingDir = "./"

const (
	IP   = "127.0.0.1"
	PORT = 9090
)

const (
	CSV_HEADER = "user,date,palette,color,r,g,b,evaluation"
	COLOR_1    = "primary"
	COLOR_2    = "text"
	COLOR_3    = "background"
	COLOR_4    = "button_primary"
	COLOR_5    = "button_secondary"
)
