package main

import (
	"fmt"
	"strings"
)

// config starts here

var conf = Config{
	countPkg:     true, // slows down by 10ms, only for Arch
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
	// catch(err)

	cBold := "\x1b[1m"
	cWhite := "\x1b[0m \x1b[1m\x1b[97m"
	cCyan := "\x1b[36m"
	cPorpur := "\x1b[35m"
	cRed := "\x1b[31m"
	cGreen := "\x1b[32m"

	//define lines
	lines := []string{
		cBold + username + hostname,
		cBold + cCyan + "os" + cWhite + distro_name + distro_version,
		cBold + cPorpur + "mem" + cWhite + memFree + mem,
		cBold + cRed + "wm" + cWhite + wm,
		cBold + cGreen + "kernel" + cWhite + kernel,
	}
	for i, v := range lines {
		fmt.Printf("%s %s\n", colASCII(asciiArray, i), v)
	} // fmt.Printf("%s@%s\n%s %s\n%s of %s\n%s\n%s\n", username, hostname, distro_name, distro_version, fmt.Sprint(memFree), fmt.Sprint(mem), wm, kernel)
}
