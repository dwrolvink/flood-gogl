#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D PfGameTexture;
uniform sampler2D PfSmellTexture;

uniform float window_width;
uniform float window_height;


void main() {
    FragColor = texture(PfGameTexture, TexCoord);
}