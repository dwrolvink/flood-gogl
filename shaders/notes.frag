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

bool equals(vec2 pos1, vec2 pos2) {
    if (abs(abs(pos1.x) - abs(pos2.x)) > 400.){
        return false;
    }
    if (abs(abs(pos1.y) - abs(pos2.y)) > 4000.){
        return false;
    }
    return true;
}