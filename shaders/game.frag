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

// 0.5 * self + 1/16th (1,2,3,4,5,6,7,8) * 1.01

vec4 blur(sampler2D image, vec2 uv, vec2 resolution) {
    //float effect = 1.0;
    float circular = 1.0/8.0;

    vec4 color = vec4(0.0);
    //vec4 icolor;

    vec2 off1 = vec2(1.3846153846) * vec2(0.0, 1.0);
    vec2 off2 = vec2(1.3846153846) * vec2(0.0, -1.0);
    vec2 off3 = vec2(1.3846153846) * vec2(1.0, 0.0);
    vec2 off4 = vec2(1.3846153846) * vec2(-1.0, 0.0);

    vec2 off5 = vec2(1.3846153846) * vec2(1.0, 1.0);
    vec2 off6 = vec2(1.3846153846) * vec2(1.0, -1.0);
    vec2 off7 = vec2(1.3846153846) * vec2(-1.0, 1.0);
    vec2 off8 = vec2(1.3846153846) * vec2(-1.0, -1.0);    

    // get pixel colors
    color += texture2D(image, uv) * 0.5;

    color += texture2D(image, uv + (off1 / resolution)) * circular;
    color += texture2D(image, uv + (off2 / resolution)) * circular;
    color += texture2D(image, uv + (off3 / resolution)) * circular;
    color += texture2D(image, uv + (off4 / resolution)) * circular;

    //color += texture2D(image, uv + (off5 / resolution)) * circular;
    //color += texture2D(image, uv + (off6 / resolution)) * circular;
    //color += texture2D(image, uv + (off7 / resolution)) * circular;
    //color += texture2D(image, uv + (off8 / resolution)) * circular;    

    return color;
}

float rand(vec2 co){
    return fract(sin(dot(co, vec2(12.9898, 78.233))) * 43758.5453);
}

void main() {
    vec2 resolution = vec2(window_width, window_height);

    vec4 GameColor = blur(PfGameTexture, TexCoord, resolution);
    vec4 SelfColor = texture(PfGameTexture, TexCoord);

    float growth = 1.2;
    float rnd = rand(TexCoord);

    // blue lighting
    if (GameColor.r > 0.01 && GameColor.g > 0.01){
        GameColor.b = 10.0*abs(GameColor.g - GameColor.r);
    }
    GameColor.b -= 0.002;

    if (GameColor.r > GameColor.g){
        // big win
        if (GameColor.r - GameColor.g > 0.01){
            GameColor.r = GameColor.r * growth ;
            GameColor.g = 0;
        }
        else {
            GameColor.r = (GameColor.r - GameColor.g) * growth ;
            GameColor.g = 0;
        }
        //GameColor.b = (GameColor.r - GameColor.g);

    } else {
        // big win
        if (GameColor.g - GameColor.r > 0.8){
            GameColor.g = GameColor.g * growth ;
            GameColor.r = 0;
        }
        else {
            GameColor.g = (GameColor.g - GameColor.r) * growth ;
            GameColor.r = 0;
        }

        //GameColor.b = (GameColor.g - GameColor.r);
    }

    if (GameColor.r > 0. || GameColor.g > 0.){
        GameColor.a = 1.0;
    }

    //GameColor.b = SelfColor.b;
    //GameColor.a = 1.0;

    FragColor =  GameColor;
    
    
}