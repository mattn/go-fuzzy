package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-fuzzy"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	out := colorable.NewColorableStdout()
	pattern := os.Args[1]
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var matched []int
		s := scanner.Text()
		rs := []rune(s)
		li := 0
		line := ""
		if m, s := fuzzy.Match(pattern, s, &matched); m {
			for _, i := range matched {
				if li < len(rs) {
					line += string(rs[li:i]) + "\x1b[31m" + string(rs[i]) + "\x1b[0m"
				}
				li = i + 1
			}
			line += string(rs[li:])
			fmt.Fprintln(out, line)
		}
	}
}
