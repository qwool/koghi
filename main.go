package main

import (
	"fmt"
	"strings"
)

// config starts here

var conf = Config{
	countPkg:     true, // if not Arch does nothing, on arch counts packages, slows down by a few ms
	presetWMName: "",   // might be useful for wayland, leave blank for no
}

// this doesnt line up properly by default, pad with spaces
var ascii string = `             
  |\'/-..--. 
 / _ _   ,  ;
'~='Y'~_<._./
 <'-....__.' `

// config ends here

func main() {

	asciiArray := strings.Split(ascii, "\n")
	// for _, v := range asciiArray {
	// 	if v == "" {
	// 		asciiArray = append(asciiArray)
	// 	}
	// }
	hostname, _ := hostname()
	// catch(err)
	username, _ := getUser()
	mem, memFree := memory()
	distro_name, distro_version := distroName()
	wm := wm()
	kernel := kernelVersion()

	cBold := "\x1b[1m"
	cWhite := "\x1b[0m \x1b[1m\x1b[97m"
	cCyan := "\x1b[36m"
	cPorpur := "\x1b[35m"
	cRed := "\x1b[31m"
	cGreen := "\x1b[32m"

	// yeah this was supposed to be a ternary operator
	// osLine := (map[bool]string{conf.countPkg || : "os", !conf.countPkg: "pkgs"})["os" > "pkgs"]

	var osLine string = "os"
	if distro_version == "" {
		osLine = "pkgs"
	}
	// define lines
	lines := []string{
		cBold + username + "@" + hostname,
		cBold + cCyan + osLine + cWhite + distro_name + " " + distro_version,
		cBold + cPorpur + "mem" + cWhite + memFree + "M/" + mem + "M",
		cBold + cRed + "wm" + cWhite + wm,
		cBold + cGreen + "kernel" + cWhite + kernel,
	}
	for i, v := range lines {
		fmt.Printf("%s %s\n", colASCII(asciiArray, i), v)
	} // fmt.Printf("%s@%s\n%s %s\n%s of %s\n%s\n%s\n", username, hostname, distro_name, distro_version, fmt.Sprint(memFree), fmt.Sprint(mem), wm, kernel)
}
