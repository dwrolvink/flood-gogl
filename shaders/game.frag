#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

// input
uniform sampler2D PfGameTexture; 
uniform sampler2D PfSmellRedTexture; 
uniform sampler2D PfSmellGreenTexture; 

uniform float window_width;
uniform float window_height;

uniform vec2 Actor1;
uniform float Actor1Radius;
uniform int TIMESTAMP;

vec2 Resolution = vec2(window_width, window_height);

int MOD4;

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
int modulo(int a, int b) {
    return a - (b * int(a/b));
}

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
vec3 package(vec2 coords, float amount) {
    // amount to send to the coordinate
    return vec3(coords.x, coords.y, amount);
}
vec2 get_coords_from_package(vec3 package){
    return vec2(package.x, package.y);
}

vec3 pick_target_1 (int self, int enemy, vec2 coords) {
    int s = self;
    int e = enemy;

    vec4 value = get_color_by_coord(PfGameTexture, coords);

    float restitution = 0.1;                            // how much to keep in the current cell
    float min_send = 0.7 - restitution;                 // don't send too small values, as rounding errors will cause loss
    float send = max(0.0, value[s] - restitution);      // how much can be sent by this cell

    if (send < min_send) {                              // don't send too small values, as rounding errors will cause loss
        return package(coords, 0.0);
    }

    MOD4 = modulo(TIMESTAMP + (2*int(coords.x)) + int(coords.y), 4);


    // -- get values in random order
    int index_0 = modulo(MOD4 + 0, 4);
    int index_1 = modulo(MOD4 + 1, 4);
    int index_2 = modulo(MOD4 + 2, 4);
    int index_3 = modulo(MOD4 + 3, 4);

    vec2 icoords[4];
    icoords[index_0] = coords + Top;
    icoords[index_1] = coords + Bottom;
    icoords[index_2] = coords + Left;
    icoords[index_3] = coords + Right;

    vec4 values[4];
    values[index_0] = get_color_by_coord(PfGameTexture, coords + Top);
    values[index_1] = get_color_by_coord(PfGameTexture, coords + Bottom);
    values[index_2] = get_color_by_coord(PfGameTexture, coords + Left);
    values[index_3] = get_color_by_coord(PfGameTexture, coords + Right);


    // -- types
    vec2 target_coords;
    float min_value;
    float max_value;

    // get smell values
    vec4 es_top; vec4 es_bottom; vec4 es_left; vec4 es_right;
    vec4 enemy_smell[4];
    if (enemy == GREEN) {
        enemy_smell[index_0] = get_color_by_coord(PfSmellGreenTexture, icoords[index_0]);
        enemy_smell[index_1] = get_color_by_coord(PfSmellGreenTexture, icoords[index_1]);
        enemy_smell[index_2] = get_color_by_coord(PfSmellGreenTexture, icoords[index_2]);
        enemy_smell[index_3] = get_color_by_coord(PfSmellGreenTexture, icoords[index_3]);
    }
    if (enemy == RED) {
        enemy_smell[index_0] = get_color_by_coord(PfSmellRedTexture, icoords[index_0]);
        enemy_smell[index_1] = get_color_by_coord(PfSmellRedTexture, icoords[index_1]);
        enemy_smell[index_2] = get_color_by_coord(PfSmellRedTexture, icoords[index_2]);
        enemy_smell[index_3] = get_color_by_coord(PfSmellRedTexture, icoords[index_3]);
    }


    // // die if surrounded by max value nbs
    // int cells_with_max_value = 0;
    // if (values[0][s] > 0.9) {
    //     cells_with_max_value += 1;
    // }
    // if (values[1][s] > 0.9) {
    //     cells_with_max_value += 1;
    // }
    // if (values[2][s] > 0.9) {
    //     cells_with_max_value += 1;
    // }
    // if (values[3][s] > 0.9) {
    //     cells_with_max_value += 1;
    // }

    // if (cells_with_max_value == 4) {
    //     // die
    //     return package(coords + vec2(2.0, 2.0), value[s]);
    // }


    // -- get empty cell
    if (values[0][s] <= 0.01 && values[0][e] == 0.0) {
        return package(icoords[0], send);
    }
    if (values[1][s] <= 0.01 && values[1][e] == 0.0) {
        return package(icoords[1], send);
    }
    if (values[2][s] <= 0.01 && values[2][e] == 0.0) {
        return package(icoords[2], send);
    }
    if (values[3][s] <= 0.01 && values[3][e] == 0.0) {
        return package(icoords[3], send);
    }

    // -- get cell with strongest enemy cell value
    target_coords = coords;
    max_value = 0.0;

    if (values[0][e] > max_value) {
        target_coords = icoords[0];
        max_value = values[0][e];
    }
    if (values[1][e] > max_value) {
        target_coords = icoords[1];
        max_value = values[1][e];
    }
    if (values[2][e] > max_value) {
        target_coords = icoords[2];
        max_value = values[2][e];
    }
    if (values[3][e] > max_value) {
        target_coords = icoords[3];
        max_value = values[3][e];
    } 

    if (max_value > 0.0) {
        return package(target_coords, send);
    }

    // -- get cell with strongest enemy cell smell
    target_coords = coords;
    max_value = 0.0;

    if (enemy_smell[0][e] > max_value) {
        target_coords = icoords[0];
        max_value = enemy_smell[0][e];
    }
    if (enemy_smell[1][e] > max_value) {
        target_coords = icoords[1];
        max_value = enemy_smell[1][e];
    }
    if (enemy_smell[2][e] > max_value) {
        target_coords = icoords[2];
        max_value = enemy_smell[2][e];
    }
    if (enemy_smell[3][e] > max_value) {
        target_coords = icoords[3];
        max_value = enemy_smell[3][e];
    }    

    if (max_value > 0.0) {
        return package(target_coords, send);
    }


    // // -- get cell with lowest own value
    // target_coords = coords;
    // min_value = 1.0;

    // if (values[0][s] > 0.0 && values[0][s] < min_value && values[0][s] + send < 1.0) {
    //     target_coords = icoords[0];
    //     min_value = values[0][s];
    // }
    // if (values[1][s] > 0.0 && values[1][s] < min_value && values[1][s] + send < 1.0) {
    //     target_coords = icoords[1];
    //     min_value = values[1][s];
    // }
    // if (values[2][s] > 0.0 && values[2][s] < min_value && values[2][s] + send < 1.0) {
    //     target_coords = icoords[2];
    //     min_value = values[2][s];
    // }
    // if (values[3][s] > 0.0 && values[3][s] < min_value && values[3][s] + send < 1.0) {
    //     target_coords = icoords[3];
    //     min_value = values[3][s];
    // } 

    // if (min_value < 1.0) {
    //     return package(target_coords, send);
    // } 


    // ?? do nothing
    return package(coords, 0.0);
}

