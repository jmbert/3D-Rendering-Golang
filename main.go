package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

var darkFactor uint8
var width int
var height int
var step float64
var camRot Matrix

func main() {
	var scene = Scene{}
	var cam = Camera{Vector3{5, 6, -6}}
	var sphere = Sphere{Vector3{0, 1, 0}, 1, 1, 2, color.RGBA{0, 0, 255, 255}}
	var cube = Cube{Vector3{2, 0, 2}, color.RGBA{0, 255, 0, 255}, 2, 2, 5}
	var lamp = Light{Vector3{0, 0, -3}}

	camRot = cam.RotateX(-math.Pi / 6)

	darkFactor = 3
	step = 0.01
	width = 200
	height = 200

	scene.objects = append(scene.objects, cube)
	scene.objects = append(scene.objects, sphere)

	scene.light = lamp

	img := cam.Render(scene, color.RGBA{0, 255, 255, 255})

	img = imaging.Rotate180(img)

	f, _ := os.Create("img.png")
	png.Encode(f, img)
}

type Object interface {
	CheckIfInside(point Vector3) bool
	GetSurfaceColour() color.RGBA
}

type Cube struct {
	position         Vector3
	surfaceColour    color.RGBA
	xLen, yLen, zLen float64
}

func (c Cube) CheckIfInside(point Vector3) bool {
	if c.position.x < point.x && point.x < c.position.x+c.xLen && c.position.y < point.y && point.y < c.position.y+c.yLen && c.position.z < point.z && point.z < c.position.z+c.zLen {
		return true
	} else {
		return false
	}
}

func (c Cube) GetSurfaceColour() color.RGBA {
	return c.surfaceColour
}

type Sphere struct {
	centre        Vector3
	a, b, c       float64
	surfaceColour color.RGBA
}

func (s Sphere) CheckIfInside(p Vector3) bool {
	if math.Pow((p.x-s.centre.x)/s.a, 2)+math.Pow((p.y-s.centre.y)/s.b, 2)+math.Pow((p.z-s.centre.z)/s.c, 2) < 1 {
		return true
	} else {
		return false
	}
}

func (s Sphere) GetSurfaceColour() color.RGBA {
	return s.surfaceColour
}

type Camera struct {
	position Vector3
}

func (c Camera) RotateY(theta float64) Matrix {
	var m Matrix
	m.r1.x = 1
	m.r1.y = 0
	m.r1.z = 0
	m.r2.x = 0
	m.r2.y = math.Cos(theta)
	m.r2.z = -math.Sin(theta)
	m.r3.x = 0
	m.r3.y = math.Sin(theta)
	m.r3.z = math.Cos(theta)
	return m
}

func (c Camera) RotateX(theta float64) Matrix {
	var m Matrix
	m.r1.x = math.Cos(theta)
	m.r1.y = 0
	m.r1.z = math.Sin(theta)
	m.r2.x = 0
	m.r2.y = 1
	m.r2.z = 0
	m.r3.x = -math.Sin(theta)
	m.r3.y = 0
	m.r3.z = math.Cos(theta)
	return m
}

func (c Camera) Rotate(thetaX, thetaY float64) Matrix {
	var m Matrix

	mx := c.RotateX(thetaX)
	my := c.RotateY(thetaY)

	m = mx.Multiply(my)

	return m
}

func (m Matrix) Multiply(m2 Matrix) Matrix {
	var m3 Matrix
	m3.r1.x = m.r1.Dot(Vector3{m2.r1.x, m2.r2.x, m2.r3.x})
	m3.r2.x = m.r2.Dot(Vector3{m2.r1.x, m2.r2.x, m2.r3.x})
	m3.r3.x = m.r3.Dot(Vector3{m2.r1.x, m2.r2.x, m2.r3.x})
	m3.r1.x = m.r1.Dot(Vector3{m2.r1.y, m2.r2.y, m2.r3.y})
	m3.r2.x = m.r2.Dot(Vector3{m2.r1.y, m2.r2.y, m2.r3.y})
	m3.r3.x = m.r3.Dot(Vector3{m2.r1.y, m2.r2.y, m2.r3.y})
	m3.r1.x = m.r1.Dot(Vector3{m2.r1.z, m2.r2.z, m2.r3.z})
	m3.r2.x = m.r2.Dot(Vector3{m2.r1.z, m2.r2.z, m2.r3.z})
	m3.r3.x = m.r3.Dot(Vector3{m2.r1.z, m2.r2.z, m2.r3.z})
	return m3
}

