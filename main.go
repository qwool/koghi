package main

import (
	"fmt"
	"strings"
)

// config starts here

type Theme struct {
	bold   string
	white  string
	cyan   string
	purple string
	red    string
	green  string
}

type Config struct {
	countPkg     bool   // counts packages on arch, adds few ms
	presetWMName string // for wayland, blank = disabled
	theme        Theme  // color theme
	format       Format // format strings for each line
}

type Format struct {
	userHost string // format for username@hostname
	distro   string // format for os/distro line
	memory   string // format for memory info
	wm       string // format for window manager
	kernel   string // format for kernel version
}

var defaultTheme = Theme{
	bold:   "\x1b[1m",
	white:  "\x1b[0m \x1b[1m\x1b[97m",
	cyan:   "\x1b[36m",
	purple: "\x1b[35m",
	red:    "\x1b[31m",
	green:  "\x1b[32m",
}

var defaultFormat = Format{
	userHost: "             %s@%s",
	distro:   " %s%s%s%s%s %s",
	memory:   "%s%s%s%s%sM/%sM",
	wm:       "%s%s%s%s%s",
	kernel:   " %s%s%s%s%s",
}

var conf = Config{
	countPkg:     true,
	presetWMName: "",
	theme:        defaultTheme,
	format:       defaultFormat,
}

// ascii art alignment done with spaces
var ascii = `
  |\'/-..--.
 / _ _   ,  ;
'~='Y'~_<._./
 <'-....__.'
 `

func main() {
	asciiLines := strings.Split(ascii, "\n")

	username, err := getUser()
	if err != nil {
		fmt.Printf("failed getting username: %v\n", err)
		return
	}

	hostname, err := hostname()
	if err != nil {
		fmt.Printf("failed getting hostname: %v\n", err)
		return
	}

	mem, memFree := get_memory()
	distroName, distroVersion := distroName()
	wmName := wm()
	kernel := kernelVersion()

	// arch vs other distros
	osLine := "os"
	if distroVersion == "" {
		osLine = "pkgs"
	}

	t := conf.theme
	f := conf.format

	lines := []string{
		fmt.Sprintf("             %s@%s", username, hostname),
	    fmt.Sprintf(f.distro, t.bold, t.cyan, osLine, t.white, distroName, distroVersion),
	    fmt.Sprintf(f.memory, t.bold, t.purple, "mem", t.white, memFree, mem),
	    fmt.Sprintf(f.wm, t.bold, t.red, "wm", t.white, wmName),
	    fmt.Sprintf(f.kernel, t.bold, t.green, "kernel", t.white, kernel),
	}

	for i, line := range lines {
		if i < len(asciiLines) {
			fmt.Printf("%s %s\n", colASCII(asciiLines, i), line)
		}
	}
}
