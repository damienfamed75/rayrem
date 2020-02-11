package physics

import r "github.com/lachee/raylib-goplus/raylib"

type Platform Rectangle

func NewPlatform(x, y, w, h float32) *Platform {
	p := &Platform{
		BasicShape: &BasicShape{},
		Rectangle:  r.NewRectangle(x, y, w, h),
	}

	return p
}
