package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"net/http"
	"os"
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

	if len(os.Args) < 2 {
		fmt.Println("Usage: dex <pokemon name>")
		return
	}

	name := os.Args[1]
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
	reqSprite := fmt.Sprintf("https://gitlab.com/phoneybadger/pokemon-colorscripts/-/raw/main/colorscripts/small/regular/%s", name)
	resSprite, errSprite := http.Get(reqSprite)
	if errSprite != nil {
		fmt.Println("Couldn't fetch sprite:", errSprite)
	}

	defer resSprite.Body.Close()

	var entry DexEntry
	json.Unmarshal(bodyInfo, &entry)

	if resSprite.StatusCode == 200 {
		bodySprite, errSprite := io.ReadAll(resSprite.Body)
		if errSprite != nil {
			fmt.Println("Error reading sprite:", errSprite)
		} else {
			fmt.Println(string(bodySprite))
		}
	} else {
		// Fallback: fetch PNG sprite from PokeAPI and render it
		if entry.Sprites.FrontDefault != "" {
			fmt.Print(renderSprite(entry.Sprites.FrontDefault))
		} else {
			fmt.Println("No sprite available")
		}
	}
	fmt.Println("Name:", strings.ToUpper(entry.Name))
	fmt.Println("ID:", entry.ID)
	fmt.Print("Type: ")
	for i := 0; i < len(entry.Types); i++ {
		if i != len(entry.Types)-1 {
			fmt.Print(strings.ToUpper(entry.Types[i].Type.Name))
			fmt.Print(", ")
		} else {
			fmt.Print(strings.ToUpper(entry.Types[i].Type.Name) + "\n\n")
		}
	}

	for i := 0; i < len(description.FlavorTextEntries); i++ {
		if description.FlavorTextEntries[i].Language.Name == "en" {
			fmt.Println("Desc:", description.FlavorTextEntries[i].FlavorText)
			break
		}
	}

}

// credit to claude for sprite render backup
// renderSprite fetches a PNG from a URL and converts it to colored terminal art
// using ANSI truecolor escape codes and half-block characters (▀▄)
// each terminal character represents 2 vertical pixels
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
