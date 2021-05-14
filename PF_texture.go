package main

import (
	"github.com/dwrolvink/gogl"
	"github.com/go-gl/gl/v4.5-core/gl"

	"image"
)

/*
	The code in this file allows us to read pixel code from the front buffer, and save it to a texture.
	This texture can then be read by the shader, to allow for an easy way to input and output data to and from the GPU.
*/

// Resets the FP_texture to initial conditions
func ResetFrame(img *image.NRGBA) {
	pixels, _ := gogl.LoadPixelDataFromImage("assets/img/start.png")

	for i := 0; i < len(img.Pix) && i < len(*pixels); i++ {
		img.Pix[i] = (*pixels)[i]
	}
}

// Gets pixels from front buffer, and puts them in `previous_frame.Pix`
func ReadFrontBuffer() {
	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(previous_frame.Pix))
}

// Fill PF_texture with new pixel data.
// ActionReset == true will cause this function to load initial data in, instead of reading pixels from the front buffer.
func UpdateFrontBufferTexture() {

	// Either reset image, or load image data from frontbuffer
	if ActionReset || tick == 1 {
		ResetFrame(previous_frame)
		ActionReset = false
	} else {
		ReadFrontBuffer()
	}

	// Bind texture
	gogl.BindTexture(PF_textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Write image data to the PF_texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(Width), int32(Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(previous_frame.Pix))
}
