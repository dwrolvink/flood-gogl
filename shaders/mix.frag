#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

uniform sampler2D PfGameTexture;
uniform sampler2D PfSmellRedTexture;
uniform sampler2D PfSmellGreenTexture;

uniform int MODE;
uniform float A;

uniform float window_width;
uniform float window_height;

uniform float ZOOM;
uniform float Y_TRANSLATE;
uniform float X_TRANSLATE;

void main() {
    // consts
    int DRAW_MODE_ADD = 1;
    int DRAW_MODE_MERGE = 2;
    int DRAW_MODE_SMELL = 3;
    
    // get colors
    vec2 transl = vec2(X_TRANSLATE, Y_TRANSLATE);

    vec4 color_smell_red = texture(PfSmellRedTexture, (TexCoord + transl)  * ZOOM );
    vec4 color_smell_green = texture(PfSmellGreenTexture, (TexCoord + transl)  * ZOOM);
    vec4 color_game = texture(PfGameTexture, (TexCoord + transl)  * ZOOM);

    // Just add together
    if (MODE == DRAW_MODE_ADD){
        FragColor = 0.8 * (color_smell_red + color_smell_green) + color_game;
    }
    // Mix smell and gamestate by ratio
    else if (MODE == DRAW_MODE_MERGE){
        FragColor = A*(color_smell_red + color_smell_green)  + (1.0-A)*color_game;
    }
    // Draw only smell (mix red & green)
    else if (MODE == DRAW_MODE_SMELL){
        FragColor = A*color_smell_red + (1.0-A)*color_smell_green;
    }        
    // Blue screen of bugs
    else {
        FragColor = vec4(0.,0.,1.0,1.0);
    }
}