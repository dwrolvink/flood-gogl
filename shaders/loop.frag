#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D PF_texture;
uniform float window_width;
uniform float window_height;


vec4 getTextureValue(float dx, float dy) {
    // get new coords
    return texture(
        PF_texture, 
        vec2(
            (TexCoord.s + (dx / window_width)), 
            (TexCoord.t + (dy / window_height))
        )
    );
}

void main() {
    vec4 texleft = getTextureValue(-3.0, -3.0);

    if (texleft.r >= 0.1){
        FragColor = texleft;
    }
    else {
        FragColor = getTextureValue(0., 0.);
    }
}