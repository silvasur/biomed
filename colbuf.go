package main

import (
	"github.com/mattn/go-gtk/gdk"
)

type ColorBuffer map[string]*gdk.Color

func (cb ColorBuffer) Color(name string) *gdk.Color {
	if col, ok := cb[name]; ok {
		return col
	}

	col := gdk.NewColor(name)
	cb[name] = col
	return col
}
