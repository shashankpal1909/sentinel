package logger

import "fmt"

const banner = `
  _____ _____ _   _ _____ _____ _   _ _____ _     
 / ____|  ____| \ | |_   _|_   _| \ | |  ____| |    
| (___ | |__  |  \| | | |   | | |  \| | |__  | |    
 \___ \|  __| | . ` + "`" + ` | | |   | | | . ` + "`" + ` |  __| | |    
 ____) | |____| |\  | | |  _| |_| |\  | |____| |____
|_____/|______|_| \_| |_| |_____|_| \_|______|______|

`

// PrintBanner prints the ASCII art banner for Sentinel.
func PrintBanner() {
	fmt.Print(banner)
}
