package main

import (
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"math"
	"unsafe"
)

const (
	zoom          = 2
	tileSize      = zoom * mcmap.ChunkSizeXZ
	halfChunkSize = mcmap.ChunkSizeXZ / 2
)

type tileCmd int

const (
	cmdUpdateTiles tileCmd = iota
	cmdFlushTiles
	cmdSave
)

func renderTile(chunk *mcmap.Chunk) (maptile, biotile *gdk.Pixmap, biocache []mcmap.Biome) {
	maptile = emptyPixmap(tileSize, tileSize, 24)
	mtDrawable := maptile.GetDrawable()
	mtGC := gdk.NewGC(mtDrawable)

	biotile = emptyPixmap(tileSize, tileSize, 24)
	btDrawable := biotile.GetDrawable()
	btGC := gdk.NewGC(btDrawable)

	biocache = make([]mcmap.Biome, mcmap.ChunkRectXZ)

	i := 0
	for z := 0; z < mcmap.ChunkSizeXZ; z++ {
	scanX:
		for x := 0; x < mcmap.ChunkSizeXZ; x++ {
			bio := chunk.Biome(x, z)
			btGC.SetRgbFgColor(bioColors[bio])
			btDrawable.DrawRectangle(btGC, true, x*zoom, z*zoom, zoom, zoom)

			biocache[i] = bio
			i++

			for y := chunk.Height(x, z); y >= 0; y-- {
				if col, ok := blockColors[chunk.Block(x, y, z).ID]; ok {
					mtGC.SetRgbFgColor(col)
					mtDrawable.DrawRectangle(mtGC, true, x*zoom, z*zoom, zoom, zoom)
					continue scanX
				}
			}

			mtGC.SetRgbFgColor(gdk.NewColor("#ffffff"))
			mtDrawable.DrawRectangle(mtGC, true, x*zoom, z*zoom, zoom, zoom)
		}
	}

	return
}

type MapWidget struct {
	dArea *gtk.DrawingArea
	w, h  int

	reportFail func(msg string)
	updateInfo func(x, z int, bio mcmap.Biome)

	isInit bool

	showBiomes bool

	offX, offZ         int
	mx1, mx2, my1, my2 int

	pixmap   *gdk.Pixmap
	pixmapGC *gdk.GC
	gdkwin   *gdk.Window

	bg *gdk.Pixmap

	region *mcmap.Region

	maptiles map[XZPos]*gdk.Pixmap
	biotiles map[XZPos]*gdk.Pixmap
	bioCache map[XZPos][]mcmap.Biome

	redraw   chan bool
	tileCmds chan tileCmd

	tool Tool
}

var (
	checker1 = gdk.NewColor("#222222")
	checker2 = gdk.NewColor("#444444")
)

func emptyPixmap(w, h, depth int) *gdk.Pixmap {
	return gdk.NewPixmap(new(gdk.Drawable), w, h, depth)
}

func (mw *MapWidget) SetShowBiomes(b bool) {
	mw.showBiomes = b
	mw.redraw <- true
}

func (mw *MapWidget) SetTool(t Tool) {
	mw.tool = t
}

func (mw *MapWidget) SetBiome(bio mcmap.Biome) {
	mw.bio = bio
}

func (mw *MapWidget) DArea() *gtk.DrawingArea { return mw.dArea }

