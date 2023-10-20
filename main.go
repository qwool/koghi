package main

import (
	"fmt"
	"strings"
)

// config starts here

var conf = Config{
	countPkg:     true,    // slows down by 10ms
	presetWMName: "ababa", // might be useful for wayland, leave blank for no
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

	//define lines
	lines := []string{
		fmt.Sprintf("\x1b[1m\x1b[97m%s@%s\x1b[0m", username, hostname),
		fmt.Sprintf("\x1b[1m\x1b[36m%s\x1b[0m \x1b[1m\x1b[97m%s %s\x1b[0m", "os", distro_name, distro_version),
		fmt.Sprintf("\x1b[1m\x1b[35m%s\x1b[0m \x1b[1m\x1b[97m%dM/%dM\x1b[0m", "mem", memFree, mem),
		fmt.Sprintf("\x1b[1m\x1b[31m%s\x1b[0m \x1b[1m\x1b[97m%s\x1b[0m", "wm", wm),
		fmt.Sprintf("\x1b[1m\x1b[32m%s\x1b[0m \x1b[1m\x1b[97m%s\x1b[0m", "kernel", kernel),
	}
	for i, v := range lines {
		fmt.Printf("%s %s\n", colASCII(asciiArray, i), v)
	} // fmt.Printf("%s@%s\n%s %s\n%s of %s\n%s\n%s\n", username, hostname, distro_name, distro_version, fmt.Sprint(memFree), fmt.Sprint(mem), wm, kernel)
}
