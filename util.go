package main

import (
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

const lsbPath string = "/etc/os-release"
const memFile string = "/proc/meminfo"

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
	currentDesktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if currentDesktop != "" {
		return currentDesktop
	}
	propsID, _ := exec.Command("xprop", "-root", "-notype", "_NET_SUPPORTING_WM_CHECK").Output()
	//tell me if regexp is wrong
	winid := regexpInByteArr(propsID, *regexp.MustCompile(`0x.*`), 0)
	wmNameRegexp := regexp.MustCompile(`.*?_NET_WM_NAME = "(.*)"`)
	wmProps, _ := exec.Command("xprop", "-id", winid, "-notype").CombinedOutput()
	WM_NAME := wmNameRegexp.FindStringSubmatch(string(wmProps))[1]
	// commands for getting name from xprop got from neofetch!! everyone say thanks dylan araps!!!
	return string(WM_NAME)
}

func memory() (int, int) {
	contents, _ := os.ReadFile(memFile)

	memTotal := regexpInByteArr(contents, *regexp.MustCompile(`.*?MemTotal:( *)(\d*)`), 2)
	memFree := regexpInByteArr(contents, *regexp.MustCompile(`.*?MemAvailable:( *)(\d*)`), 2)

	intTotalMatch, _ := strconv.Atoi(memTotal)
	intFreeMatch, _ := strconv.Atoi(memFree)

	nonAvailable := intTotalMatch - intFreeMatch

	return intTotalMatch / 1024, nonAvailable / 1024
}

func kernelVersion() string {
	file := "/proc/version"
	content, _ := os.ReadFile(file)
	version := regexpInByteArr(content, *regexp.MustCompile(`version (.*?) `), 1)
	return version
}

func distroName() (string, string) {
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
