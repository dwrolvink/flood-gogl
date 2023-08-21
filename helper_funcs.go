package main

import (
	"image"
	"os"
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func abs[T Number](a T) T {
	if a < 0.0 {
		return -a
	}
	return a
}

func max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func makeTimestamp() int64 {
	return time.Now().UnixMilli()
}

// Easy way to create a quad with a certain size and offset
func CreateQuadVertexMatrix(size float32, x_offset float32, y_offset float32) []float32 {
	screen_left := -size + x_offset
	screen_bottom := -size + y_offset
	screen_right := size + x_offset
	screen_top := size + y_offset
	texture_top := float32(1.0)
	texture_bottom := float32(0.0)
	texture_left := float32(0.0)
	texture_right := float32(1.0)
	//z := float32(0.0)

	vertices := []float32{
		// x, y, z, texcoordx, texcoordy
		screen_left, screen_top, texture_left, texture_top,
		screen_right, screen_top, texture_right, texture_top,
		screen_left, screen_bottom, texture_left, texture_bottom,
		screen_right, screen_bottom, texture_right, texture_bottom,
	}

	return vertices
}

// Tests if int element is in slice
func contains_si(slice []int, element int) bool {
	for _, val := range slice {
		if val == element {
			return true
		}
	}
	return false
}

type img_dims struct {
	Width  int
	Height int
}

func get_image_dimensions(image_path string) img_dims {
	reader, err := os.Open(StartImageSrc)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	bounds := m.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	return img_dims{
		Width:  w,
		Height: h,
	}
}
