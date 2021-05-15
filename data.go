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
func SetData() []gogl.DataObject {

	// List of datasets
	datalist := make([]gogl.DataObject, 3)

	// Fist dataset: Vertex type: Quad, uses Sprites
	// -----------------------------------------------------------
	datalist[0] = gogl.DataObject{
		ProgramName: "Game",
		Type:        gogl.GOGL_QUADS,
		Vertices:    CreateQuadVertexMatrix(1.0, 0.0, 0.0),
		Indices: []uint32{
			1, 0, 3, // triangle 1
			0, 2, 3, // triangle 2
		},
		VertexShaderSource:   "shaders/quad.vert",
		FragmentShaderSource: "shaders/game.frag",
	}

	datalist[1] = gogl.DataObject{
		ProgramName: "Smell",
		Type:        gogl.GOGL_QUADS,
		Vertices:    CreateQuadVertexMatrix(1.0, 0.0, 0.0),
		Indices: []uint32{
			1, 0, 3, // triangle 1
			0, 2, 3, // triangle 2
		},
		VertexShaderSource:   "shaders/quad.vert",
		FragmentShaderSource: "shaders/smell.frag",
	}

	datalist[2] = gogl.DataObject{
		ProgramName: "BlitGameState",
		Type:        gogl.GOGL_QUADS,
		Vertices:    CreateQuadVertexMatrix(1.0, 0.0, 0.0),
		Indices: []uint32{
			1, 0, 3, // triangle 1
			0, 2, 3, // triangle 2
		},
		VertexShaderSource:   "shaders/quad.vert",
		FragmentShaderSource: "shaders/blit_game.frag",
	}

	for i := 0; i < len(datalist); i++ {
		datalist[i].ProcessData()
	}

	return datalist
}
