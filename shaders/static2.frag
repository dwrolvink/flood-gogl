#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

uniform sampler2D PfGameTexture;
uniform sampler2D PfSmellTexture;



uniform float window_width;
uniform float window_height;

void main() {
    vec2 resolution = vec2(window_width, window_height);

    vec4 Tex1Color = texture(PfGameTexture, TexCoord);
    vec4 Tex2Color = texture(PfSmellTexture, TexCoord);

    if (TexCoord.x > 0.5){
        FragColor = Tex1Color;
    }
    else {
        FragColor =  Tex2Color;
    }
    
}