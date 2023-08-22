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
    float effect = 0.99;//0.985;
    float circular = 0.99;

    vec4 color = vec4(0.0, 0.0, 0.0, 1.0);
    vec4 color_smell = vec4(0.0);
    vec4 color_game = vec4(0.0);
    vec4 s_color;
    vec4 g_color;

    vec2 off1 = vec2(1.3846153846) * vec2(0.0, 1.0);
    vec2 off2 = vec2(1.3846153846) * vec2(0.0, -1.0);
    vec2 off3 = vec2(1.3846153846) * vec2(1.0, 0.0);
    vec2 off4 = vec2(1.3846153846) * vec2(-1.0, 0.0);

    vec2 off5 = vec2(1.3846153846) * vec2(1.0, 1.0);
    vec2 off6 = vec2(1.3846153846) * vec2(1.0, -1.0);
    vec2 off7 = vec2(1.3846153846) * vec2(-1.0, 1.0);
    vec2 off8 = vec2(1.3846153846) * vec2(-1.0, -1.0);    

    // get pixel colors 
    color_game += texture2D(PfGameTexture, uv);
    color_smell += texture2D(PfSmellTexture, uv);

    // get pixel colors - smell
    vec4 s_color1 = texture2D(PfSmellTexture, uv + (off1 / resolution));
    vec4 s_color2 = texture2D(PfSmellTexture, uv + (off2 / resolution));
    vec4 s_color3 = texture2D(PfSmellTexture, uv + (off3 / resolution));
    vec4 s_color4 = texture2D(PfSmellTexture, uv + (off4 / resolution));

    vec4 s_color5 = texture2D(PfSmellTexture, uv + (off5 / resolution)) * circular;
    vec4 s_color6 = texture2D(PfSmellTexture, uv + (off6 / resolution)) * circular;
    vec4 s_color7 = texture2D(PfSmellTexture, uv + (off7 / resolution)) * circular;
    vec4 s_color8 = texture2D(PfSmellTexture, uv + (off8 / resolution)) * circular;  

    // get pixel colors - game
    vec4 g_color1 = texture2D(PfGameTexture, uv + (off1 / resolution));
    vec4 g_color2 = texture2D(PfGameTexture, uv + (off2 / resolution));
    vec4 g_color3 = texture2D(PfGameTexture, uv + (off3 / resolution));
    vec4 g_color4 = texture2D(PfGameTexture, uv + (off4 / resolution));

    vec4 g_color5 = texture2D(PfGameTexture, uv + (off5 / resolution)) * circular;
    vec4 g_color6 = texture2D(PfGameTexture, uv + (off6 / resolution)) * circular;
    vec4 g_color7 = texture2D(PfGameTexture, uv + (off7 / resolution)) * circular;
    vec4 g_color8 = texture2D(PfGameTexture, uv + (off8 / resolution)) * circular;       

    // find highest value that was retrieved
    float max_green;
    max_green = max5(color.g, s_color1.g, s_color2.g, s_color3.g, s_color4.g);
    max_green = max5(max_green, s_color5.g, s_color6.g, s_color7.g, s_color8.g);
    max_green = max5(max_green, g_color1.g, g_color2.g, g_color3.g, g_color4.g);
    max_green = max5(max_green, g_color5.g, g_color6.g, g_color7.g, g_color8.g);

    // Set output color
    color.g = max_green;
    color.g *= effect;

    // Evaporate    
    color.g -= 0.002;    

      // alpha fix
    color.a = 1.0;

    return color;
}

void main() {
    vec2 resolution = vec2(window_width, window_height);
    vec4 SmellColor = blur(PfSmellTexture, TexCoord, resolution);

    FragColor = SmellColor;
    
}