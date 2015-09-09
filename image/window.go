package image

import (
	"fmt"
	img "image"
	"image/jpeg"
	"os"
	"time"

	"github.com/acsellers/splish"
)

func NewWindow(size img.Rectangle) splish.Window {
	return &Window{
		size: size,
	}
}

type Window struct {
	textures map[int]Texture
	index    int
	size     img.Rectangle
}

func (w *Window) Draw(scenes []*splish.SubScene, r img.Rectangle) {
	out := img.NewRGBA(r)

	f, err := os.Create(fmt.Sprintf("output-%d.jpg", time.Now().UnixNano()))
	if err != nil {
		fmt.Println("Image Output:", err)
		return
	}
	jpeg.Encode(f, out, nil)

}

func (w *Window) UploadTexture(im *img.RGBA) splish.Texture {
	t := Texture{w: w, img: im, id: w.index}
	w.textures[w.index] = t
	w.index++

	return &t
}

func (w *Window) OnSizeChange(chan splish.SizeEvent) {
}
func (w *Window) OnKey(chan splish.KeyEvent) {
}
func (w *Window) OnMouseKey(chan splish.MouseKeyEvent) {
}
func (w *Window) OnMouseMove(chan splish.MouseMoveEvent) {
}

type Texture struct {
	id   int
	w    *Window
	img  *img.RGBA
	subs []img.Rectangle
}

func (t *Texture) Expire() {
	delete(t.w.textures, t.id)
}

func (t *Texture) MarkSubTexture(r img.Rectangle) interface{} {
	t.subs = append(t.subs, r)
	return len(t.subs) - 1
}
