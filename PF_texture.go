package main

import (
	"github.com/dwrolvink/gogl"
	"github.com/go-gl/gl/v4.5-core/gl"

	"image"
)

/*
	The code in this file allows us to read pixel code from the front buffer, and save it to a texture (or an image).
	This texture can then be read by the shader, to allow for an easy way to input and output data to and from the GPU.
*/

// Wrapper function to create a texture with default settings for this project
func NewDefaultTexture() gogl.TextureID {
	textureID := gogl.GenTexture()
	gl.BindTexture(gl.TEXTURE_2D, uint32(textureID))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return textureID
}

// Load given image into the given texture. Used to reset the game state.
func ResetFrame(textureID gogl.TextureID, imgSource string) {
	// load image data from source image
	pixels, _ := gogl.LoadPixelDataFromImage(imgSource)

	// Put the pixels in the texture
	gogl.BindTexture(textureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(Width), int32(Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(*pixels))
}

// Gets pixels from front buffer, and puts them in the pixel array of the given image.
func ReadBuffer(bufferID uint32, imgStorage *image.NRGBA) {
	gl.ReadBuffer(bufferID)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgStorage.Pix))
}

// Bind FBO, and copy its pixel data over to the given texture
func TakeBufferSnapshot(bufferID uint32, textureID gogl.TextureID) {
	// We are going to overwrite the read buffer, get the current one
	// so that we can set it back afterwards
	var drawFboId int32
	gl.GetIntegerv(gl.DRAW_FRAMEBUFFER_BINDING, &drawFboId)

	// Bind given FBO to read buffer so that we can get the pixels
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, bufferID)

	// Copy pixel data over to the given texture
	gl.BindTexture(gl.TEXTURE_2D, uint32(textureID)) //A texture you have already created storage for
	gl.CopyTexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, 0, 0, Width, Height)

	// Restore defaults before function call
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, uint32(drawFboId))

}

// Inits framebuffer @ given id, with the given texture as the first and only color attachment.
// - Make sure that the given texture is already initialized! (ResetFrame(...)).
func CreateFramebuffer(frameBufferId *uint32, textureId gogl.TextureID) {
	// generate a framebuffer
	gl.GenFramebuffers(1, frameBufferId)

	// attach texture to fbo
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, *frameBufferId)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, uint32(textureId), 0)

	// check state of framebuffer
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("framebuffer not complete!")
	}

	// bind the default framebuffer for normal operation
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
