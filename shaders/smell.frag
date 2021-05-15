#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D PfGameTexture;
uniform sampler2D PfSmellTexture;

uniform float window_width;
uniform float window_height;

float max5 (float v0, float v1, float v2, float v3, float v4) {
  return max( max( max( max(v0, v1), v2), v3), v4);
}

vec4 blur(sampler2D image, vec2 uv, vec2 resolution) {
    float effect = 0.96;
    float circular = 0.982;
    vec4 color = vec4(0.0);
    vec4 icolor;

    vec2 off1 = vec2(1.3846153846) * vec2(0.0, 1.0);
    vec2 off2 = vec2(1.3846153846) * vec2(0.0, -1.0);
    vec2 off3 = vec2(1.3846153846) * vec2(1.0, 0.0);
    vec2 off4 = vec2(1.3846153846) * vec2(-1.0, 0.0);

    vec2 off5 = vec2(1.3846153846) * vec2(1.0, 1.0);
    vec2 off6 = vec2(1.3846153846) * vec2(1.0, -1.0);
    vec2 off7 = vec2(1.3846153846) * vec2(-1.0, 1.0);
    vec2 off8 = vec2(1.3846153846) * vec2(-1.0, -1.0);    

    // get pixel colors
    color += texture2D(image, uv);
    vec4 icolor1 = texture2D(image, uv + (off1 / resolution));
    vec4 icolor2 = texture2D(image, uv + (off2 / resolution));
    vec4 icolor3 = texture2D(image, uv + (off3 / resolution));
    vec4 icolor4 = texture2D(image, uv + (off4 / resolution));

    vec4 icolor5 = texture2D(image, uv + (off5 / resolution)) * circular;
    vec4 icolor6 = texture2D(image, uv + (off6 / resolution)) * circular;
    vec4 icolor7 = texture2D(image, uv + (off7 / resolution)) * circular;
    vec4 icolor8 = texture2D(image, uv + (off8 / resolution)) * circular;    

    // blur red
    float max_red;
    max_red = max5(color.r, icolor1.r, icolor2.r, icolor3.r, icolor4.r);
    max_red = max5(max_red, icolor5.r, icolor6.r, icolor7.r, icolor8.r);
    if (max_red != color.r){
        color.r = max_red * effect;
    }    

    // blur green
    float max_green;
    max_green = max5(color.g, icolor1.g, icolor2.g, icolor3.g, icolor4.g);
    max_green = max5(max_green, icolor5.g, icolor6.g, icolor7.g, icolor8.g);
    if (max_green != color.g){
        color.g = max_green * effect;
    }        

    // alpha fix, might do nothing
    float max_alpha = max5(color.a, icolor1.a, icolor2.a, icolor3.a, icolor4.a);
    color.a = max_alpha;

    return color;
}

void main() {
    vec2 resolution = vec2(window_width, window_height);
    vec4 SmellColor = blur(PfSmellTexture, TexCoord, resolution);
    vec4 GameColor = texture(PfGameTexture, TexCoord);

    if (GameColor.r > 0.01){
        FragColor = GameColor;
    } else {
        FragColor = SmellColor;
    }
}