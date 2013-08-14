package main

import (
	"github.com/kch42/gomcmap/mcmap"
)

type BiomeGetter interface {
	GetBiome(x, z int) (mcmap.Biome, bool)
}

type XZPos struct {
	X, Z int
}

type Change map[XZPos]mcmap.Biome

type Tool interface {
	SingleClick() bool // Whether only one click should be performed (true) or the action should be repeated, if the mouse is dragged
	Do(bio mcmap.Biome, bioget BiomeGetter, x, z int) Change
}

type drawTool struct {
	radGetter func() int
}

func (d *drawTool) SingleClick() bool { return false }

func (d *drawTool) Do(bio mcmap.Biome, bioget BiomeGetter, x, z int) Change {
	rad := d.radGetter()
	if rad <= 0 {
		return nil
	}

	change := make(Change)
	for xp := x - (rad - 1); xp < x+rad; xp++ {
		for zp := z - (rad - 1); zp < z+rad; zp++ {
			change[XZPos{xp, zp}] = bio
		}
	}

	return change
}

func NewDrawTool(radGetter func() int) *drawTool {
	return &drawTool{radGetter}
}

type fillTool struct{}

func (f *fillTool) SingleClick() bool { return true }

func (f *fillTool) Do(bio mcmap.Biome, bioget BiomeGetter, x, z int) Change {
	oldbio, ok := bioget.GetBiome(x, z)
	if !ok {
		return nil
	}

	change := make(Change)
	floodfill(oldbio, bio, bioget, x, z, change)
	return change
}

func floodfill(oldbio, newbio mcmap.Biome, bioget BiomeGetter, x, z int, change Change) {
	pos := XZPos{x, z}
	if _, ok := change[pos]; ok {
		return
	}

	if bio, ok := bioget.GetBiome(x, z); ok && (bio == oldbio) {
		change[pos] = newbio

		floodfill(oldbio, newbio, bioget, x-1, z, change)
		floodfill(oldbio, newbio, bioget, x+1, z, change)
		floodfill(oldbio, newbio, bioget, x, z-1, change)
		floodfill(oldbio, newbio, bioget, x, z+1, change)
	}
}

func NewFillTool() *fillTool {
	return new(fillTool)
}
