#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

uniform sampler2D PfGameTexture;
uniform sampler2D PfSmellRedTexture;
uniform int MODE;
uniform float A;

void main() {
    // consts
    int DRAW_MODE_ADD = 1;
    int DRAW_MODE_MERGE = 2;
    int DRAW_MODE_SMELL = 3;
    

    // get colors
    vec4 color_smell_red = texture(PfSmellRedTexture, TexCoord);
    vec4 color_game = texture(PfGameTexture, TexCoord);

    // Mix the two textures by ratio
    if (MODE == DRAW_MODE_MERGE){
        FragColor = A*color_smell_red + (1.0-A)*color_game;
    }
    // Just add together
    else if (MODE == DRAW_MODE_ADD){
        FragColor = color_smell_red + color_game;
    }   
 
    // Blue screen of bugs
    else {
        FragColor = vec4(0.,0.,1.0,1.0);
    }


    
}