func (mw *MapWidget) doTileCmds() {
	for cmd := range mw.tileCmds {
		switch cmd {
		case cmdSave:
			mw.region.Save()
		case cmdFlushTiles:
			gdk.ThreadsEnter()
			for _, mt := range mw.maptiles {
				mt.Unref()
			}
			for _, bt := range mw.biotiles {
				bt.Unref()
			}
			gdk.ThreadsLeave()

			mw.maptiles = make(map[XZPos]*gdk.Pixmap)
			mw.biotiles = make(map[XZPos]*gdk.Pixmap)
			mw.bioCache = make(map[XZPos][]mcmap.Biome)
		case cmdUpdateTiles:
			todelete := make(map[XZPos]bool)

			startX := int(math.Floor(float64(mw.offX) / tileSize))
			startZ := int(math.Floor(float64(mw.offZ) / tileSize))
			endX := int(math.Ceil(float64(mw.offX+mw.w) / tileSize))
			endZ := int(math.Ceil(float64(mw.offZ+mw.h) / tileSize))

			for pos := range mw.maptiles {
				if (pos.X < startX) || (pos.Z < startZ) || (pos.X >= endX) || (pos.Z >= endZ) {
					todelete[pos] = true
				}
			}

			gdk.ThreadsEnter()
			for pos := range todelete {
				if tile, ok := mw.maptiles[pos]; ok {
					tile.Unref()
					delete(mw.maptiles, pos)
				}

				if tile, ok := mw.biotiles[pos]; ok {
					tile.Unref()
					delete(mw.biotiles, pos)
				}

				if _, ok := mw.bioCache[pos]; ok {
					delete(mw.bioCache, pos)
				}
			}

			if mw.region != nil {
				for z := startZ; z < endZ; z++ {
				scanX:
					for x := startX; x < endX; x++ {
						pos := XZPos{x, z}

						if _, ok := mw.biotiles[pos]; ok {
							continue scanX
						}

						chunk, err := mw.region.Chunk(x, z)
						switch err {
						case nil:
						case mcmap.NotAvailable:
							continue scanX
						default:
							mw.reportFail(fmt.Sprintf("Could not get chunk %d, %d: %s", x, z, err))
							return
						}

						mw.maptiles[pos], mw.biotiles[pos], mw.bioCache[pos] = renderTile(chunk)
						chunk.MarkUnused()

						gdk.ThreadsLeave()
						mw.redraw <- true
						gdk.ThreadsEnter()
					}
				}
			}

			gdk.ThreadsLeave()
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

	mw.pixmap = gdk.NewPixmap(mw.dArea.GetWindow().GetDrawable(), mw.w, mw.h, 24)
	mw.pixmapGC = gdk.NewGC(mw.pixmap.GetDrawable())

	mw.drawBg()
	gdk.ThreadsLeave()
	mw.redraw <- true
	gdk.ThreadsEnter()
}

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

func (mw *MapWidget) compose() {
	drawable := mw.pixmap.GetDrawable()
	gc := mw.pixmapGC

	drawable.DrawDrawable(gc, mw.bg.GetDrawable(), 0, 0, 0, 0, -1, -1)

	var tiles map[XZPos]*gdk.Pixmap
	if mw.showBiomes {
		tiles = mw.biotiles
	} else {
		tiles = mw.maptiles
	}

	for pos, tile := range tiles {
		x := (pos.X * tileSize) - mw.offX
		y := (pos.Z * tileSize) - mw.offZ

		drawable.DrawDrawable(gc, tile.GetDrawable(), 0, 0, x, y, tileSize, tileSize)
	}
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
	cx, cz, cbx, cbz := mcmap.BlockToChunk(x, z)
	bio := mcmap.Biome(mcmap.BioUncalculated)
	if bc, ok := mw.bioCache[XZPos{cx, cz}]; ok {
		bio = bc[cbz*mcmap.ChunkSizeXZ+cbx]
	}
	mw.updateInfo(x, z, bio)

	switch {
	case mt&gdk.BUTTON1_MASK != 0:
	case mt&gdk.BUTTON2_MASK != 0:
		if (mw.mx1 != -1) && (mw.my1 != -1) {
			mw.offX += mw.mx1 - mw.mx2
			mw.offZ += mw.my1 - mw.my2

			gdk.ThreadsLeave()
			mw.tileCmds <- cmdUpdateTiles
			gdk.ThreadsEnter()
		}
	}

	mw.mx1, mw.my1 = mw.mx2, mw.my2
}

func (mw *MapWidget) expose() {
	mw.dArea.GetWindow().GetDrawable().DrawDrawable(mw.pixmapGC, mw.pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
}

func (mw *MapWidget) guiUpdater() {
	for _ = range mw.redraw {
		gdk.ThreadsEnter()
		mw.compose()
		mw.expose()
		mw.dArea.GetWindow().Invalidate(nil, false)
		gdk.ThreadsLeave()
	}
}

func (mw *MapWidget) init() {
	mw.redraw = make(chan bool)
	mw.tileCmds = make(chan tileCmd)

	mw.maptiles = make(map[XZPos]*gdk.Pixmap)
	mw.biotiles = make(map[XZPos]*gdk.Pixmap)
	mw.bioCache = make(map[XZPos][]mcmap.Biome)

	mw.showBiomes = true

	mw.mx1, mw.my1 = -1, -1

	go mw.doTileCmds()
	go mw.guiUpdater()

	mw.dArea = gtk.NewDrawingArea()
	mw.dArea.Connect("configure-event", mw.configure)
	mw.dArea.Connect("expose-event", mw.expose)
	mw.dArea.Connect("motion-notify-event", mw.movement)

	mw.dArea.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.POINTER_MOTION_HINT_MASK | gdk.BUTTON_PRESS_MASK))
}

func (mw *MapWidget) setRegion(region *mcmap.Region) {
	mw.tileCmds <- cmdFlushTiles
	mw.region = region
	mw.tileCmds <- cmdUpdateTiles
}

func (mw *MapWidget) GetBiome(x, z int) (mcmap.Biome, bool) {
	cx, cz, bx, bz := mcmap.BlockToChunk(x, z)
	pos := XZPos{cx, cz}

	if bc, ok := mw.bioCache[pos]; ok {
		return bc[bz*mcmap.ChunkSizeXZ+bx], true
	}

	chunk, err := mw.region.Chunk(x, z)
	switch err {
	case nil:
	case mcmap.NotAvailable:
		return mcmap.BioUncalculated, false
	default:
		mw.reportFail(fmt.Sprintf("Error while getting chunk %d, %d: %s", cx, cz, err))
		return mcmap.BioUncalculated, false
	}
	defer chunk.MarkUnused()

	bc := make([]mcmap.Biome, mcmap.ChunkRectXZ)
	i := 0
	for z := 0; z < mcmap.ChunkSizeXZ; z++ {
		for x := 0; x < mcmap.ChunkSizeXZ; x++ {
			bc[i] = chunk.Biome(x, z)
			i++
		}
	}
	mw.bioCache[pos] = bc

	return chunk.Biome(bx, bz), true
}

func (mw *MapWidget) SetBiome(x, z int, bio mcmap.Biome) {
	cx, cz, bx, bz := mcmap.BlockToChunk(x, z)
	pos := XZPos{cx, cz}

	// Update cache
	if bc, ok := mw.bioCache[pos]; ok {
		bc[bz*mcmap.ChunkSizeXZ+bx] = bio
	}

	chunk, err := mw.region.Chunk(x, z)
	switch err {
	case nil:
	case mcmap.NotAvailable:
		return
	default:
		mw.reportFail(fmt.Sprintf("Error while getting chunk %d, %d: %s", cx, cz, err))
		return
	}
	defer chunk.MarkUnused()

	chunk.SetBiome(bx, bz, bio)
}

func NewMapWidget(reportFail func(msg string), updateInfo func(x, z int, bio mcmap.Biome)) *MapWidget {
	mw := &MapWidget{reportFail: reportFail, updateInfo: updateInfo}
	mw.init()
	return mw
}
