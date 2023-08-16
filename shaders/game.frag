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

// settings
float MinGrowth = 1.01;

// reusables
vec2 Top = vec2(0.0, 1.0);
vec2 Bottom = vec2(0, -1.0);
vec2 Left = vec2(-1, 0);
vec2 Right = vec2(1, 0);


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
vec2 pick_highest_value_coords(sampler2D tex, vec2 coords) {
    vec4 top_value = get_color_by_coord(tex, coords + Top);
    vec4 bottom_value = get_color_by_coord(tex, coords + Bottom);
    vec4 left_value = get_color_by_coord(tex, coords + Left);
    vec4 right_value = get_color_by_coord(tex, coords + Right);

    float max_value = 0.0;
    vec2 receiver_coords = coords;

    if (top_value.g > max_value) {
        max_value = top_value.g;
        receiver_coords = coords + Top;
    }
    if (bottom_value.g > max_value) {
        max_value = bottom_value.g;
        receiver_coords = coords + Bottom;
    }
    if (left_value.g > max_value) {
        max_value = left_value.g;
        receiver_coords = coords + Left;
    }
    if (right_value.g > max_value) {
        max_value = right_value.g;
        receiver_coords = coords + Right;
    }    

    return receiver_coords;
}

// vec2 pick_lowest_value_coords(sampler2D tex, vec2 coords) {
//     vec4 top_value = get_color_by_coord(tex, coords + Top);
//     vec4 bottom_value = get_color_by_coord(tex, coords + Bottom);
//     vec4 left_value = get_color_by_coord(tex, coords + Left);
//     vec4 right_value = get_color_by_coord(tex, coords + Right);

//     float min_value = 1.0;
//     vec2 receiver_coords = coords;

//     if (top_value.r < min_value) {
//         min_value = top_value.r;
//         receiver_coords = Top;
//     }
//     if (bottom_value.r < min_value) {
//         min_value = bottom_value.r;
//         receiver_coords = Bottom;
//     }
//     if (left_value.r < min_value) {
//         min_value = left_value.r;
//         receiver_coords = Left;
//     }
//     if (right_value.r < min_value) {
//         min_value = right_value.r;
//         receiver_coords = Right;
//     }    

//     return receiver_coords;
// }

// vec2 pick_receiver_coords(vec2 coords) {
//     // pick the cell with the strongest enemy smell
//     // vec2 receiver_coords = pick_highest_value_coords(PfSmellGreenTexture, coords);
//     // if (receiver_coords != coords){
//     //     return receiver_coords;
//     // }
    

//     // pick cell with lowest own value
//     vec2 receiver_coords = pick_lowest_value_coords(PfGameTexture, coords);
//     if (receiver_coords != coords){
//         return receiver_coords;
//     }

//     return coords;
// }

vec2 pick_bottom(vec2 coords) {
    vec4 color = get_color_by_coord(PfGameTexture, coords);
    if (color.r > 0.){
        return coords + Bottom;
    }
    return coords + Left;


}

vec2 pick_top(vec2 coords) {
    vec4 color_bottom = get_color_by_coord(PfGameTexture, coords + Top);
    if (color_bottom.r > 0.){
        return coords + Top;
    }
    return coords + Left;
}

