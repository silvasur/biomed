package main

import (
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
)

const cacheSize = 4

type RegionWrapper struct {
	region *CachedRegion

	tileUpdates chan bool

	Maptiles map[XZPos]*gdk.Pixmap
	Biotiles map[XZPos]*gdk.Pixmap
	bioCache map[XZPos][]mcmap.Biome

	redraw func()
	guicbs GUICallbacks

	toolsEnabled bool
	tool         Tool

	fixSnowIce bool

	bio mcmap.Biome

	startX, startZ, endX, endZ int
}

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

func (rw *RegionWrapper) tileUpdater() {
	for _ = range rw.tileUpdates {
		todelete := make(map[XZPos]bool)

		for pos := range rw.Maptiles {
			if (pos.X < rw.startX) || (pos.Z < rw.startZ) || (pos.X >= rw.endX) || (pos.Z >= rw.endZ) {
				todelete[pos] = true
			}
		}

		gdk.ThreadsEnter()
		for pos := range todelete {
			if tile, ok := rw.Maptiles[pos]; ok {
				tile.Unref()
				delete(rw.Maptiles, pos)
			}

			if tile, ok := rw.Biotiles[pos]; ok {
				tile.Unref()
				delete(rw.Biotiles, pos)
			}

			if _, ok := rw.bioCache[pos]; ok {
				delete(rw.bioCache, pos)
			}
		}

		if rw.region != nil {
			for z := rw.startZ; z < rw.endZ; z++ {
			scanX:
				for x := rw.startX; x < rw.endX; x++ {
					pos := XZPos{x, z}

					if _, ok := rw.Biotiles[pos]; ok {
						continue scanX
					}

					chunk, err := rw.region.Chunk(x, z)
					switch err {
					case nil:
					case mcmap.NotAvailable:
						continue scanX
					default:
						rw.guicbs.reportFail(fmt.Sprintf("Could not get chunk %d, %d: %s", x, z, err))
						return
					}

					rw.Maptiles[pos], rw.Biotiles[pos], rw.bioCache[pos] = renderTile(chunk)
					chunk.MarkUnused()

					rw.redraw()
				}
			}
		}

		gdk.ThreadsLeave()
	}
}

func (rw *RegionWrapper) SetRegion(region *mcmap.Region) {
	if rw.RegionLoaded() {
		rw.flushTiles()
	}
	rw.region = NewCachedRegion(region, cacheSize)

	rw.tileUpdates <- true
}

func (rw *RegionWrapper) SetChunkBounds(startX, startZ, endX, endZ int) {
	rw.startX = startX
	rw.startZ = startZ
	rw.endX = endX
	rw.endZ = endZ
}

func (rw *RegionWrapper) SetTool(t Tool)           { rw.tool = t }
func (rw *RegionWrapper) SetBiome(bio mcmap.Biome) { rw.bio = bio }
func (rw *RegionWrapper) SetFixSnowIce(b bool)     { rw.fixSnowIce = b }

func (rw *RegionWrapper) RegionLoaded() bool    { return rw.region != nil }
func (rw *RegionWrapper) ToolSingleClick() bool { return rw.tool.SingleClick() }

func (rw *RegionWrapper) flushTiles() {
	if err := rw.region.Flush(); err != nil {
		rw.guicbs.reportFail(fmt.Sprintf("Error while flushing cache: %s", err))
		return
	}

	for _, mt := range rw.Maptiles {
		mt.Unref()
	}
	for _, bt := range rw.Biotiles {
		bt.Unref()
	}

	rw.Maptiles = make(map[XZPos]*gdk.Pixmap)
	rw.Biotiles = make(map[XZPos]*gdk.Pixmap)
	rw.bioCache = make(map[XZPos][]mcmap.Biome)
}

func (rw *RegionWrapper) Save() {
	rw.flushTiles()

	if err := rw.region.Flush(); err != nil {
		rw.guicbs.reportFail(fmt.Sprintf("Error while flushing cache: %s", err))
		return
	}
	if err := rw.region.Region.Save(); err != nil {
		rw.guicbs.reportFail(fmt.Sprintf("Error while saving: %s", err))
		return
	}
}

func (rw *RegionWrapper) UseTool(x, z int) {
	gdk.ThreadsLeave()
	defer gdk.ThreadsEnter()

	if !rw.toolsEnabled {
		return
	}

	if rw.tool.IsSlow() {
		rw.toolsEnabled = false
		rw.guicbs.setBusy(true)

		go func() {
			rw.tool.Do(rw.bio, rw, x, z)
			rw.guicbs.setBusy(false)
			rw.toolsEnabled = true

			gdk.ThreadsEnter()
			rw.redraw()
			gdk.ThreadsLeave()
		}()
	} else {
		rw.tool.Do(rw.bio, rw, x, z)
	}
}

