varying vec2 v_texcoord;
void main() {
    gl_Position = ftransform();
    v_texcoord = gl_MultiTexCoord0.st;
}