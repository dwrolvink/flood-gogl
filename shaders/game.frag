#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

// input
uniform sampler2D PfGameTexture; 
uniform sampler2D PfSmellRedTexture; 
uniform sampler2D PfSmellGreenTexture; 

uniform int iTime;
uniform float window_width;
uniform float window_height;

uniform vec2 Actor1;
uniform float Actor1Radius;

vec2 Resolution = vec2(window_width, window_height);

int RED = 0;
int GREEN = 1;

// settings
float MinGrowth = 1.01;

// reusables
vec2 Top = vec2(0.0, 1.0);
vec2 Bottom = vec2(0.0, -1.0);
vec2 Left = vec2(-1.0, 0.0);
vec2 Right = vec2(1.0, 0.0);


// helper functions
// float max5 (float v0, float v1, float v2, float v3, float v4) {
//   return max( max( max( max(v0, v1), v2), v3), v4);
// }

// float rand(vec2 co){
//     return fract(sin(dot(co, vec2(12.9898, 78.233))) * 43758.5453);
// }

vec2 uv_to_coords(vec2 uv) {
    vec2 coords = (TexCoord * Resolution);
    coords.x = floor(coords.x);
    coords.y = floor(coords.y);
    
    return coords;
}

vec4 get_color_by_coord(sampler2D tex, vec2 coords) {
// vec2 uv = (vec2(pixelX, pixelY) + .5) / resolutionOfTexture;
    vec2 uv = (coords + .5) / Resolution;
    return texture2D(tex, uv);
}

bool equals(vec2 pos1, vec2 pos2) {
    if (abs(abs(pos1.x) - abs(pos2.x)) > 400.){
        return false;
    }
    if (abs(abs(pos1.y) - abs(pos2.y)) > 4000.){
        return false;
    }
    return true;
}

vec4 value_at(sampler2D tex, vec2 uv) {
    return texture2D(tex, uv);
}

vec4 average(sampler2D image, vec2 uv) {
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
    color += texture2D(image, uv + (off1 / Resolution)) * circular;
    color += texture2D(image, uv + (off2 / Resolution)) * circular;
    color += texture2D(image, uv + (off3 / Resolution)) * circular;
    color += texture2D(image, uv + (off4 / Resolution)) * circular;

    return color;
}


// float ToroidalDistance (vec2 P1, vec2 P2)
// {
//     float dx = abs(P2.s - P1.s);
//     float dy = abs(P2.t - P1.t);
 
//     if (dx > 0.5)
//         dx = 1.0 - dx;
 
//     if (dy > 0.5)
//         dy = 1.0 - dy;
 
//     return sqrt(dx*dx + dy*dy);
// }

// game logic

vec2 pick_target (int c, vec2 coords) {
    vec4 value = get_color_by_coord(PfGameTexture, coords);

    // pick no target if no red in cell
    if (value[c] < 0.5) {
        return coords;
    }

    vec4 value_top = get_color_by_coord(PfGameTexture, coords + Top);
    vec4 value_bottom = get_color_by_coord(PfGameTexture, coords + Bottom);
    vec4 value_left = get_color_by_coord(PfGameTexture, coords + Left);
    vec4 value_right = get_color_by_coord(PfGameTexture, coords + Right);

    vec2 target_coords = coords;
    float min_value = 1000.0;

    if (value_top[c] < min_value) {
        target_coords = coords + Top;
        min_value = value_top[c];
    }
    if (value_bottom[c] < min_value) {
        target_coords = coords + Bottom;
        min_value = value_bottom[c];
    }
    if (value_left[c] < min_value) {
        target_coords = coords + Left;
        min_value = value_left[c];
    }
    if (value_right[c] < min_value) {
        target_coords = coords + Right;
        min_value = value_right[c];
    }

    return target_coords;
}


float strat2 (int c) {
    vec2 uv = TexCoord;
    vec2 coords = uv_to_coords(uv);

    vec4 cell = get_color_by_coord(PfGameTexture, coords);
    float g = 1.01;
    float a = 0.5;

    vec2 test;
    // send
    test = pick_target(c, coords);
    if (test.x != coords.x || test.y != coords.y) {
        cell[c] = max(0.0, cell[c] - a);
    }

    // receive
    test = pick_target(c, coords + Top);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }
    test = pick_target(c, coords + Bottom);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }
    test = pick_target(c, coords + Left);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }
    test = pick_target(c, coords + Right);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }


    // return
    return cell[c] * g;
}



