#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

// input
uniform sampler2D PfGameTexture; 
uniform sampler2D PfSmellRedTexture; 
uniform sampler2D PfSmellGreenTexture; 

uniform int iTime;
uniform float window_width;
uniform float window_height;

uniform vec2 Actor1;
uniform float Actor1Radius;

vec2 Resolution = vec2(window_width, window_height);

// settings
float MinGrowth = 1.01;

// reusables
vec2 Top = (vec2(0, 1)/ Resolution);
vec2 Bottom = (vec2(0, -1)/ Resolution);
vec2 Left = (vec2(-1, 0)/ Resolution);
vec2 Right = (vec2(1, 0)/ Resolution);


vec2 uv_to_coords(vec2 uv) {
    vec2 coords = (TexCoord * Resolution);
    coords.x = floor(coords.x);
    coords.y = floor(coords.y);
    
    return coords;
}

vec4 get_color_by_coord(sampler2D tex, vec2 coords) {
// vec2 uv = (vec2(pixelX, pixelY) + .5) / resolutionOfTexture;
    vec2 uv = (coords + .5) / Resolution;
    return texture2D(tex, uv);
}

void main() {
    vec4 GameColor = texture2D(PfGameTexture, TexCoord);
    // GameColor.r = 1.0;
    // GameColor.b = 0.0;
    // GameColor.g = 0.0;
    // GameColor.a = 1.0;

    vec2 coords = uv_to_coords(TexCoord);
    // if (coords.x > 20.0 && coords.y == 400.0) {
    //     GameColor = vec4(1.0,1.0,0.0,1.0);
    // }

    // else {
        GameColor = texture2D(PfGameTexture, ((TexCoord * Resolution) + vec2(0.000,1.00000)) / Resolution);
        // GameColor = get_color_by_coord(PfGameTexture, vec2(coords.x, coords.y + 1./Resolution ));
    // }

    FragColor = GameColor;
    
}

