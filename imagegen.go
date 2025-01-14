package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"os"

	"path/filepath"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

type Mod struct {
	DisplayName        string
	RankTotal          int
	DownloadsTotal     int
	DownloadsYesterday int
}

var myClient = &http.Client{Timeout: 10 * time.Second}

var mods []Mod

func generateImage(steamId string) {
	getJson("https://tmlapis.thelonelysheep.repl.co/author_api/"+steamId, &mods)
	fmt.Printf("mods : %+v", mods)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

var imageWidth float64
var imageHeight float64

const margin float64 = 20.0
const padding float64 = 5.0

func run() error {
	imageWidth = 878.0
	imageHeight = (35.0+padding)*float64(len(mods)) + (35 * 2) + margin*2 + 10
	dc := gg.NewContext(int(imageWidth), int(imageHeight))

	// Draw light gray rounded rectangle
	dc.SetColor(color.RGBA{35, 39, 42, 255})
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, 15)
	dc.Fill()

	// Load font
	fontPath := filepath.Join("fonts", "Andy Bold.ttf")
	if err := dc.LoadFontFace(fontPath, 35); err != nil {
		return err
	}

	// Draw dark gray rectangle and leave 20px border
	x := margin
	y := margin
	w := float64(imageWidth) - (2.0 * margin)
	h := float64(imageHeight) - (2.0 * margin)
	dc.SetColor(color.RGBA{25, 28, 30, 255})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	// Draw Text
	DrawText(dc, "NotLe0n's Stats", imageWidth/3, margin*2+10)
	for i := 0; i < len(mods); i++ {
		_, nameTextHeight := dc.MeasureString(mods[i].DisplayName)
		dowloadsTextWidth, _ := dc.MeasureString(strconv.Itoa(mods[i].DownloadsTotal))

		DrawText(dc, strconv.Itoa(mods[i].RankTotal), 30, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+margin*2+10)
		DrawText(dc, mods[i].DisplayName, 120, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+margin*2+10)
		DrawText(dc, strconv.Itoa(mods[i].DownloadsTotal), imageWidth-dowloadsTextWidth-50, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+margin*2+10)
	}

	// save image as output.png
	if err := dc.SavePNG("output.png"); err != nil {
		return err
	}

	return nil
}

func DrawText(dc *gg.Context, s string, x float64, y float64) {
	dc.SetColor(color.White)
	textWidth, textHeight := dc.MeasureString(s)
	x = ClampFloat(x, 0, imageWidth-textWidth)
	y = ClampFloat(y, textHeight, imageHeight-textHeight)
	dc.DrawString(s, x, y)
}

func ClampFloat(v float64, min float64, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
