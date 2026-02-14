package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type DexEntry struct {
	Name    string        `json:"name"`
	ID      int           `json:"id"`
	Types   []PokemonType `json:"types"`
	Sprites Sprites       `json:"sprites"`
}

type Sprites struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
}

// Go needs multiple structs for multi level json
type PokemonType struct {
	Type TypeName `json:"type"`
}

type TypeName struct {
	Name string `json:"name"`
}

type SpeciesData struct {
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
}

type FlavorTextEntry struct {
	FlavorText string   `json:"flavor_text"`
	Language   Language `json:"language"`
}

type Language struct {
	Name string `json:"name"`
}

func main() {
	shiny := flag.Bool("shiny", false, "Show shiny variant")
	flag.Parse()

	args := flag.Args() // non-flag arguments
	if len(args) < 1 {
		fmt.Println("Usage: dex <--shiny> <pokemon name>")
		return
	}

	name := args[0]
	fmt.Println("Searching for:", strings.ToLower(name)+"...")

	// ---------- POKEMON INFO FETCH ----------
	reqInfo := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	resInfo, errInfo := http.Get(reqInfo)

	if errInfo != nil {
		fmt.Println("Error fetching data:", errInfo)
		return
	}

	defer resInfo.Body.Close()
	if resInfo.StatusCode == 404 {
		fmt.Println("Could not find that pokemon...")
		return
	}

	bodyInfo, errInfo := io.ReadAll(resInfo.Body)

	// ---------- SPECIES INFO FETCH ----------
	reqSpeciesInfo := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%s", name)
	resSpeciesInfo, errSpeciesInfo := http.Get(reqSpeciesInfo)

	if errSpeciesInfo != nil {
		fmt.Println("Error fetching data:", errSpeciesInfo)
		return
	}

	defer resSpeciesInfo.Body.Close()
	if resSpeciesInfo.StatusCode == 404 {
		fmt.Println("Could not find that pokemon...")
		return
	}

	bodySpeciesInfo, errSpeciesInfo := io.ReadAll(resSpeciesInfo.Body)

	var description SpeciesData
	json.Unmarshal(bodySpeciesInfo, &description)

	// ---------- SPRITE FETCH ----------
	var reqSprite string
	if *shiny {
		reqSprite = fmt.Sprintf("https://gitlab.com/phoneybadger/pokemon-colorscripts/-/raw/main/colorscripts/small/shiny/%s", name)
	} else {
		reqSprite = fmt.Sprintf("https://gitlab.com/phoneybadger/pokemon-colorscripts/-/raw/main/colorscripts/small/regular/%s", name)
	}

	resSprite, errSprite := http.Get(reqSprite)
	if errSprite != nil {
		fmt.Println("Couldn't fetch sprite:", errSprite)
	}

	defer resSprite.Body.Close()

	var entry DexEntry
	json.Unmarshal(bodyInfo, &entry)

	// get sprite as a string
	var sprite string
	if resSprite.StatusCode == 200 {
		bodySprite, _ := io.ReadAll(resSprite.Body)
		sprite = string(bodySprite)
	} else if !*shiny && entry.Sprites.FrontDefault != "" {
		sprite = renderSprite(entry.Sprites.FrontDefault)
	} else if *shiny && entry.Sprites.FrontShiny != "" {
		sprite = renderSprite(entry.Sprites.FrontShiny)
	} else {
		sprite = "No sprite available"
	}

	// build info lines
	var infoLines []string
	infoLines = append(infoLines, fmt.Sprintf("Name: %s", strings.ToUpper(entry.Name)))
	infoLines = append(infoLines, fmt.Sprintf("ID: %d", entry.ID))

	typeStr := "Type: "
	for i := 0; i < len(entry.Types); i++ {
		typeStr += strings.ToUpper(entry.Types[i].Type.Name)
		if i != len(entry.Types)-1 {
			typeStr += ", "
		}
	}
	infoLines = append(infoLines, typeStr)

	for i := 0; i < len(description.FlavorTextEntries); i++ {
		if description.FlavorTextEntries[i].Language.Name == "en" {
			infoLines = append(infoLines, "")
			// clean up newlines/form feeds from API text and word-wrap
			desc := strings.ReplaceAll(description.FlavorTextEntries[i].FlavorText, "\n", " ")
			desc = strings.ReplaceAll(desc, "\f", " ")
			// wrap so the box doesnt break
			wrapped := wordWrap(desc, 35)
			for j, line := range wrapped {
				if j == 0 {
					infoLines = append(infoLines, "Desc: "+line)
				} else {
					infoLines = append(infoLines, "      "+line)
				}
			}
			break
		}
	}

	printSideBySide(sprite, infoLines)

}

