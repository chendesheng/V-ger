uniform sampler2D tex_y;
uniform sampler2D tex_u;
uniform sampler2D tex_v;
varying vec2 v_texcoord;
void main() {
	//check: http://en.wikipedia.org/wiki/YCbCr

	float y = texture2D(tex_y, v_texcoord).r*255.0-16.0;
	float u = texture2D(tex_u, v_texcoord).r*255.0-128.0;
	float v = texture2D(tex_v, v_texcoord).r*255.0-128.0;
	
	float r = 255.0/219.0*y + 255.0/112.0*0.701*v;
	float g = 255.0/219.0*y - 255.0/112.0*0.886*0.114/0.587*u - 255.0/112.0*0.701*0.299/0.587*v;
	float b = 255.0/219.0*y + 255.0/112.0*0.886*u;

	vec3 rgb = vec3(r, g, b)/255.0;

    gl_FragColor = vec4(rgb, 1.0);
}