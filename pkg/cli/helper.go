package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	Header    = "\033[95m"
	Blue      = "\033[94m"
	Cyan      = "\033[96m"
	Green     = "\033[32m"
	Yellow    = "\033[93m"
	Red       = "\033[91m"
	End       = "\033[0m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
)

func GetColorCode(c string) color.Attribute {
	switch strings.ToLower(c) {
	case "red":
		return color.FgRed
	case "green":
		return color.FgGreen
	case "cyan":
		return color.FgCyan
	case "magenta":
		return color.FgMagenta
	case "white":
		return color.FgWhite
	default:
		return color.FgCyan
	}
}
func stringTorune(str string) (converter []rune) {

	for _, runeValue := range str {
		converter = append(converter, rune(runeValue))
	}
	return
}
func clearLine() {
	fmt.Print("\033[2K\r") // ANSI escape sequence to clear the line
}

func PrintUserMessage(username, c, content string) {
	cc := GetColorCode(c)
	col := color.New(cc).SprintFunc()

	clearLine()
	fmt.Printf("%s@%s %s\n", Room, col(username), content)
	l.Refresh()
}

func Print(content string) {
	clearLine()
	fmt.Println(content)
	l.Refresh()
}

func commonPrefix(strs ...string) string {
	// بررسی طول آرایه ورودی
	if len(strs) == 0 {
		return ""
	} else if len(strs) == 1 {
		return strs[0]
	}

	// مشترک‌ترین قسمت را بین دو رشته پیدا کرده و در متغیر prefix ذخیره می‌کنیم
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}
