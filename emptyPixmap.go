package main

import (
	"github.com/mattn/go-gtk/gdk"
)

// emptyPixmap creates an empty pixmap.
func emptyPixmap(w, h, depth int) *gdk.Pixmap {
	// The underlying C function would create an empty, unbound pixmap, if the drawable paramater was a NULL pointer.
	// Since simply passing a nil value would result in a panic (dereferencing a nil pointer), we use a new gdk.Drawable.
	// gdk.Drawable contains a C pointer which is NULL by default. So passing a new(gdk.Drawable) does the trick.
	return gdk.NewPixmap(new(gdk.Drawable), w, h, depth)
}
