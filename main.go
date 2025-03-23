package main

import (
	"fmt"
	"strings"
)

type Theme struct {
	bold   string
	white  string
	cyan   string
	purple string
	red    string
	green  string
}

type Config struct {
	countPkg     bool   // counts packages on arch, adds few ms. does nothing if not on arch
	presetWMName string // for wayland, blank = disabled
	theme        Theme  // color theme
}

type Line struct {
	label    string
	format   string
	dataFunc func() []any
}

var defaultTheme = Theme{
	bold:   "\x1b[1m",
	white:  "\x1b[0m \x1b[1m\x1b[97m",
	cyan:   "\x1b[36m",
	purple: "\x1b[35m",
	red:    "\x1b[31m",
	green:  "\x1b[32m",
}

var conf = Config{
	countPkg:     true,
	presetWMName: "",
	theme:        defaultTheme,
}

var ascii = `
  |\'/-..--.
 / _ _   ,  ;
'~='Y'~_<._./
 <'-....__.'
`

const asciiWidth = 15 // padding from the start

func main() {
	asciiLines := strings.Split(ascii, "\n")

	lines := []Line{
		// line config starts here!!
		// label: obvious
		// format: each %s stands for a value in return []any{...data}, you can add extra characters for formatting
		// dataFunc: getters for values
		// return []any{...data}:
		{
			label:  "userHost",
			format: "%s@%s",
			dataFunc: func() []any {
				username, _ := getUser()
				hostname, _ := hostname()
				return []any{username, hostname}
			},
		},
		{
			label:  "os",
			format: "%s%s%s%s%s %s",
			dataFunc: func() []any {
				osName, osVersion := osName()
				osLine := "os"
				if osVersion == "" {
					osLine = "pkgs"
				}
				t := conf.theme
				return []any{t.bold, t.cyan, osLine, t.white, osName, osVersion}
			},
		},
		{
			label:  "memory",
			format: "%s%s%s%s%sM/%sM",
			dataFunc: func() []any {
				mem, memFree := getMemory()
				t := conf.theme
				return []any{t.bold, t.purple, "mem", t.white, memFree, mem}
			},
		},
		{
			label:  "wm",
			format: "%s%s%s%s%s",
			dataFunc: func() []any {
				wmName := wm()
				t := conf.theme
				return []any{t.bold, t.red, "wm", t.white, wmName}
			},
		},
		{
			label:  "kernel",
			format: "%s%s%s%s%s",
			dataFunc: func() []any {
				kernel := kernelVersion()
				t := conf.theme
				return []any{t.bold, t.green, "kernel", t.white, kernel}
			},
		},
	}

	for i, line := range lines {
		data := line.dataFunc()
		formattedLine := fmt.Sprintf(line.format, data...)

		padding := asciiWidth
		if i < len(asciiLines) {
			padding -= len(strings.TrimRight(asciiLines[i], " "))
		}

		fmt.Printf("%s%s%s\n", colASCII(asciiLines, i), strings.Repeat(" ", padding), formattedLine)
	}
}
