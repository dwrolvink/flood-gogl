package main

import (
	"errors"
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

// Used in PrtSc.
// Reads out the pixel data in gl.FRONT, and saves it to recording/temp/image<Tick>.png when mode == RECORDING_GIF
// When mode == RECORDING_PRTSC, the save location is "recording/printscreens/<datetime>.png
func CreateImage(number int, mode int) {
	// Init
	width := Width
	height := Height

	// Set folder and filename
	var filename string
	var folder string
	if mode == RECORDING_PRTSC {
		currentTime := time.Now()
		filename = currentTime.Format("2006-01-02 15:04:05.000000") + ".png"
		folder = "recording/printscreens/"
	} else {
		filename = fmt.Sprintf("image%03d.png", number)
		folder = "recording/temp/"
	}

	// Create folder
	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Create empty image to receive the pixel data
	img := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	// Copy the pixel data from the default front buffer to the image
	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, int32(Width), int32(Height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	img = imaging.FlipV(img)

	// Make bacground black (in case starting texture has an alpha layer)
	byteIndex := 3
	for byteIndex < len(img.Pix) {
		img.Pix[byteIndex] = byte(255)
		byteIndex += 4
	}

	// Encode as PNG.
	f, _ := os.Create(folder + filename)
	png.Encode(f, img)

	// Terminal feedback
	if mode == RECORDING_GIF {
		rl := "?"
		if record_length > 0 {
			rl = fmt.Sprint(record_length)
		}
		fmt.Println("Recording frame", tick-record_start-1, "of", rl)
	} else if mode == RECORDING_PRTSC {
		fmt.Println("Created " + folder + filename)
	}

}

// Takes all the frame images in recording/temp and makes an mp4/palletted-gif out of it using ffmpeg.
func CompileGif() {
	// terminal feedback
	fmt.Println("Compiling gif, don't close the window.")

	// run script to compile output and cleanup separate images
	filename := time.Now().Unix()
	cmd, err := exec.Command("/bin/sh", "scripts/make_gif.sh", fmt.Sprint(filename), fmt.Sprint(fps), fmt.Sprint(record_start)).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}

	fmt.Println(cmd)                 // empty when no problems
	fmt.Println("Compilation done.") //
}