type Matrix struct {
	r1, r2, r3 Vector3
}

func (m Matrix) Transform(v Vector3) Vector3 {
	var transformed Vector3

	transformed.x = v.Dot(m.r1)
	transformed.y = v.Dot(m.r2)
	transformed.z = v.Dot(m.r3)

	return transformed
}

type Vector3 struct {
	x, y, z float64
}

func (vec Vector3) Dot(vec2 Vector3) float64 {
	var dot = (vec.x * vec2.x) + (vec.y * vec2.y) + (vec.z * vec2.z)
	return dot
}

func (vec Vector3) Print() string {
	return strconv.FormatFloat(vec.x, 'E', 2, 64) + ", " + strconv.FormatFloat(vec.y, 'E', 2, 64) + ", " + strconv.FormatFloat(vec.z, 'E', 2, 64)
}

func (vec Vector3) Add(vec2 Vector3) Vector3 {
	return Vector3{vec.x + vec2.x, vec.y + vec2.y, vec.z + vec2.z}
}

func (vec Vector3) Scale(scalar float64) Vector3 {
	return Vector3{vec.x * scalar, vec.y * scalar, vec.z * scalar}
}

func (vec Vector3) Normalise() Vector3 {
	var normal = vec.Scale(1 / vec.Magnitude())
	return normal
}

func (vec Vector3) Magnitude() float64 {
	var mag = math.Sqrt(math.Pow(vec.x, 2) + math.Pow(vec.y, 2) + math.Pow(vec.z, 2))
	return mag
}

type Scene struct {
	objects []Object
	light   Light
}

func (cam Camera) Render(scene Scene, skybox color.RGBA) image.Image {
	var img = image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0.0; y < float64(height); y++ {
		for x := 0.0; x < float64(width); x++ {
			hit, objHit, hitPos := Raycast(cam.position, camRot.Transform(Vector3{x/(float64(width)/2) - 1, y/(float64(height)/2) - 1, 1}), 50, scene, step, step)
			if hit {
				lHit, _, _ := Raycast(hitPos, scene.light.position.Add(hitPos.Scale(-1)), 50, scene, step, step*10)
				if lHit {
					var darkSurface = color.RGBA{objHit.GetSurfaceColour().R / darkFactor, objHit.GetSurfaceColour().G / darkFactor, objHit.GetSurfaceColour().B / darkFactor, objHit.GetSurfaceColour().A}
					img.Set(int(x), int(y), darkSurface)
					println("dark")
				} else {
					img.Set(int(x), int(y), objHit.GetSurfaceColour())
					println("light")
				}
			} else {
				img.Set(int(x), int(y), skybox)
			}
		}
	}

	return img
}

func Raycast(origin Vector3, direction Vector3, distance float64, scene Scene, step float64, initalstep float64) (bool, Object, Vector3) {
	var rayPosition = origin
	var tStep float64
	for x := 0.0; x < distance; x += step {
		if x == 0.0 {
			tStep = initalstep
		} else {
			tStep = step
		}
		rayPosition = rayPosition.Add(direction.Normalise().Scale(tStep))
		for i := 0; i < len(scene.objects); i++ {
			if scene.objects[i].CheckIfInside(rayPosition) {
				return true, scene.objects[i], rayPosition
			}
		}
	}
	return false, nil, Vector3{0, 0, 0}
}

type Light struct {
	position Vector3
}
