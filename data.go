package main

import (
	"github.com/dwrolvink/gogl"
)

/*
	Multiple datasets can be defined.
	Each set contains all that it needs to draw to the screen,
	think of: Program, VOA, VBO, EBO, Textures, Sprites, etc

	Below, each dataset is defined, and added to datalist.
*/

// Define the DataObjects that contain our Programs, Shaders, Sprites, etc
func SetData() [](*gogl.DataObject) {

	// List of datasets
	datalist := make([]*gogl.DataObject, 4)

	// Defaults
	var vertex_type = gogl.GOGL_QUADS
	var verts = CreateQuadVertexMatrix(1.0, 0.0, 0.0)
	var vShader = "shaders/quad.vert"
	var indices = []uint32{
		1, 0, 3, // triangle 1
		0, 2, 3, // triangle 2
	}

	// Game state
	datalist[0] = &gogl.DataObject{
		ProgramName:          "Game",
		Type:                 vertex_type,
		Vertices:             verts,
		Indices:              indices,
		VertexShaderSource:   vShader,
		FragmentShaderSource: "shaders/game.frag",
	}

	// Smell red (blurred version of gamestate)
	datalist[1] = &gogl.DataObject{
		ProgramName:          "Smell Red",
		Type:                 vertex_type,
		Vertices:             verts,
		Indices:              indices,
		VertexShaderSource:   vShader,
		FragmentShaderSource: "shaders/smell_red.frag",
	}

	// Smell green (blurred version of gamestate)
	datalist[2] = &gogl.DataObject{
		ProgramName:          "Smell Green",
		Type:                 vertex_type,
		Vertices:             verts,
		Indices:              indices,
		VertexShaderSource:   vShader,
		FragmentShaderSource: "shaders/smell_green.frag",
	}

	// Mix. Can be used to do advanced merging of textures
	datalist[3] = &gogl.DataObject{
		ProgramName:          "Mix",
		Type:                 vertex_type,
		Vertices:             verts,
		Indices:              indices,
		VertexShaderSource:   vShader,
		FragmentShaderSource: "shaders/mix.frag",
	}

	// Compile Shader programs
	for i := 0; i < len(datalist); i++ {
		datalist[i].ProcessData()
	}

	return datalist
}
