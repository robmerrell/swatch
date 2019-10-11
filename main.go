package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type swatchEntry struct {
	Hue        float64 `json:"hue"`
	Brightness float64 `json:"brightness"`
	Saturation float64 `json:"saturation"`
	Alpha      int     `json:"alpha"`
	ColorSpace int     `json:"colorSpace"`
}

type swatch struct {
	Name     string         `json:"name"`
	Swatches []*swatchEntry `json:"swatches"`
}

func main() {
	// make sure the user has given enough arguments
	if len(os.Args) < 4 {
		usage()
		os.Exit(1)
	}

	// read the arguments given
	filename := os.Args[1]
	swatchName := os.Args[2]
	colors := os.Args[3]

	fmt.Printf("Generating '%s' as %s\n", swatchName, filename)

	splitColors := strings.Split(colors, " ")
	palette := []*swatch{&swatch{
		Name:     swatchName,
		Swatches: make([]*swatchEntry, 30),
	}}

	// parse the colors
	for i, color := range splitColors {
		// if the color doesn't begin with a #, stick it on there
		if color[0] != '#' {
			color = "#" + color
		}

		// parse the color and save the output
		parsedColor, err := colorful.Hex(color)
		if err != nil {
			log.Fatal(err)
		}

		h, s, v := parsedColor.Hsv()
		palette[0].Swatches[i] = &swatchEntry{Hue: h / 360.0, Saturation: s, Brightness: v, Alpha: 1, ColorSpace: 0}
	}

	// convert the palette to json
	jsonOutput, err := json.Marshal(palette)
	if err != nil {
		log.Fatal(err)
	}

	// write the palette file
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	file, err := writer.Create("Swatches.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(jsonOutput)
	if err != nil {
		log.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("usage: swatch filename swatch-name colors")
	fmt.Println(`ex: swatch ocean.swatches "Ocean Nights" "#000000 #123123 #ff00ff"`)
}