float strat_1 (int self, int enemy) {
    int s = self;
    int e = enemy;

    vec2 uv = TexCoord;
    vec2 coords = uv_to_coords(uv);
    vec4 cell = get_color_by_coord(PfGameTexture, coords);

    float cell_start = cell[s];


    // growth multiplier
    float g = 1.01;

    // growth adder
    float a = 0.001;


    vec3 test;
    vec2 test_coords;
    float test_val;

    // send
    test = pick_target_1(self, enemy, coords);
    cell[s] = max(0.0, cell[s] - test.z);
    
    // receive
    test = pick_target_1(self, enemy, coords + Top);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z); // bug: can lead to loss of amount
    }
    test = pick_target_1(self, enemy, coords + Bottom);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z);
    }
    test = pick_target_1(self, enemy, coords + Left);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z);
    }
    test = pick_target_1(self, enemy, coords + Right);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z);
    }

    // grow
    float val = cell[s];
    a += min(0.1, max(0.0, (cell[s] - cell_start)/cell_start));
    if (cell[s] > 0.0) {
        val = max(0.0, min(1.0, cell[s] * g + a ));
    }
    return val;
}


vec3 pick_target_2 (int self, int enemy, vec2 coords) {
    /*
        if no smell:
          - if 2 neighbours with max value:
            - [send] to lowest neighbor (any)
          - else
             [send to cell with highest value]

        - send to cell of own color with lowest value
        - 


    */
    int s = self;
    int e = enemy;

    vec4 value = get_color_by_coord(PfGameTexture, coords);

    float restitution = 0.2 ;                            // how much to keep in the current cell
    float min_send = 0.3;                                // don't send too small values, as rounding errors will cause loss
    float min_send_new_cell = 0.6;                       // to colonize an empty cell we need at least this amount (for this strat)
    float send = max(0.0, value[s] - restitution);       // how much can be sent by this cell

    if (send < min_send) {                              // don't send too small values, as rounding errors will cause loss
        return package(coords, 0.0);
    }

    MOD4 = modulo(TIMESTAMP + int(coords.x) + int(coords.y), 4);


    // -- get values in random order
    int index_0 = modulo(MOD4 + 0, 4);
    int index_1 = modulo(MOD4 + 1, 4);
    int index_2 = modulo(MOD4 + 2, 4);
    int index_3 = modulo(MOD4 + 3, 4);

    vec2 icoords[4];
    icoords[index_0] = coords + Top;
    icoords[index_1] = coords + Bottom;
    icoords[index_2] = coords + Left;
    icoords[index_3] = coords + Right;

    vec4 values[4];
    values[index_0] = get_color_by_coord(PfGameTexture, coords + Top);
    values[index_1] = get_color_by_coord(PfGameTexture, coords + Bottom);
    values[index_2] = get_color_by_coord(PfGameTexture, coords + Left);
    values[index_3] = get_color_by_coord(PfGameTexture, coords + Right);


    // -- types
    vec2 target_coords;
    float min_value;
    float max_value;

    // get smell values
    vec4 es_top; vec4 es_bottom; vec4 es_left; vec4 es_right;
    vec4 enemy_smell[4];
    if (enemy == GREEN) {
        enemy_smell[index_0] = get_color_by_coord(PfSmellGreenTexture, icoords[index_0]);
        enemy_smell[index_1] = get_color_by_coord(PfSmellGreenTexture, icoords[index_1]);
        enemy_smell[index_2] = get_color_by_coord(PfSmellGreenTexture, icoords[index_2]);
        enemy_smell[index_3] = get_color_by_coord(PfSmellGreenTexture, icoords[index_3]);
    }
    if (enemy == RED) {
        enemy_smell[index_0] = get_color_by_coord(PfSmellRedTexture, icoords[index_0]);
        enemy_smell[index_1] = get_color_by_coord(PfSmellRedTexture, icoords[index_1]);
        enemy_smell[index_2] = get_color_by_coord(PfSmellRedTexture, icoords[index_2]);
        enemy_smell[index_3] = get_color_by_coord(PfSmellRedTexture, icoords[index_3]);
    }

    // NO SMELL
    // --------------
    if (enemy_smell[0][e] + enemy_smell[1][e] + enemy_smell[2][e] + enemy_smell[3][e] == 0.0){
        int cells_with_max_value = 0;
        if (values[0][s] > 0.9) {
            cells_with_max_value += 1;
        }
        if (values[1][s] > 0.9) {
            cells_with_max_value += 1;
        }
        if (values[2][s] > 0.9) {
            cells_with_max_value += 1;
        }
        if (values[3][s] > 0.9) {
            cells_with_max_value += 1;
        }

        if (cells_with_max_value == 4) {
            // die
            return package(coords + vec2(2.0, 2.0), value[s]);
        }

        if (cells_with_max_value >= 2) {
            // get cell with lowest value
            target_coords = coords;
            min_value = 1.0;

            if (values[0][s] < min_value && values[0][s] + send < 1.0) {
                target_coords = icoords[0];
                min_value = values[0][s];
            }
            if (values[1][s] < min_value && values[1][s] + send < 1.0) {
                target_coords = icoords[1];
                min_value = values[1][s];
            }
            if (values[2][s] < min_value && values[2][s] + send < 1.0) {
                target_coords = icoords[2];
                min_value = values[2][s];
            }
            if (values[3][s] < min_value && values[3][s] + send < 1.0) {
                target_coords = icoords[3];
                min_value = values[3][s];
            } 

            if (min_value < 1.0) {
                return package(target_coords, send);
            } 
        }
    }



    // -- get empty cell
    // if (send >= min_send_new_cell) {
        if (values[0][s] <= 0.01 && values[0][e] == 0.0) {
            return package(icoords[0], send);
        }
        if (values[1][s] <= 0.01 && values[1][e] == 0.0) {
            return package(icoords[1], send);
        }
        if (values[2][s] <= 0.01 && values[2][e] == 0.0) {
            return package(icoords[2], send);
        }
        if (values[3][s] <= 0.01 && values[3][e] == 0.0) {
            return package(icoords[3], send);
        }
    // }


    // // -- get cell with strongest enemy cell smell
    target_coords = coords;
    max_value = 0.0;

    if (enemy_smell[0][e] > max_value) {
        target_coords = icoords[0];
        max_value = enemy_smell[0][e];
    }
    if (enemy_smell[1][e] > max_value) {
        target_coords = icoords[1];
        max_value = enemy_smell[1][e];
    }
    if (enemy_smell[2][e] > max_value) {
        target_coords = icoords[2];
        max_value = enemy_smell[2][e];
    }
    if (enemy_smell[3][e] > max_value) {
        target_coords = icoords[3];
        max_value = enemy_smell[3][e];
    }    

    if (max_value > 0.0) {
        return package(target_coords, send);
    }


    // -- get cell with lowest own value
    target_coords = coords;
    min_value = 1.0;

    if (values[0][s] > 0.0 && values[0][s] < min_value && values[0][s] + send < 1.0) {
        target_coords = icoords[0];
        min_value = values[0][s];
    }
    if (values[1][s] > 0.0 && values[1][s] < min_value && values[1][s] + send < 1.0) {
        target_coords = icoords[1];
        min_value = values[1][s];
    }
    if (values[2][s] > 0.0 && values[2][s] < min_value && values[2][s] + send < 1.0) {
        target_coords = icoords[2];
        min_value = values[2][s];
    }
    if (values[3][s] > 0.0 && values[3][s] < min_value && values[3][s] + send < 1.0) {
        target_coords = icoords[3];
        min_value = values[3][s];
    } 

    if (min_value < 1.0) {
        return package(target_coords, send);
    } 

    // -- get cell with strongest enemy value
    target_coords = coords;
    max_value = 0.0;

    if (values[0][e] > max_value) {
        target_coords = icoords[0];
        max_value = values[0][e];
    }
    if (values[1][e] > max_value) {
        target_coords = icoords[1];
        max_value = values[1][e];
    }
    if (values[2][e] > max_value) {
        target_coords = icoords[2];
        max_value = values[2][e];
    }
    if (values[3][e] > max_value) {
        target_coords = icoords[3];
        max_value = values[3][e];
    } 

    if (max_value > 0.0) {
        return package(target_coords, send);
    }


    // ?? don't know what to do: do nothing

    return package(coords, 0.0);


    // return target_coords;
}

