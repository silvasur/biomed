package main

import (
	"github.com/kch42/gomcmap/mcmap"
)

type BiomeGetSetter interface {
	GetBiomeAt(x, z int) (mcmap.Biome, bool)
	SetBiomeAt(x, z int, bio mcmap.Biome)
}

type XZPos struct {
	X, Z int
}

type Tool interface {
	SingleClick() bool // Whether only one click should be performed (true) or the action should be repeated, if the mouse is dragged
	IsSlow() bool
	Do(bio mcmap.Biome, biogs BiomeGetSetter, x, z int)
}

type drawTool struct {
	radGetter func() int
}

func (d *drawTool) SingleClick() bool { return false }
func (d *drawTool) IsSlow() bool      { return false }

func (d *drawTool) Do(bio mcmap.Biome, biogs BiomeGetSetter, x, z int) {
	rad := d.radGetter()
	if rad <= 0 {
		return
	}

	for xp := x - (rad - 1); xp < x+rad; xp++ {
		for zp := z - (rad - 1); zp < z+rad; zp++ {
			biogs.SetBiomeAt(xp, zp, bio)
		}
	}
}

func NewDrawTool(radGetter func() int) *drawTool {
	return &drawTool{radGetter}
}

type fillTool struct{}

func (f *fillTool) SingleClick() bool { return true }
func (f *fillTool) IsSlow() bool      { return true }

func chkBounds(x, z, xStart, zStart, xEnd, zEnd int) bool {
	return (x >= xStart) && (z >= zStart) && (x < xEnd) && (z < zEnd)
}

func (f *fillTool) Do(bio mcmap.Biome, biogs BiomeGetSetter, x, z int) {
	oldbio, ok := biogs.GetBiomeAt(x, z)
	if (!ok) || (oldbio == bio) {
		return
	}

	inChunkQueue := []XZPos{}
	outOfChunkQueue := []XZPos{{x, z}}

	for {
		oocqL := len(outOfChunkQueue) - 1
		if oocqL < 0 {
			break
		}

		pos := outOfChunkQueue[oocqL]
		inChunkQueue = []XZPos{pos}
		outOfChunkQueue = outOfChunkQueue[:oocqL]

		cx, cz, _, _ := mcmap.BlockToChunk(pos.X, pos.Z)
		xStart := cx * mcmap.ChunkSizeXZ
		zStart := cz * mcmap.ChunkSizeXZ
		xEnd := xStart + mcmap.ChunkSizeXZ
		zEnd := zStart + mcmap.ChunkSizeXZ

		for {
			icqL := len(inChunkQueue) - 1
			if icqL < 0 {
				break
			}

			pos := inChunkQueue[icqL]
			inChunkQueue = inChunkQueue[:icqL]

			px, pz := pos.X, pos.Z

			if haveBio, ok := biogs.GetBiomeAt(px, pz); ok && (haveBio == oldbio) {
				biogs.SetBiomeAt(px, pz, bio)

				nx, nz := px+1, pz
				if chkBounds(nx, nz, xStart, zStart, xEnd, zEnd) {
					inChunkQueue = append(inChunkQueue, XZPos{nx, nz})
				} else {
					outOfChunkQueue = append(outOfChunkQueue, XZPos{nx, nz})
				}

				nx, nz = px-1, pz
				if chkBounds(nx, nz, xStart, zStart, xEnd, zEnd) {
					inChunkQueue = append(inChunkQueue, XZPos{nx, nz})
				} else {
					outOfChunkQueue = append(outOfChunkQueue, XZPos{nx, nz})
				}

				nx, nz = px, pz+1
				if chkBounds(nx, nz, xStart, zStart, xEnd, zEnd) {
					inChunkQueue = append(inChunkQueue, XZPos{nx, nz})
				} else {
					outOfChunkQueue = append(outOfChunkQueue, XZPos{nx, nz})
				}

				nx, nz = px, pz-1
				if chkBounds(nx, nz, xStart, zStart, xEnd, zEnd) {
					inChunkQueue = append(inChunkQueue, XZPos{nx, nz})
				} else {
					outOfChunkQueue = append(outOfChunkQueue, XZPos{nx, nz})
				}
			}
		}
	}
}

func NewFillTool() *fillTool {
	return new(fillTool)
}
