package sdl2

import (
	"image"
	"log"
	"unsafe"

	"github.com/acsellers/splish"
	"github.com/veandco/go-sdl2/sdl"
)

func NewWindow(bounds splish.Rectangle, fullscreen bool) (*Window, error) {
	size := bounds.Size()
	flags := uint32(sdl.WINDOW_OPENGL | sdl.RENDERER_ACCELERATED)
	if fullscreen {
		flags |= sdl.WINDOW_FULLSCREEN
	}
	w, r, err := sdl.CreateWindowAndRenderer(int(size.X), int(size.Y), flags)
	if err != nil {
		return nil, err
	}
	return &Window{
		window:   w,
		renderer: r,
		w:        size.X,
		h:        size.Y,
	}, nil
}

type Window struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	scene    *splish.Scene
	w, h     float32
}

func (w *Window) Draw(ss []*splish.SubScene, r splish.Rectangle) {
	baseX, baseY := r.Min.X, r.Min.Y
	scale := r.Size().X / w.w
	for _, sub := range ss {
		sceneX, sceneY := scale*(sub.Location.X-baseX), scale*(sub.Location.Y-baseY)
		for _, layer := range sub.Layers {
			for _, sprite := range layer.Sprites {
				sd := w.scene.Sprites[sprite.Sprite]
				st := sd.Tex.(*Texture)
				sr := sdl.Rect{
					X: int32(sprite.Location.X*scale + sceneX),
					Y: int32(sprite.Location.Y*scale + sceneY),
					W: int32(sd.Size.X * scale),
					H: int32(sd.Size.Y * scale),
				}
				w.renderer.Copy(st.Tex, nil, &sr)
			}
		}
	}
	w.renderer.Present()
}

func (w *Window) UploadTexture(img *image.RGBA) splish.Texture {
	tex, err := w.NewTexture(img)
	if err != nil {
		log.Fatal(tex)
	}
	return tex
}

func (w *Window) NewTexture(img *image.RGBA) (*Texture, error) {
	tex := &Texture{Src: img}
	size := img.Bounds().Size()
	var err error
	tex.Tex, err = w.renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STATIC,
		size.X,
		size.Y,
	)
	if err == nil {
		err = tex.Tex.Update(nil, unsafe.Pointer(&img.Pix), img.Stride)
	}
	tex.Subs = append(tex.Subs, tex.Tex)

	return tex, err
}

type Texture struct {
	Win  *Window
	Tex  *sdl.Texture
	Subs []*sdl.Texture
	Src  *image.RGBA
}

func (t *Texture) Expire() {
	for _, s := range t.Subs {
		s.Destroy()
	}
}

func (t *Texture) MarkSubTexture(r splish.Rectangle) interface{} {
	si := t.Src.SubImage(r.ToImage()).(*image.RGBA)

	size := r.Size()
	var subTex *sdl.Texture
	var err error
	subTex, err = t.Win.renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STATIC,
		int(size.X),
		int(size.Y),
	)
	if err == nil {
		err = subTex.Update(nil, unsafe.Pointer(&si.Pix), t.Src.Stride)
	}
	if err != nil {
		return 0
	}
	t.Subs = append(t.Subs, subTex)

	return len(t.Subs) - 1
}
