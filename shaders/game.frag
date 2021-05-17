#version 330 core
out vec4 FragColor;
in vec2 TexCoord;


uniform sampler2D PfGameTexture;
uniform sampler2D PfSmellRedTexture;
uniform sampler2D PfSmellGreenTexture;

uniform float window_width;
uniform float window_height;

uniform vec2 Actor1;
uniform float Actor1Radius;

float max5 (float v0, float v1, float v2, float v3, float v4) {
  return max( max( max( max(v0, v1), v2), v3), v4);
}

float rand(vec2 co){
    return fract(sin(dot(co, vec2(12.9898, 78.233))) * 43758.5453);
}

vec4 average(sampler2D image, vec2 uv, vec2 resolution) {
    float circular = 1.0/5.0;
    
    vec2 off1 = vec2(1.3846153846) * vec2(0.0, 1.0);
    vec2 off2 = vec2(1.3846153846) * vec2(0.0, -1.0);
    vec2 off3 = vec2(1.3846153846) * vec2(1.0, 0.0);
    vec2 off4 = vec2(1.3846153846) * vec2(-1.0, 0.0);

    vec2 off5 = vec2(1.3846153846) * vec2(1.0, 1.0);
    vec2 off6 = vec2(1.3846153846) * vec2(1.0, -1.0);
    vec2 off7 = vec2(1.3846153846) * vec2(-1.0, 1.0);
    vec2 off8 = vec2(1.3846153846) * vec2(-1.0, -1.0);    

    // get pixel colors
    vec4 color = vec4(0.0);
    color += texture2D(image, uv) * circular;
    color += texture2D(image, uv + (off1 / resolution)) * circular;
    color += texture2D(image, uv + (off2 / resolution)) * circular;
    color += texture2D(image, uv + (off3 / resolution)) * circular;
    color += texture2D(image, uv + (off4 / resolution)) * circular;

    return color;
}

float ToroidalDistance (vec2 P1, vec2 P2)
{
    float dx = abs(P2.s - P1.s);
    float dy = abs(P2.t - P1.t);
 
    if (dx > 0.5)
        dx = 1.0 - dx;
 
    if (dy > 0.5)
        dy = 1.0 - dy;
 
    return sqrt(dx*dx + dy*dy);
}

void main() {
    vec2 resolution = vec2(window_width, window_height);
    float growth = 1.3;
    float rnd = rand(TexCoord);    

    // Get average of 5 cells
    vec4 GameColor = average(PfGameTexture, TexCoord, resolution);

    // blue lighting where battle takes place
    if (GameColor.r > 0.01 && GameColor.g > 0.01){
        GameColor.b = 10.0*abs(GameColor.g - GameColor.r);
    }
    GameColor.b -= 0.002;

    // subtract red from green and vice versa to have them battle it out
    if (GameColor.r > GameColor.g){
        // big win
        if (GameColor.r - GameColor.g > 0.5){
            GameColor.r = GameColor.r * growth ;
            GameColor.g = 0;
        }
        else {
            GameColor.r = (GameColor.r - GameColor.g) * growth ;
            GameColor.g = 0;
        }
    } 
    else {
        // big win
        if (GameColor.g - GameColor.r > 0.8){
            GameColor.g = GameColor.g * growth ;
            GameColor.r = 0;
        }
        else {
            GameColor.g = (GameColor.g - GameColor.r) * growth ;
            GameColor.r = 0;
        }
    }

    GameColor.a = 1.0;

    float d = ToroidalDistance(TexCoord, Actor1);
    float r = Actor1Radius;
    float sharpness = 100.0 * r;
    if (d <= r){        
        float c = 1 - sharpness*max(0, (d - r)/r);
            GameColor.r = 0.0;
            GameColor.g = c;
            GameColor.b = 1.0;
    }
    

    // Done
    FragColor = GameColor;
    
    
}