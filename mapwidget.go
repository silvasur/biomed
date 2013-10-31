package main

import (
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"math"
	"unsafe"
)

const (
	zoom     = 2
	tileSize = zoom * mcmap.ChunkSizeXZ
)

type GUICallbacks struct {
	reportFail func(msg string)
	updateInfo func(x, z int, bio mcmap.Biome, name string)
	setBusy    func(bool)
}

type MapWidget struct {
	dArea *gtk.DrawingArea
	w, h  int

	guicbs GUICallbacks

	isInit bool

	showBiomes bool

	offX, offZ            int
	mx1, mx2, my1, my2    int
	continueTool, panning bool

	pixmap   *gdk.Pixmap
	pixmapGC *gdk.GC
	gdkwin   *gdk.Window

	bg *gdk.Pixmap

	regWrap *RegionWrapper

	bioLookup BiomeLookup
}

func (mw *MapWidget) calcChunkRect() {

}

func (mw *MapWidget) DArea() *gtk.DrawingArea { return mw.dArea }

func (mw *MapWidget) SetShowBiomes(b bool) {
	mw.showBiomes = b
	mw.updateGUI()
}

func (mw *MapWidget) SetFixSnowIce(b bool)           { mw.regWrap.SetFixSnowIce(b) }
func (mw *MapWidget) SetBiome(bio mcmap.Biome)       { mw.regWrap.SetBiome(bio) }
func (mw *MapWidget) SetRegion(region *mcmap.Region) { mw.regWrap.SetRegion(region) }
func (mw *MapWidget) SetTool(t Tool)                 { mw.regWrap.SetTool(t) }

func (mw *MapWidget) Save() { mw.regWrap.Save() }

func (mw *MapWidget) updateChunkBounds() {
	startX := int(math.Floor(float64(mw.offX) / tileSize))
	startZ := int(math.Floor(float64(mw.offZ) / tileSize))
	endX := int(math.Ceil(float64(mw.offX+mw.w) / tileSize))
	endZ := int(math.Ceil(float64(mw.offZ+mw.h) / tileSize))
	mw.regWrap.SetChunkBounds(startX, startZ, endX, endZ)
}

func (mw *MapWidget) movement(ctx *glib.CallbackContext) {
	if mw.gdkwin == nil {
		mw.gdkwin = mw.dArea.GetWindow()
	}
	arg := ctx.Args(0)
	mev := *(**gdk.EventMotion)(unsafe.Pointer(&arg))
	var mt gdk.ModifierType
	if mev.IsHint != 0 {
		mw.gdkwin.GetPointer(&(mw.mx2), &(mw.my2), &mt)
	} else {
		mw.mx2, mw.my2 = int(mev.X), int(mev.Y)
	}

	x := (mw.offX + mw.mx2) / zoom
	z := (mw.offZ + mw.my2) / zoom

	bio := mcmap.Biome(mcmap.BioUncalculated)
	if _bio, ok := mw.regWrap.GetBiomeAt(x, z); ok {
		bio = _bio
	}
	mw.guicbs.updateInfo(x, z, bio, mw.bioLookup.Name(bio))

	if mw.panning {
		if (mw.mx1 != -1) && (mw.my1 != -1) {
			mw.offX += mw.mx1 - mw.mx2
			mw.offZ += mw.my1 - mw.my2

			mw.updateGUI()
		}
	}

	if mw.continueTool {
		mw.regWrap.UseTool(x, z)

		mw.updateGUI()
	}

	mw.mx1, mw.my1 = mw.mx2, mw.my2
}

