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

vec2 coords_to_uv(vec2 coords) {
    return (coords + .5) / Resolution;
}

vec4 get_color_by_coord(sampler2D tex, vec2 coords) {
    return texture2D(tex, coords_to_uv(coords));
}



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


vec2 pick_target2 (int self, int enemy, vec2 coords) {
    int s = self;
    int e = enemy;

    vec4 value = get_color_by_coord(PfGameTexture, coords);

    // pick no target if no red in cell
    if (value[s] < 0.5) {
        return coords;
    }

    // -- get values
    vec4 value_top = get_color_by_coord(PfGameTexture, coords + Top);
    vec4 value_bottom = get_color_by_coord(PfGameTexture, coords + Bottom);
    vec4 value_left = get_color_by_coord(PfGameTexture, coords + Left);
    vec4 value_right = get_color_by_coord(PfGameTexture, coords + Right);


    //vec4 es_top; vec4 es_bottom; vec4 es_left; vec4 es_right;
    // if (enemy == GREEN) {
        vec4 es_top = get_color_by_coord(PfSmellGreenTexture, coords + Top);
        vec4 es_bottom = get_color_by_coord(PfSmellGreenTexture, coords + Bottom);
        vec4 es_left = get_color_by_coord(PfSmellGreenTexture, coords + Left);
        vec4 es_right = get_color_by_coord(PfSmellGreenTexture, coords + Right);
    // }
    // if (enemy == RED) {
    //     es_top = get_color_by_coord(PfSmellRedTexture, coords + Top);
    //     es_bottom = get_color_by_coord(PfSmellRedTexture, coords + Bottom);
    //     es_left = get_color_by_coord(PfSmellRedTexture, coords + Left);
    //     es_right = get_color_by_coord(PfSmellRedTexture, coords + Right);
    // }


    // -- types
    vec2 target_coords;
    float min_value;
    float max_value;



    // -- get empty cell
    if (value_left[s] == 0.0 && value_left[e] == 0.0) {
        return coords + Left;
    }
    if (value_right[s] == 0.0 && value_right[e] == 0.0) {
        return coords + Right;
    }
    if (value_top[s] == 0.0 && value_top[e] == 0.0) {
        return coords + Top;
    }
    if (value_bottom[s] == 0.0 && value_bottom[e] == 0.0) {
        return coords + Bottom;
    }


    // -- get cell with strongest enemy cell smell
    target_coords = coords;
    max_value = 0.0;


    if (es_top[e] > max_value) {
        target_coords = coords + Top;
        max_value = es_top[e];
    }
    if (es_bottom[e] > max_value) {
        target_coords = coords + Bottom;
        max_value = es_bottom[e];
    }
    if (es_left[e] > max_value) {
        target_coords = coords + Left;
        max_value = es_left[e];
    }
    if (es_right[e] > max_value) {
        target_coords = coords + Right;
        max_value = es_right[e];
    }

    if (max_value > 0.0) {
        return target_coords;
    }

    return coords;


    // // -- get cell with lowest own value
    // target_coords = coords;
    // min_value = 1000.0;

    // if (value_top[s] < min_value) {
    //     target_coords = coords + Top;
    //     min_value = value_top[s];
    // }
    // if (value_bottom[s] < min_value) {
    //     target_coords = coords + Bottom;
    //     min_value = value_bottom[s];
    // }
    // if (value_left[s] < min_value) {
    //     target_coords = coords + Left;
    //     min_value = value_left[s];
    // }
    // if (value_right[s] < min_value) {
    //     target_coords = coords + Right;
    //     min_value = value_right[s];
    // }

    // return target_coords;
}

float strat3 (int self, int enemy) {
    int s = self;
    int e = enemy;

    vec2 uv = TexCoord;
    vec2 coords = uv_to_coords(uv);

    vec4 cell = get_color_by_coord(PfGameTexture, coords);
    float g = 1.01;
    float a = 0.5;

    vec2 test;
    // send
    test = pick_target2(self, enemy, coords);
    if (test.x != coords.x || test.y != coords.y) {
        cell[s] = max(0.0, cell[s] - a);
    }

    // receive
    test = pick_target2(self, enemy, coords + Top);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + a);
    }
    test = pick_target2(self, enemy, coords + Bottom);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + a);
    }
    test = pick_target2(self, enemy, coords + Left);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + a);
    }
    test = pick_target2(self, enemy, coords + Right);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + a);
    }


    // return
    return cell[s] * g;
}


void main() {
    vec4 GameColor = texture2D(PfGameTexture, TexCoord);

    GameColor.g = strat2(GREEN);
    GameColor.r = strat3(RED, GREEN);

    // GameColor.b = get_color_by_coord(PfSmellGreenTexture, uv_to_coords(TexCoord)).g;

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