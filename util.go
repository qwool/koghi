package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/pbnjay/memory"
)

const lsbPath string = "/etc/os-release"
const memFile string = "/proc/meminfo"

// const hostnameFile string = "/etc/hostname"
const pacmanFile string = "/var/lib/pacman/local"
const kernelVerFile string = "/proc/version"

// pre-compiled regexps for speed
var (
	wmNameRegexp     = regexp.MustCompile(`.*?_NET_WM_NAME = "(.*)"`)
	windowIDRegexp   = regexp.MustCompile(`0x.*`)
	kernelRegexp     = regexp.MustCompile(`version (.*?) `)
	lsbNameRegexp    = regexp.MustCompile(`PRETTY_NAME="([^"]+)"`)
	lsbVersionRegexp = regexp.MustCompile(`(?m)VERSION_ID="?([^"\n]+)"?`)
)

func catch(err error) {
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func colASCII(arr []string, line int) string {
	if line >= len(arr) {
		return ""
	}
	return fmt.Sprintf("\x1b[1m\x1b[90m%s\x1b[0m", arr[line])
}

/*
this cut code by like 10loc instantly and made me able to read it
also i had less names to assign so i kinda follow style
TODO: do it on every function
*/
func regexpInByteArr(array []byte, regexp *regexp.Regexp, index int) string {
	matches := regexp.FindStringSubmatch(string(array))
	if len(matches) <= index {
		return ""
	}
	return matches[index]
}

func wm() string {
	if runtime.GOOS == "darwin" {
		return "WindowServer" //TODO: support for custom mac WMs
	}
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
	winid := regexpInByteArr(propsID, windowIDRegexp, 0)
	wmProps, err := exec.Command("xprop", "-id", winid, "-notype").CombinedOutput()
	if err != nil {
		return "no wm"
	}
	// commands for getting name from xprop got from neofetch!! everyone say thanks dylan araps!!!
	return regexpInByteArr(wmProps, wmNameRegexp, 1)
}

func get_memory() (string, string) {
	total := memory.TotalMemory()
	used := total - memory.FreeMemory()
	return strconv.Itoa(int(used / 1024 / 1024)), strconv.Itoa(int(total / 1024 / 1024))
}

func kernelVersion() string {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("uname", "-r")
		output, err := cmd.Output()
		if err != nil {
			return "darwin"
		}
		return strings.TrimSpace("darwin " + string(output))
	}
	content, err := os.ReadFile(kernelVerFile)
	if err != nil {
		return ""
	}
	return regexpInByteArr(content, kernelRegexp, 1)
}

func archCountPkgs() (string, error) {
	// yes, pacman -Q | wc -l exists but it takes 50ms to run
	// and this seems considerably faster
	files, err := os.ReadDir(pacmanFile)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(files)), nil
}

func distroName() (string, string) {
	switch runtime.GOOS {
	case "darwin":
		nameCmd := exec.Command("sw_vers", "-productName")
		name, err := nameCmd.Output()
		if err != nil {
			return "macOS", ""
		}

		versionCmd := exec.Command("sw_vers", "-productVersion")
		version, err := versionCmd.Output()
		if err != nil {
			return "macOS", ""
		}

		return strings.TrimSpace(string(name)), strings.TrimSpace(string(version))

	case "linux":
		// check if arch for packages
		if conf.countPkg {
			_, err := os.Stat(pacmanFile)
			if err == nil {
				pkgs, err := archCountPkgs()
				if err != nil {
					return "arch", ""
				}
				return pkgs, ""
			}
		}

		lsbFile, err := os.ReadFile(lsbPath)
		if err != nil {
			cmd := exec.Command("lsb_release", "-ds")
			output, err := cmd.Output()
			if err != nil {
				return "Linux", ""
			}
			return strings.TrimSpace(string(output)), ""
		}

		name := regexpInByteArr(lsbFile, lsbNameRegexp, 1)
		version := regexpInByteArr(lsbFile, lsbVersionRegexp, 1)

		if name == "" {
			name = "Linux"
		}

		return name, version

	default:
		// use uname for any unsupported oses
		osCmd := exec.Command("uname", "-s")
		os, _ := osCmd.Output()

		verCmd := exec.Command("uname", "-r")
		ver, _ := verCmd.Output()

		if len(os) == 0 {
			return "idk but", "listen to earl more"
		}

		return strings.TrimSpace(string(os)), strings.TrimSpace(string(ver))
	}
}

func hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "puter", err
	}
	return strings.TrimSpace(hostname), nil
}

func getUser() (string, error) {
	curUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return curUser.Username, nil
}
