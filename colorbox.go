package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

func colorBox(c *gdk.Color) *gtk.DrawingArea {
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
		gc.SetRgbFgColor(c)
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
	})

	dArea.Connect("expose-event", func() {
		if pixmap != nil {
			dArea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	return dArea
}
