package logger

import "fmt"

const banner = `                 _   _            _ 
  ___  ___ _ __ | |_(_)_ __   ___| |
 / __|/ _ \ '_ \| __| | '_ \ / _ \ |
 \__ \  __/ | | | |_| | | | |  __/ |
 |___/\___|_| |_|\__|_|_| |_|\___|_|

`

// PrintBanner prints the ASCII art banner for Sentinel.
func PrintBanner() {
	fmt.Print(banner)
}
