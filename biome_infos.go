package main

import (
	"github.com/kch42/gomcmap/mcmap"
	"github.com/kch42/kagus"
)

var bioList = []mcmap.Biome{
	mcmap.BioOcean,
	mcmap.BioPlains,
	mcmap.BioDesert,
	mcmap.BioExtremeHills,
	mcmap.BioForest,
	mcmap.BioTaiga,
	mcmap.BioSwampland,
	mcmap.BioRiver,
	mcmap.BioHell,
	mcmap.BioSky,
	mcmap.BioFrozenOcean,
	mcmap.BioFrozenRiver,
	mcmap.BioIcePlains,
	mcmap.BioIceMountains,
	mcmap.BioMushroomIsland,
	mcmap.BioMushroomIslandShore,
	mcmap.BioBeach,
	mcmap.BioDesertHills,
	mcmap.BioForestHills,
	mcmap.BioTaigaHills,
	mcmap.BioExtremeHillsEdge,
	mcmap.BioJungle,
	mcmap.BioJungleHills,
	mcmap.BioUncalculated,
}

var bioColors = map[mcmap.Biome]kagus.RGB0x{
	mcmap.BioOcean:               0x0000ff,
	mcmap.BioPlains:              0x9fe804,
	mcmap.BioDesert:              0xf5ff58,
	mcmap.BioExtremeHills:        0xa75300,
	mcmap.BioForest:              0x006f2a,
	mcmap.BioTaiga:               0x05795a,
	mcmap.BioSwampland:           0x6a7905,
	mcmap.BioRiver:               0x196eff,
	mcmap.BioHell:                0xd71900,
	mcmap.BioSky:                 0x871eb3,
	mcmap.BioFrozenOcean:         0xd6f0ff,
	mcmap.BioFrozenRiver:         0x8fb6cd,
	mcmap.BioIcePlains:           0xfbfbfb,
	mcmap.BioIceMountains:        0xc6bfb1,
	mcmap.BioMushroomIsland:      0x9776a4,
	mcmap.BioMushroomIslandShore: 0x9e8ebc,
	mcmap.BioBeach:               0xfffdc9,
	mcmap.BioDesertHills:         0xadb354,
	mcmap.BioForestHills:         0x40694f,
	mcmap.BioTaigaHills:          0x5b8578,
	mcmap.BioExtremeHillsEdge:    0xa77748,
	mcmap.BioJungle:              0x22db04,
	mcmap.BioJungleHills:         0x63bf54,
	mcmap.BioUncalculated:        0x333333,
}