func (mw *MapWidget) buttonChanged(ctx *glib.CallbackContext) {
	arg := ctx.Args(0)
	bev := *(**gdk.EventButton)(unsafe.Pointer(&arg))

	switch gdk.EventType(bev.Type) {
	case gdk.BUTTON_RELEASE:
		if mw.panning {
			mw.panning = false

			mw.updateChunkBounds()

			gdk.ThreadsLeave()
			mw.regWrap.UpdateTiles()
			gdk.ThreadsEnter()
		}

		mw.continueTool = false
	case gdk.BUTTON_PRESS:
		switch bev.Button {
		case 1:
			if !mw.regWrap.RegionLoaded() {
				return
			}
			x := (mw.offX + int(bev.X)) / zoom
			z := (mw.offZ + int(bev.Y)) / zoom
			mw.regWrap.UseTool(x, z)

			mw.updateGUI()

			if !mw.regWrap.ToolSingleClick() {
				mw.continueTool = true
			}
		case 2:
			mw.panning = true
		}
	}
}

var (
	checker1 = gdk.NewColor("#222222")
	checker2 = gdk.NewColor("#444444")
)

func (mw *MapWidget) drawBg() {
	if mw.bg != nil {
		mw.bg.Unref()
	}

	mw.bg = emptyPixmap(mw.w, mw.h, 24)
	drawable := mw.bg.GetDrawable()
	gc := gdk.NewGC(drawable)

	w := int(math.Ceil(float64(mw.w) / 32))
	h := int(math.Ceil(float64(mw.h) / 32))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (x % 2) == (y % 2) {
				gc.SetRgbFgColor(checker1)
			} else {
				gc.SetRgbFgColor(checker2)
			}
			drawable.DrawRectangle(gc, true, x*32, y*32, 32, 32)
		}
	}
}

func (mw *MapWidget) configure() {
	if mw.pixmap != nil {
		mw.pixmap.Unref()
	}

	alloc := mw.dArea.GetAllocation()
	mw.w = alloc.Width
	mw.h = alloc.Height

	if !mw.isInit {
		mw.offX = -(mw.w / 2)
		mw.offZ = -(mw.h / 2)
		mw.isInit = true
	}

	mw.updateChunkBounds()

	mw.pixmap = gdk.NewPixmap(mw.dArea.GetWindow().GetDrawable(), mw.w, mw.h, 24)
	mw.pixmapGC = gdk.NewGC(mw.pixmap.GetDrawable())

	mw.drawBg()
	mw.updateGUI()
}

func (mw *MapWidget) compose() {
	drawable := mw.pixmap.GetDrawable()
	gc := mw.pixmapGC

	drawable.DrawDrawable(gc, mw.bg.GetDrawable(), 0, 0, 0, 0, -1, -1)

	var tiles map[XZPos]*gdk.Pixmap
	if mw.showBiomes {
		tiles = mw.regWrap.Biotiles
	} else {
		tiles = mw.regWrap.Maptiles
	}

	for pos, tile := range tiles {
		x := (pos.X * tileSize) - mw.offX
		y := (pos.Z * tileSize) - mw.offZ

		drawable.DrawDrawable(gc, tile.GetDrawable(), 0, 0, x, y, tileSize, tileSize)
	}
}

func (mw *MapWidget) expose() {
	mw.dArea.GetWindow().GetDrawable().DrawDrawable(mw.pixmapGC, mw.pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
}

func (mw *MapWidget) updateGUI() {
	mw.compose()
	mw.expose()
	mw.dArea.GetWindow().Invalidate(nil, false)
}

func NewMapWidget(guicbs GUICallbacks, bioLookup BiomeLookup) *MapWidget {
	dArea := gtk.NewDrawingArea()

	mw := &MapWidget{
		dArea:      dArea,
		guicbs:     guicbs,
		showBiomes: true,
		mx1:        -1,
		my1:        -1,
		bioLookup:  bioLookup,
	}

	mw.regWrap = NewRegionWrapper(mw.updateGUI, guicbs)
	mw.regWrap.bioLookup = bioLookup

	dArea.Connect("configure-event", mw.configure)
	dArea.Connect("expose-event", mw.expose)
	dArea.Connect("motion-notify-event", mw.movement)
	dArea.Connect("button-press-event", mw.buttonChanged)
	dArea.Connect("button-release-event", mw.buttonChanged)
	dArea.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.POINTER_MOTION_HINT_MASK | gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK))

	return mw
}
