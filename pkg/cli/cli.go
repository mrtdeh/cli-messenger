package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

var (
	filePath string
	Room     string
	// room     string
)
var (
	l *readline.Instance
)

func SetReadLine(a *readline.Instance) {
	l = a
}

func SetRoome(room string) {
	Room = room
}

func Close() error {
	return l.Close()
}

func ListernerReadline(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {

	switch int(key) {
	// 1. If press TAB
	case readline.CharTab:
		// 2. If file Path is not empty
		if filePath != "" {
			// 3. Complete file path :
			c := filePathCompleter(filePath)
			// 3.1 Return one value
			if len(c) == 1 { // 3.1 Then
				b := c[0]
				filePath = string(b)
				// Return completed path
				return stringTorune(b), len(b), true
			}

			// 3.2 If path not exist
			// Add @ to first if not contain it
			if len(filePath) > 0 && filePath[:1] != "@" {
				filePath = "@" + filePath
			}
			// Return unchanged path
			return stringTorune(filePath), len(filePath), true
		}

		// 4. Else, Return line without TAB character
		return line[:len(line)-1], len(line) - 1, true

	}

	// If line is not empty
	if len(line) > 0 {
		// Retrive sign
		sign := string(line[0])
		// Check sign is @
		if sign == "@" && len(string(line)) > 1 {
			// Fill file path variable
			filePath = string(line)[1:]
		} else {
			filePath = ""
		}
	}
	return line, pos, true
}

// func filterInput(r rune) (rune, bool) {
// 	switch r {
// 	case readline.CharCtrlZ:
// 		return r, false
// 	}
// 	return r, true
// }

func filePathCompleter(line string) (c []string) {
	var pathPart, dir string
	EndWithSlash := false
	if line[len(line)-1:] == "/" {
		EndWithSlash = true
	}

	if line[:1] == "@" {
		pathPart = line[1:]
	} else {
		pathPart = line
	}

	if pathPart == "/" {
		dir = pathPart
	} else {
		dir = filepath.Dir(strings.TrimRight(pathPart, "/"))
	}

	base := filepath.Base(pathPart)
	prefix := filepath.Join(dir, base)

	if EndWithSlash {
		prefix += "/"
	}

	matches, err := filepath.Glob(prefix + "*")
	if err != nil {
		return
	}

	if len(matches) > 1 {
		fmt.Println()
		for _, m := range matches {
			fmt.Printf(path.Base(m) + " ")
		}
		fmt.Println()
		c = append(c, "@"+commonPrefix(matches...))
		return
	}
	for _, match := range matches {
		if EndWithSlash {
			match = match + "/"
		}

		c = append(c, "@"+match)

		if f, _ := os.Stat(c[0][1:]); f != nil {
			if !f.IsDir() {
				c[0] = c[0] + " "
			} else {
				if c[0][len(c[0])-1:] != "/" {
					c[0] = c[0] + "/"
				}
			}
		}
	}

	return
}
