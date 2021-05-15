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
func ResetFrame(textureID gogl.TextureID, imgStorage *image.NRGBA, imgSource string, textureIndex uint32) {
	// load image data from source image
	pixels, _ := gogl.LoadPixelDataFromImage(imgSource)
	for i := 0; i < len(*pixels); i++ {
		imgStorage.Pix[i] = (*pixels)[i]
	}

	WritePixelsToTexture(textureID, imgStorage, textureIndex)

}

func WritePixelsToTexture(textureID gogl.TextureID, imgStorage *image.NRGBA, textureIndex uint32) {
	// Bind texture
	gl.ActiveTexture(gl.TEXTURE0 + textureIndex)
	gogl.BindTexture(textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Write image data to the PF_texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(Width), int32(Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgStorage.Pix))
}

// -----------------------------------
// Gets pixels from front buffer, and puts them in `previous_frame.Pix`
func ReadBuffer(bufferID uint32, imgStorage *image.NRGBA) {
	gl.ReadBuffer(bufferID)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgStorage.Pix))
}

func TakeBufferSnapshot(bufferID uint32, textureID gogl.TextureID, imgStorage *image.NRGBA, textureIndex uint32) {

	// Either reset image, or load image data from frontbuffer
	ReadBuffer(bufferID, imgStorage)

	WritePixelsToTexture(textureID, imgStorage, textureIndex)
}