// stripAnsi removes ANSI escape codes so we can measure visible string length
func stripAnsi(s string) string {
	// wtf is this regex ;(
	re := regexp.MustCompile(`\x1b\[[\d;]*[A-Za-z]`)
	return re.ReplaceAllString(s, "")
}

// wordWrap splits a string into lines of at most maxWidth characters, breaking at spaces
func wordWrap(s string, maxWidth int) []string {
	words := strings.Fields(s)
	var lines []string
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

// displayWidth counts the visible column width of a string
func displayWidth(s string) int {
	return len([]rune(stripAnsi(s)))
}

// printSideBySide prints the sprite on the left and info in a box on the right
func printSideBySide(sprite string, infoLines []string) {
	spriteLines := strings.Split(strings.TrimRight(sprite, "\n"), "\n")

	// find the widest sprite line (visible characters only)
	spriteMaxWidth := 0
	for _, line := range spriteLines {
		w := displayWidth(line)
		if w > spriteMaxWidth {
			spriteMaxWidth = w
		}
	}

	// find the widest info line to size the box
	boxContentWidth := 0
	for _, line := range infoLines {
		if len([]rune(line)) > boxContentWidth {
			boxContentWidth = len([]rune(line))
		}
	}
	boxContentWidth += 2 // padding inside box

	// build the box lines: top border, content rows, bottom border
	white := "\033[97m" // bright white
	reset := "\033[0m"
	var boxLines []string
	boxLines = append(boxLines, white+"╭"+strings.Repeat("─", boxContentWidth)+"╮"+reset)
	for _, line := range infoLines {
		padding := strings.Repeat(" ", boxContentWidth-len([]rune(line))-1)
		boxLines = append(boxLines, white+"│ "+line+padding+"│"+reset)
	}
	boxLines = append(boxLines, white+"╰"+strings.Repeat("─", boxContentWidth)+"╯"+reset)

	// print lines side by side
	totalLines := len(spriteLines)
	if len(boxLines) > totalLines {
		totalLines = len(boxLines)
	}

	gap := "    " // space between sprite and box

	for i := 0; i < totalLines; i++ {
		spritePart := ""
		if i < len(spriteLines) {
			spritePart = spriteLines[i]
		}

		// pad sprite side to consistent width
		visibleLen := displayWidth(spritePart)
		padding := strings.Repeat(" ", spriteMaxWidth-visibleLen)

		boxPart := ""
		if i < len(boxLines) {
			boxPart = boxLines[i]
		}

		fmt.Println(spritePart + padding + gap + boxPart)
	}
}

// renderSprite fetches a PNG from a URL and converts it to colored terminal art
// using ANSI truecolor escape codes and half-block characters (▀▄)
// each terminal character represents 2 vertical pixels
// CAUTION: this part was largely gen'd by claude. thanks claude
func renderSprite(url string) string {
	res, err := http.Get(url)
	if err != nil {
		return "Could not fetch sprite"
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return "Could not decode sprite"
	}

	bounds := img.Bounds()

	// find the bounding box of non-transparent pixels to trim empty space
	minX, minY, maxX, maxY := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a >= 128 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	var result strings.Builder

	// scale factor — 1 means full size, 2 means half, etc.
	scale := 1

	for y := minY; y <= maxY; y += 2 * scale {
		for x := minX; x <= maxX; x += scale {
			topR, topG, topB, topA := img.At(x, y).RGBA()
			topTransparent := topA < 128

			// bottom pixel may not exist if image height is odd
			var botR, botG, botB, botA uint32
			botTransparent := true
			if y+scale < bounds.Max.Y {
				botR, botG, botB, botA = img.At(x, y+scale).RGBA()
				botTransparent = botA < 128
			}

			if topTransparent && botTransparent {
				result.WriteString(" ")
			} else if topTransparent {
				// only bottom pixel visible — use ▄ with foreground color
				br, bg, bb := botR>>8, botG>>8, botB>>8
				result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm▄\033[0m", br, bg, bb))
			} else if botTransparent {
				// only top pixel visible — use ▀ with foreground color
				tr, tg, tb := topR>>8, topG>>8, topB>>8
				result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm▀\033[0m", tr, tg, tb))
			} else {
				// both pixels visible — ▀ foreground=top, background=bottom
				tr, tg, tb := topR>>8, topG>>8, topB>>8
				br, bg, bb := botR>>8, botG>>8, botB>>8
				result.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀\033[0m", tr, tg, tb, br, bg, bb))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}
