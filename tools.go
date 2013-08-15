package main

import (
	"github.com/kch42/gomcmap/mcmap"
)

type BiomeGetSetter interface {
	GetBiome(x, z int) (mcmap.Biome, bool)
	SetBiome(x, z int, bio mcmap.Biome)
}

type XZPos struct {
	X, Z int
}

type Tool interface {
	SingleClick() bool // Whether only one click should be performed (true) or the action should be repeated, if the mouse is dragged
	Do(bio mcmap.Biome, biogs BiomeGetSetter, x, z int)
}

type drawTool struct {
	radGetter func() int
}

func (d *drawTool) SingleClick() bool { return false }

func (d *drawTool) Do(bio mcmap.Biome, biogs BiomeGetSetter, x, z int) {
	rad := d.radGetter()
	if rad <= 0 {
		return
	}

	for xp := x - (rad - 1); xp < x+rad; xp++ {
		for zp := z - (rad - 1); zp < z+rad; zp++ {
			biogs.SetBiome(xp, zp, bio)
		}
	}
}

func NewDrawTool(radGetter func() int) *drawTool {
	return &drawTool{radGetter}
}

type fillTool struct{}

func (f *fillTool) SingleClick() bool { return true }

func (f *fillTool) Do(bio mcmap.Biome, biogs BiomeGetSetter, x, z int) {
	if oldbio, ok := biogs.GetBiome(x, z); ok {
		floodfill(oldbio, bio, biogs, x, z)
	}
}

func floodfill(oldbio, newbio mcmap.Biome, biogs BiomeGetSetter, x, z int) {
	if bio, ok := biogs.GetBiome(x, z); ok && (bio == oldbio) {
		biogs.SetBiome(x, z, newbio)

		floodfill(oldbio, newbio, biogs, x-1, z)
		floodfill(oldbio, newbio, biogs, x+1, z)
		floodfill(oldbio, newbio, biogs, x, z-1)
		floodfill(oldbio, newbio, biogs, x, z+1)
	}
}

func NewFillTool() *fillTool {
	return new(fillTool)
}