float strat_2 (int self, int enemy) {
    int s = self;
    int e = enemy;

    vec2 uv = TexCoord;
    vec2 coords = uv_to_coords(uv);
    vec4 cell = get_color_by_coord(PfGameTexture, coords);


    // growth multiplier
    float g = 1.01;

    // growth adder
    float a = 0.001;


    vec3 test;
    vec2 test_coords;
    float test_val;

    // send
    test = pick_target_2(self, enemy, coords);
    if (test.x != coords.x || test.y != coords.y) {
        cell[s] = cell[s] - test.z;
    }
  

    // receive
    test = pick_target_2(self, enemy, coords + Top);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z); // bug: can lead to loss of amount
    }
    test = pick_target_2(self, enemy, coords + Bottom);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z);
    }
    test = pick_target_2(self, enemy, coords + Left);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z);
    }
    test = pick_target_2(self, enemy, coords + Right);
    if (test.x == coords.x && test.y == coords.y) {
        cell[s] = min(1.0, cell[s] + test.z);
    }

    // return
    if (cell[s] > 0.0) {
        return max(0.0, min(1.0, cell[s] * g + a));
    }
    return 0.0;
}

float clamp_v(float a) {
    return min(1.0, max(0.0, a));
}

void main() {
    vec4 GameColor = texture2D(PfGameTexture, TexCoord);

    GameColor.r = strat_1(RED, GREEN);
    GameColor.g = strat_1(GREEN, RED); //strat2(GREEN);
    
    // blue lighting where battle takes place
    if (GameColor.r > 0.01 && GameColor.g > 0.01){
        GameColor.b = 10.0*abs(GameColor.g - GameColor.r);
    }
    GameColor.b -= 0.002;


    float buff_r = GameColor.r;
    float buff_g = GameColor.g;

    if (GameColor.r > 0.01 && GameColor.g > 0.01){
        GameColor.g = clamp_v(GameColor.g - (0.2 * buff_r));
        GameColor.r = clamp_v(GameColor.r - (0.2 * buff_g));
    }
    


    // // subtract red from green and vice versa to have them battle it out
    // if (GameColor.r > GameColor.g){
    //     // big win
    //     if (GameColor.r - GameColor.g > 0.5){
    //         // GameColor.r = GameColor.r * growth ;
    //         GameColor.r = (GameColor.r - GameColor.g) ;
    //         GameColor.g = 0;
    //     }
    //     else {
    //         GameColor.r = (GameColor.r - GameColor.g) ;
    //         GameColor.g = 0;
    //     }
    // } 
    // else {
    //     // big win
    //     if (GameColor.g - GameColor.r > 0.8){
    //         // GameColor.g = GameColor.g * growth ;
    //         GameColor.g = (GameColor.g - GameColor.r);
    //         GameColor.r = 0;
    //     }
    //     else {
    //         GameColor.g = (GameColor.g - GameColor.r);
    //         GameColor.r = 0;
    //     }
    // }

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