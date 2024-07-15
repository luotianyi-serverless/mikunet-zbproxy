package color

import (
	"strconv"
)

// Foreground text colors
const (
	FgBlack = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

const (
	escapeHead    = "\x1b["
	escapeTail    = "m"
	compiledReset = "\x1b[0m"
)

var List = [...]int64{
	FgRed, FgGreen, FgYellow,
	FgBlue, FgMagenta, FgCyan, FgWhite,

	FgHiRed, FgHiGreen, FgHiYellow,
	FgHiBlue, FgHiMagenta, FgHiCyan, FgHiWhite,
}

func Apply(color int64, str string) string {
	return escapeHead + strconv.FormatInt(color, 10) + escapeTail + str + compiledReset
}
