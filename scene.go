package splish

import (
	"errors"
	"image"
	"image/draw"
)

type Scene struct {
	// Foreground is the highest level SubScene
	Foreground *SubScene
	// Background is the lowest level SubScene
	Background *SubScene

	window Window
	// regions should be a sorted by ascending priority
	regions []Region
	view    Rectangle
	Sprites []Sprite
}

func NewScene(w Window) Scene {
	s := Scene{window: w}

	return s
}
func (s *Scene) NewSprite(img image.Image) *Sprite {
	var rgba *image.RGBA
	if ri, ok := img.(*image.RGBA); ok {
		rgba = ri
	} else {
		rgba = image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)
	}

	tex := s.window.UploadTexture(rgba)
	spr := Sprite{Tex: tex, Img: rgba}
	spr.SubSprite(NewRectangle(img.Bounds()))

	s.Sprites = append(s.Sprites, spr)
	return &spr
}

func (s *Scene) Draw() {
	scenes := []*SubScene{}
	if s.Background != nil {
		scenes = append(scenes, s.Background)
	}
	for _, region := range s.regions {
		if region.Area.Overlaps(s.view) && region.Scene != nil {
			scenes = append(scenes, region.Scene)
		}
	}
	if s.Foreground != nil {
		scenes = append(scenes, s.Foreground)
	}

	s.window.Draw(
		scenes,
		s.view,
	)
}

func (s *Scene) TranslateView(x, y float32) {
	s.view.Min = s.view.Min.AddXY(x, y)
	s.view.Max = s.view.Max.AddXY(x, y)
}

func (s *Scene) ZoomView(z int) {
}

type SubScene struct {
	Global   bool
	Location Point
	Layers   []Layer2d
}

type Layer2d struct {
	Sprites   []SpriteInstance
	Transform image.Rectangle
}

type Point struct {
	X, Y float32
}

func NewPoint(ip image.Point) Point {
	return Point{float32(ip.X), float32(ip.Y)}
}

func (p Point) Add(s Point) Point {
	return Point{p.X + s.X, p.Y + s.Y}
}
func (p *Point) AddXY(x, y float32) Point {
	return Point{p.X + x, p.Y + y}
}
func (p Point) ToImage() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

type Rectangle struct {
	Min, Max Point
}

func NewRectangle(ir image.Rectangle) Rectangle {
	return Rectangle{NewPoint(ir.Min), NewPoint(ir.Max)}
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rectangle) Size() Point {
	return Point{r.Max.X - r.Min.X, r.Max.Y - r.Min.Y}
}

func (r Rectangle) ToImage() image.Rectangle {
	return image.Rectangle{r.Min.ToImage(), r.Max.ToImage()}
}

type SpriteInstance struct {
	Sprite   SpriteID
	Location Point
	Index    int
	Rotation int
	Scale    float32
}

type SpriteID int

type Sprite struct {
	Tex  Texture
	Img  *image.RGBA
	Size Point
	Sub  []interface{}
}

func (s *Sprite) SubSprite(rect Rectangle) int {
	s.Sub = append(s.Sub, s.Tex.MarkSubTexture(rect))
	return len(s.Sub) - 1
}

type Region struct {
	// Regions are drawn starting with low priority
	At       Point
	Priority int
	Scene    *SubScene
	Area     Rectangle
}

var todoErr = errors.New("TODO")
