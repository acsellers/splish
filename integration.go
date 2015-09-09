package splish

import "image"

type Window interface {
	Draw([]*SubScene, Rectangle)

	UploadTexture(*image.RGBA) Texture
}

type Texture interface {
	Expire()
	MarkSubTexture(Rectangle) interface{}
}
