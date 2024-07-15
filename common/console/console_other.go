//go:build !windows
// +build !windows

package console

import "os"

func SetTitle(title string) {
	os.Stdout.WriteString("\033]0;" + title + "\007")
}
