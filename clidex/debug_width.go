//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	res, _ := http.Get("https://gitlab.com/phoneybadger/pokemon-colorscripts/-/raw/main/colorscripts/small/regular/pikachu")
	body, _ := io.ReadAll(res.Body)
	sprite := string(body)
	re := regexp.MustCompile(`\x1b\[[\d;]*[A-Za-z]`)
	lines := strings.Split(strings.TrimRight(sprite, "\n"), "\n")
	for i, line := range lines {
		stripped := re.ReplaceAllString(line, "")
		// what the fk is going on ;(
		fmt.Printf("line %d: rune_len=%d byte_len=%d stripped=[%s]\n", i, len([]rune(stripped)), len(stripped), stripped)
	}
}
