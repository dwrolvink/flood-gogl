package main

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
