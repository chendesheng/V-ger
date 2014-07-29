package gui

import (
	"fmt"
	"log"

	"github.com/go-gl/gl"
)

//opengl shader helper

func shaderCompileFromFile(typ gl.GLenum, source string) (gl.Shader, error) {
	shader := gl.CreateShader(typ)
	shader.Source(source)
	shader.Compile()

	if shader.Get(gl.COMPILE_STATUS) == gl.FALSE {
		defer shader.Delete()

		content := shader.GetInfoLog()

		return 0, fmt.Errorf("Compile shader error: %s", content)
	}

	return shader, nil
}

func shaderAttachFromFile(program gl.Program, typ gl.GLenum, source string) {
	shader, err := shaderCompileFromFile(typ, source)
	if err == nil {
		program.AttachShader(shader)

		/* delete the shader - it won't actually be
		 * destroyed until the program that it's attached
		 * to has been destroyed */
		shader.Delete()
	} else {
		log.Fatal(err)
	}
}

const FSHADER = `uniform sampler2D tex_y;
uniform sampler2D tex_u;
uniform sampler2D tex_v;
varying vec2 v_texcoord;
void main() {
	//check: http://en.wikipedia.org/wiki/YCbCr

	float y = texture2D(tex_y, v_texcoord).r-16.0/255.0;
	float u = texture2D(tex_u, v_texcoord).r-128.0/255.0;
	float v = texture2D(tex_v, v_texcoord).r-128.0/255.0;

	gl_FragColor = vec4(y, u, v, 1.0) * mat4(
		1.16438356164384,  0			   ,  1.59602678571428, 0,
		1.16438356164384, -0.39176229009491, -0.81296764723777, 0,
		1.16438356164384,  2.01723214285714,  0               , 0,
		0               ,  0               ,  0               , 1.0
	);
}`
const VSHADER = `varying vec2 v_texcoord;
void main() {
    gl_Position = ftransform();
    v_texcoord = vec2(gl_MultiTexCoord0.s, 1.0-gl_MultiTexCoord0.t);
}`

type yuvRender struct {
	texYUV  [3]gl.Texture
	program gl.Program
}

func (r *yuvRender) delete() {
	r.program.Delete()

	for _, tex := range r.texYUV {
		tex.Delete()
	}
}

func (r *yuvRender) draw(img []byte, width, height int) {
	if len(img) != width*height*3/2 {
		return
	}

	channels := getYUVChannels(img, width, height)

	for i, tex := range r.texYUV {
		c := channels[i]

		uniform := r.program.GetUniformLocation(c.name)
		gl.ActiveTexture(gl.TEXTURE0 + gl.GLenum(i))
		tex.Bind(gl.TEXTURE_2D)
		uniform.Uniform1i(i)

		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, c.width, c.height, gl.LUMINANCE,
			gl.UNSIGNED_BYTE, c.data)
	}
}

type channel struct {
	data []byte

	width  int
	height int
	name   string
}

func getYUVChannels(img []byte, width, height int) [3]channel {
	var channels [3]channel
	channels[0].data = img[:width*height]
	channels[0].width = width
	channels[0].height = height
	channels[0].name = "tex_y"

	channels[1].data = img[width*height : width*height*5/4]
	channels[1].width = width / 2
	channels[1].height = height / 2
	channels[1].name = "tex_u"

	channels[2].data = img[width*height*5/4:]
	channels[2].width = width / 2
	channels[2].height = height / 2
	channels[2].name = "tex_v"

	return channels
}

func NewYUVRender(width, height int) imageRender {
	size := width * height * 3 / 2
	img := make([]byte, size)
	// if cap(img) > size {
	// 	img = img[:size]
	// } else {

	// }
	log.Print("NewYUVRender:", len(img), size, width*height)

	r := yuvRender{}

	channels := getYUVChannels(img, width, height)

	for i, _ := range r.texYUV {
		tex := gl.GenTexture()
		tex.Bind(gl.TEXTURE_2D)

		c := channels[i]

		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.LUMINANCE, c.width, c.height, 0,
			gl.LUMINANCE, gl.UNSIGNED_BYTE, c.data)

		r.texYUV[i] = tex
	}

	r.program = gl.CreateProgram()

	shaderAttachFromFile(r.program, gl.FRAGMENT_SHADER, FSHADER)
	shaderAttachFromFile(r.program, gl.VERTEX_SHADER, VSHADER)
	r.program.Link()
	r.program.Use()

	return &r
}
