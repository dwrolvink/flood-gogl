package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.5-core/gl"
)

const (
	RECORDING_GIF   = 0
	RECORDING_PRTSC = 1
)

/*
	These functions allow us to read from the front buffer and save it as a file.
	Continually making a picture every frame, and compiling all the frames into a gif is also supported.
*/

// Ensures that the recording folder (recording/temp) is present.
func InitRecording() {
	err := os.Mkdir("recording/temp/", 0755)
	if err != nil {
		//panic(err)
	}
}

// Reads out the pixel data in gl.FRONT, and saves it to recording/temp/image<Tick>.png when mode == RECORDING_GIF
// When mode == RECORDING_PRTSC, the save location is "recording/printscreens/<datetime>.png
func CreateImage(number int, mode int) {
	if mode == RECORDING_GIF {
		fmt.Println("Recording frame", tick, "of", record_length)
	}

	filename := fmt.Sprintf("image%03d.png", number)
	width := Width
	height := Height

	var folder string
	if mode == RECORDING_PRTSC {
		folder = "recording/printscreens/"

		currentTime := time.Now()
		filename = currentTime.Format("2006-01-02 15:04:05.000000") + ".png"
	} else {
		folder = "recording/temp/"
	}

	img := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	img = imaging.FlipV(img)

	// Make bacground black (starting texture has an alpha layer)
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	byteIndex := 0
	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			//pixels[byteIndex] = byte(r / 256)
			byteIndex++
			//pixels[byteIndex] = byte(g / 256)
			byteIndex++
			//pixels[byteIndex] = byte(b / 256)
			byteIndex++
			img.Pix[byteIndex] = byte(255)
			byteIndex++
		}
	}

	// Encode as PNG.
	f, _ := os.Create(folder + filename)
	png.Encode(f, img)

	if mode == RECORDING_PRTSC {
		fmt.Println("Created " + folder + filename)
	}

}

// Takes all the frame images in recording/temp and makes a palletted gif out of it using ffmpeg.
func CompileGif() {
	filename := time.Now().Unix()

	fmt.Println("Compiling gif, don't close the window.")

	cmd, err := exec.Command("/bin/sh", "scripts/make_gif.sh", fmt.Sprint(filename)).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fmt.Println(cmd)

	fmt.Println("Compilation done.")
}