func (rw *RegionWrapper) GetBiomeAt(x, z int) (mcmap.Biome, bool) {
	cx, cz, bx, bz := mcmap.BlockToChunk(x, z)
	pos := XZPos{cx, cz}

	if bc, ok := rw.bioCache[pos]; ok {
		return bc[bz*mcmap.ChunkSizeXZ+bx], true
	}

	if !rw.RegionLoaded() {
		return mcmap.BioUncalculated, false
	}

	chunk, err := rw.region.Chunk(cx, cz)
	switch err {
	case nil:
	case mcmap.NotAvailable:
		return mcmap.BioUncalculated, false
	default:
		rw.guicbs.reportFail(fmt.Sprintf("Error while getting chunk %d, %d: %s", cx, cz, err))
		return mcmap.BioUncalculated, false
	}

	bc := make([]mcmap.Biome, mcmap.ChunkRectXZ)
	i := 0
	for z := 0; z < mcmap.ChunkSizeXZ; z++ {
		for x := 0; x < mcmap.ChunkSizeXZ; x++ {
			bc[i] = chunk.Biome(x, z)
			i++
		}
	}
	rw.bioCache[pos] = bc

	return chunk.Biome(bx, bz), true
}

func fixFreeze(bx, bz int, chunk *mcmap.Chunk) (newcol *gdk.Color) {
	for y := mcmap.ChunkSizeY - 1; y >= 0; y-- {
		if blk := chunk.Block(bx, y, bz); blk.ID != mcmap.BlkAir {
			if (blk.ID == mcmap.BlkStationaryWater) || (blk.ID == mcmap.BlkWater) {
				blk.ID = mcmap.BlkIce
				newcol = blockColors[mcmap.BlkIce]
			} else if blockCanSnowIn[blk.ID] {
				if yFix := y + 1; yFix < mcmap.ChunkSizeY {
					blkFix := chunk.Block(bx, yFix, bz)
					blkFix.ID = mcmap.BlkSnow
					blkFix.Data = 0x0
					newcol = blockColors[mcmap.BlkSnow]
				}
			}

			break
		}
	}

	return
}

func fixMelt(bx, bz int, chunk *mcmap.Chunk) (newcol *gdk.Color) {
	for y := mcmap.ChunkSizeY - 1; y >= 0; y-- {
		if blk := chunk.Block(bx, y, bz); blk.ID != mcmap.BlkAir {
			if blk.ID == mcmap.BlkIce {
				blk.ID = mcmap.BlkStationaryWater
				blk.Data = 0x0
				newcol = blockColors[mcmap.BlkStationaryWater]
			} else if blk.ID == mcmap.BlkSnow {
				blk.ID = mcmap.BlkAir
				for y2 := y - 1; y2 >= 0; y2-- {
					if col, ok := blockColors[chunk.Block(bx, y2, bz).ID]; ok {
						newcol = col
						break
					}
				}
			}

			break
		}
	}

	return
}

func (rw *RegionWrapper) SetBiomeAt(x, z int, bio mcmap.Biome) {
	cx, cz, bx, bz := mcmap.BlockToChunk(x, z)
	pos := XZPos{cx, cz}

	chunk, err := rw.region.Chunk(cx, cz)
	switch err {
	case nil:
	case mcmap.NotAvailable:
		return
	default:
		rw.guicbs.reportFail(fmt.Sprintf("Error while getting chunk %d, %d: %s", cx, cz, err))
		return
	}

	chunk.SetBiome(bx, bz, bio)

	var newcol *gdk.Color
	if rw.fixSnowIce {
		if coldBiome[bio] {
			newcol = fixFreeze(bx, bz, chunk)
		} else {
			newcol = fixMelt(bx, bz, chunk)
		}
	}

	chunk.MarkModified()

	// Update cache
	if bc, ok := rw.bioCache[pos]; ok {
		bc[bz*mcmap.ChunkSizeXZ+bx] = bio
	}

	// Update tile
	if biotile, ok := rw.Biotiles[pos]; ok {
		gdk.ThreadsEnter()

		drawable := biotile.GetDrawable()
		gc := gdk.NewGC(drawable)
		gc.SetRgbFgColor(bioColors[bio])
		drawable.DrawRectangle(gc, true, bx*zoom, bz*zoom, zoom, zoom)

		if newcol != nil {
			drawable = rw.Maptiles[pos].GetDrawable()
			gc = gdk.NewGC(drawable)
			gc.SetRgbFgColor(newcol)
			drawable.DrawRectangle(gc, true, bx*zoom, bz*zoom, zoom, zoom)
		}

		gdk.ThreadsLeave()
	}
}

func (rw *RegionWrapper) UpdateTiles() {
	rw.tileUpdates <- true
}

func NewRegionWrapper(redraw func(), guicbs GUICallbacks) *RegionWrapper {
	rw := &RegionWrapper{
		redraw:       redraw,
		tileUpdates:  make(chan bool),
		toolsEnabled: true,
		guicbs:       guicbs,
		Maptiles:     make(map[XZPos]*gdk.Pixmap),
		Biotiles:     make(map[XZPos]*gdk.Pixmap),
		bioCache:     make(map[XZPos][]mcmap.Biome),
	}
	go rw.tileUpdater()
	return rw
}
