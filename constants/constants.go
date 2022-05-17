package constants

import (
	"github.com/bamboutech/golog"
)

var Log *golog.TypLog

const PORT = 8080

const (
	CSV_HEADER = "user,date,color,value,evaluation"
	COLOR_1    = "primary"
	COLOR_2    = "text_primary"
	COLOR_3    = "text_secondary"
	COLOR_4    = "button_primary"
	COLOR_5    = "button_secondary"
)
