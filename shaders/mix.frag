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

uniform int SHOW_HUD;

void main() {
    // consts
    int DRAW_MODE_ADD = 1;
    int DRAW_MODE_MERGE = 2;
    int DRAW_MODE_SMELL = 3;
    
    // get colors
    vec2 zoom_transl = vec2(0.5 * ( 1 - ZOOM), 0.5 * ( 1 - ZOOM));
    vec2 transl = vec2(X_TRANSLATE, Y_TRANSLATE);
    vec2 FrameCoords = (TexCoord  * ZOOM)  + transl + zoom_transl;

    vec4 color_smell_red = texture(PfSmellRedTexture, FrameCoords);
    vec4 color_smell_green = texture(PfSmellGreenTexture, FrameCoords);
    vec4 color_game = texture(PfGameTexture, FrameCoords);

    // Just add together
    if (MODE == DRAW_MODE_ADD){
        FragColor = 0.5 * (color_smell_red + color_smell_green) + color_game;
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

    // draw cross hairs
    if (SHOW_HUD == 1) {
        float cross_size_px = 10.0;
        float line_thickness_px = 1.0;
        float cross_width = cross_size_px / window_width;
        float cross_height = cross_size_px / window_height;
        float line_width = line_thickness_px / window_width;
        float line_height = line_thickness_px / window_height;

        vec4 cross_color = vec4(.3,.3,.3,0.0);

            float x_d = abs(0.5 - TexCoord.x);
        float y_d = abs(0.5 - TexCoord.y);

        if (x_d < line_width && y_d < cross_height) {
            FragColor += cross_color;
        }
        if (y_d < line_height && x_d < cross_width) {
            FragColor += cross_color;
        }
    }
}