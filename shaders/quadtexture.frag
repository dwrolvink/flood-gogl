#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D texture1;

uniform float x;
uniform float tex_x;
uniform float tex_y;
uniform float tex_divisions;
uniform float tex_fliph;

void main() {
    FragColor = texture(texture1, TexCoord);
}