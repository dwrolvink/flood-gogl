#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

uniform sampler2D PfTexture;

void main() {
    FragColor = texture(PfTexture, TexCoord);
}