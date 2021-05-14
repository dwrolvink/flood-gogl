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
    vec4 texleft = getTextureValue(-2.5, -1.0);
    vec4 tex = getTextureValue(0., 0.);

    if (texleft.r >= 0.1){
        FragColor = vec4((texleft.r + tex.r)/2.0 + 0.01, texleft.g, texleft.b, texleft.a);
    }
    else {
        FragColor = getTextureValue(0., 0.);
    }
}