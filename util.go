package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

const lsbPath string = "/etc/os-release"
const memFile string = "/proc/meminfo"

type Config struct {
	countPkg     bool
	presetWMName string
}

func catch(err error) {
	if err != nil {
		fmt.Println("err: %v\n", err)
	}
}

func colASCII(arr []string, line int) string {
	return fmt.Sprintf("\x1b[1m\x1b[90m" + arr[line] + "\x1b[0m")
}

/*
this cut code by like 10loc instantly and made me able to read it
also i had less names to assign so i kinda follow style
TODO: do it on every function
*/
func regexpInByteArr(array []byte, regexp regexp.Regexp, index int) string {
	matches := regexp.FindStringSubmatch(string(array))
	return matches[index]
}

func wm() string {
	if conf.presetWMName != "" {
		return conf.presetWMName
	}
	currentDesktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if currentDesktop != "" {
		return currentDesktop
	} // if you have env variable XDG_CURRENT_DESKTOP it shaves off like 10ms sometimes
	propsID, err := exec.Command("xprop", "-root", "-notype", "_NET_SUPPORTING_WM_CHECK").Output()
	if err != nil {
		return ""
	}
	//tell me if regexp is wrong
	winid := regexpInByteArr(propsID, *regexp.MustCompile(`0x.*`), 0)
	wmNameRegexp := regexp.MustCompile(`.*?_NET_WM_NAME = "(.*)"`)
	wmProps, err := exec.Command("xprop", "-id", winid, "-notype").CombinedOutput()
	if err != nil {
		return ""
	}
	WM_NAME := wmNameRegexp.FindStringSubmatch(string(wmProps))[1]
	// commands for getting name from xprop got from neofetch!! everyone say thanks dylan araps!!!
	return string(WM_NAME)
}

func memory() (string, string) {
	contents, _ := os.ReadFile(memFile)

	memTotal := regexpInByteArr(contents, *regexp.MustCompile(`.*?MemTotal:( *)(\d*)`), 2)
	memFree := regexpInByteArr(contents, *regexp.MustCompile(`.*?MemAvailable:( *)(\d*)`), 2)

	intTotalMatch, _ := strconv.Atoi(memTotal)
	intFreeMatch, _ := strconv.Atoi(memFree)

	nonAvailable := intTotalMatch - intFreeMatch

	return strconv.Itoa(intTotalMatch / 1024), strconv.Itoa(nonAvailable / 1024)
}

func kernelVersion() string {
	file := "/proc/version"
	content, _ := os.ReadFile(file)
	version := regexpInByteArr(content, *regexp.MustCompile(`version (.*?) `), 1)
	return version
}

func archCountPkgs() (string, error) {
	// yes, pacman -Q | wc -l exists but it takes 50ms to run
	// and this seems considerably faster
	files, err := os.ReadDir("/var/lib/pacman/local")
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(files)), nil
}

func distroName() (string, string) {
	if conf.countPkg == true {
		_, err := os.Stat("/var/lib/pacman/local")
		if err == nil {
			pkgs, err := archCountPkgs()
			if err != nil {
				return "arch", ""
			}
			return "arch, pkgs:", pkgs
		}
	}
	lsbFile, _ := os.ReadFile(lsbPath)
	lsbContent := string(lsbFile)

	regexpName := regexp.MustCompile(`PRETTY_NAME="([^"]+)"`)
	regexpVersion := regexp.MustCompile(`(?m)BUILD_ID=([^"]+$)`)

	nameMatch := regexpName.FindStringSubmatch(lsbContent)
	versionMatch := regexpVersion.FindStringSubmatch(lsbContent)
	if nameMatch != nil {
		return nameMatch[1], versionMatch[1]
	}
	return "undefined OS", "undefined ver"
}

func hostname() (string, error) {
	hostnameFile := "/etc/hostname"
	hostname, err := os.ReadFile(hostnameFile)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(hostname), "\n", ""), nil
}

func getUser() (string, error) {
	curUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return curUser.Username, err
}