float action_red () {
    vec2 uv = TexCoord;
    vec2 coords = uv_to_coords(uv);

    //vec4 cell = value_at(PfGameTexture, uv);
    vec4 cell = get_color_by_coord(PfGameTexture, coords);

    float send_threshold = 0.4;

    // [send] don't send if not at the minimum
    // if (cell.r > send_threshold) {
    //     vec2 receiver_coords = pick_receiver_coords(uv);
    //     if (equals(mock_receiver_coords,uv)) {
    //         vec4 receiver = value_at(PfGameTexture, receiver_coords);
    //         float max_send = (1.0 - receiver.r) * 0.25;
    //         float min_send = cell.r - 0.05;
    //         cell.r -= min(min_send, max_send);
    //     }
    // }

    // [receive] for each neighbour, we will check if we will receive any of them and add it to our total
    vec4 receiver = get_color_by_coord(PfGameTexture, coords);

    vec4 sender = get_color_by_coord(PfGameTexture, coords + Top + Right);

    vec2 picked_coords = pick_bottom(coords + Top);
    if (picked_coords == coords) {
        vec4 color = get_color_by_coord(PfGameTexture, coords + Top);
        if (color.r > 0.0) {
            cell.r += color.r;
        }
    }

    picked_coords = pick_bottom(coords + Left);
    if (picked_coords == coords) {
        vec4 color = get_color_by_coord(PfGameTexture, coords + Left);
        if (color.r > 0.0) {
            cell.r += color.r;
        }
    }

    // picked_coords = pick_top(coords + Bottom);
    // if (picked_coords == coords) {
    //     vec4 color = get_color_by_coord(PfGameTexture, coords + Bottom);
    //     if (color.r > 0.0) {
    //         cell.r += color.r;
    //     }
    // }



    // if (sender.r > send_threshold) {
    //     vec2 mock_receiver_coords = pick_receiver_coords(coords + Top);
    //     if (equals(mock_receiver_coords,coords)){
    //         float max_send = (1.0 - receiver.r) * 0.25;
    //         float min_send = sender.r - 0.05;
    //         cell.r += min(min_send, max_send);
    //     }
    // }

    // sender = get_color_by_coord(PfGameTexture, coords + Bottom);
    // if (sender.r > send_threshold) {
    //     vec2 mock_receiver_coords = pick_receiver_coords(coords + Bottom);
    //     if (equals(mock_receiver_coords,coords)){
    //         float max_send = (1.0 - receiver.r) * 0.25;
    //         float min_send = sender.r - 0.05;
    //         cell.r += min(min_send, max_send);
    //     }
    // }

    // sender = get_color_by_coord(PfGameTexture, coords + Left);
    // if (sender.r > send_threshold) {
    //     vec2 mock_receiver_coords = pick_receiver_coords(coords + Left);
    //     if (equals(mock_receiver_coords,coords)){
    //         float max_send = (1.0 - receiver.r) * 0.25;
    //         float min_send = sender.r - 0.05;
    //         cell.r += min(min_send, max_send);
    //     }
    // }

    // sender = get_color_by_coord(PfGameTexture, coords + Right);
    // if (sender.r > send_threshold) {
    //     vec2 mock_receiver_coords = pick_receiver_coords(coords + Right);
    //     if (equals(mock_receiver_coords,coords)){
    //         float max_send = (1.0 - receiver.r) * 0.25;
    //         float min_send = sender.r - 0.05;
    //         cell.r += min(min_send, max_send);
    //     }
    // }

    // growth 
    float growth_factor = 1.0 + (0.01 * cell.r);
    growth_factor = max(growth_factor, MinGrowth);
    cell.r = cell.r * growth_factor;
    
    // return
    return cell.r;
}


float action_green () {
    vec4 GameColor = average(PfGameTexture, TexCoord);

    // growth 
    float growth_factor = 1.0 + (0.1 * GameColor.g);
    growth_factor = max(growth_factor, MinGrowth);
    GameColor.g = GameColor.g * growth_factor;

    return GameColor.g;
}

// float action_red (vec2 uv) {
//     vec4 GameColor = average(PfGameTexture, uv);

//     // growth 
//     float growth_factor = 1.0 + (0.1 * GameColor.r);
//     growth_factor = max(growth_factor, MinGrowth);
//     GameColor.r = GameColor.r * growth_factor;

//     return GameColor.r;
// }

void main() {
    // float rnd = rand(TexCoord);    

    // Get average of 5 cells
    // vec4 GameColor = average(PfGameTexture, TexCoord);
    vec4 GameColor = value_at(PfGameTexture, TexCoord);
    // GameColor = GameColor * vec4(1.3, 1.3, 0.99, 1.0);
    

    // if (iTime >= 0 && iTime < 10) {
    //     GameColor.r = action_red(TexCoord + Bottom + Right);
    // }
    // else if (iTime >= 10 && iTime < 20) {
    //     GameColor.r = action_red(TexCoord + Bottom);
    // }
    // else if (iTime >= 20 && iTime < 30) {
    //     GameColor.r = action_red(TexCoord + Bottom + Left);
    // }
    // else if (iTime >= 30 && iTime < 40) {
    //     GameColor.r = action_red(TexCoord + Left);
    // }
    // else if (iTime >= 40 && iTime < 50) {
    //     GameColor.r = action_red(TexCoord + Left + Top);
    // }
    // else if (iTime >= 50 && iTime < 60) {
    //     GameColor.r = action_red(TexCoord + Top);
    // }
    // else if (iTime >= 60 && iTime < 70) {
    //     GameColor.r = action_red(TexCoord + Top + Right);
    // }
    // else if (iTime >= 70 && iTime < 80) {
    //     GameColor.r = action_red(TexCoord + Right);
    // }
    // else if (iTime >= 80 && iTime < 90) {
    //     GameColor.r = action_red(TexCoord + Left);
    // }
    // else if (iTime >= 90 && iTime < 100) {
    //     GameColor.r = action_red(TexCoord + Left + Top);
    // }


    GameColor.g = action_green();
    GameColor.r = action_red();

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