vec2 pick_target2 (int c, vec2 coords) {
    vec4 value = get_color_by_coord(PfGameTexture, coords);

    // pick no target if no red in cell
    if (value[c] < 0.5) {
        return coords;
    }

    vec4 value_top = get_color_by_coord(PfGameTexture, coords + Top);
    vec4 value_bottom = get_color_by_coord(PfGameTexture, coords + Bottom);
    vec4 value_left = get_color_by_coord(PfGameTexture, coords + Left);
    vec4 value_right = get_color_by_coord(PfGameTexture, coords + Right);

    vec2 target_coords = coords;
    float min_value = 1000.0;

    if (value_top[c] < min_value) {
        target_coords = coords + Top;
        min_value = value_top[c];
    }
    if (value_bottom[c] < min_value) {
        target_coords = coords + Bottom;
        min_value = value_bottom[c];
    }
    if (value_left[c] < min_value) {
        target_coords = coords + Left;
        min_value = value_left[c];
    }
    if (value_right[c] < min_value) {
        target_coords = coords + Right;
        min_value = value_right[c];
    }

    return target_coords;
}

float strat3 (int c) {
    vec2 uv = TexCoord;
    vec2 coords = uv_to_coords(uv);

    vec4 cell = get_color_by_coord(PfGameTexture, coords);
    float g = 1.01;
    float a = 0.5;

    vec2 test;
    // send
    test = pick_target2(c, coords);
    if (test.x != coords.x || test.y != coords.y) {
        cell[c] = max(0.0, cell[c] - a);
    }

    // receive
    test = pick_target2(c, coords + Top);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }
    test = pick_target2(c, coords + Bottom);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }
    test = pick_target2(c, coords + Left);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }
    test = pick_target2(c, coords + Right);
    if (test.x == coords.x && test.y == coords.y) {
        cell[c] = min(1.0, cell[c] + a);
    }


    // return
    return cell[c] * g;
}


void main() {
    vec4 GameColor = value_at(PfGameTexture, TexCoord);

    GameColor.g = strat2(GREEN);
    GameColor.r = strat3(RED);

    // blue lighting where battle takes place
    if (GameColor.r > 0.01 && GameColor.g > 0.01){
        GameColor.b = 10.0*abs(GameColor.g - GameColor.r);
    }
    GameColor.b -= 0.002;

    // subtract red from green and vice versa to have them battle it out
    if (GameColor.r > GameColor.g){
        // big win
        if (GameColor.r - GameColor.g > 0.5){
            // GameColor.r = GameColor.r * growth ;
            GameColor.r = (GameColor.r - GameColor.g) ;
            GameColor.g = 0;
        }
        else {
            GameColor.r = (GameColor.r - GameColor.g) ;
            GameColor.g = 0;
        }
    } 
    else {
        // big win
        if (GameColor.g - GameColor.r > 0.8){
            // GameColor.g = GameColor.g * growth ;
            GameColor.g = (GameColor.g - GameColor.r);
            GameColor.r = 0;
        }
        else {
            GameColor.g = (GameColor.g - GameColor.r);
            GameColor.r = 0;
        }
    }

    GameColor.a = 1.0;


    // // draw player
    // float d = ToroidalDistance(TexCoord, Actor1);
    // float d2 = ToroidalDistance(TexCoord, vec2(0.5, 0.5));
    // float r = Actor1Radius;
    // float sharpness = 10.0 * r;
    // float c = 0.0;
    // if (d <= r){        
    //     c += 1.0 - sharpness*max(0.0, (d - r)/r);
    //     GameColor.r = 0.0;
    //     GameColor.g = c;
    //     GameColor.b = c;
    // } else if (d2 <= r){        
    //     c += 1.0 - sharpness*max(0.0, (d2 - r)/r);
    //     GameColor.r = 0.0;
    //     GameColor.g = c;
    //     GameColor.b = c;
    // } 
    

    // Done
    FragColor = GameColor;
    
    
}