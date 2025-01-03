package scanner

import "fmt"

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
	bold         = "\033[1m"
	underline    = "\033[4m"
)

func PrintInfo(format string, a ...interface{}) {
	fmt.Printf(colorBlue+"[*] "+format+colorReset+"\n", a...)
}

func PrintSuccess(format string, a ...interface{}) {
	fmt.Printf(colorGreen+"[+] "+format+colorReset+"\n", a...)
}

func PrintError(format string, a ...interface{}) {
	fmt.Printf(colorRed+"[-] "+format+colorReset+"\n", a...)
}

func PrintVulnFound(format string, a ...interface{}) {
	fmt.Printf("\n"+bold+"\033[41m\033[37m[!] "+format+colorReset+"\n", a...)
}

func PrintVulnDetail(name, value string) {
	fmt.Printf(colorCyan+"    %-15s"+colorReset+" : "+colorYellow+"%s"+colorReset+"\n", name, value)
}

// 导出常量供其他包使用
const (
	ColorReset  = colorReset
	Bold        = bold
	ColorYellow = colorYellow
	ColorCyan   = colorCyan
)
