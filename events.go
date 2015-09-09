package splish

import "image"

type SizeEvent image.Rectangle

type KeyEvent struct {
	IsDown bool
	Key    interface{}
}

type MouseKeyEvent struct {
}

type MouseMoveEvent struct {
}
