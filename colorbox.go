package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	"image/color"
)

func colorBox(c color.Color) *gtk.DrawingArea {
	r, g, b, _ := c.RGBA()
	colstring := fmt.Sprintf("#%02x%02x%02x", (r>>8)&0xff, (g>>8)&0xff, (b>>8)&0xff)

	dArea := gtk.NewDrawingArea()
	var pixmap *gdk.Pixmap
	var gc *gdk.GC

	dArea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		alloc := dArea.GetAllocation()
		pixmap = gdk.NewPixmap(dArea.GetWindow().GetDrawable(), alloc.Width, alloc.Height, 24)
		gc = gdk.NewGC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.NewColor(colstring))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
	})

	dArea.Connect("expose-event", func() {
		if pixmap != nil {
			dArea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	return dArea
}
