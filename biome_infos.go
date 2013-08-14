package main

import (
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
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

var bioColors = map[mcmap.Biome]*gdk.Color{
	mcmap.BioOcean:               gdk.NewColor("#0000ff"),
	mcmap.BioPlains:              gdk.NewColor("#9fe804"),
	mcmap.BioDesert:              gdk.NewColor("#f5ff58"),
	mcmap.BioExtremeHills:        gdk.NewColor("#a75300"),
	mcmap.BioForest:              gdk.NewColor("#006f2a"),
	mcmap.BioTaiga:               gdk.NewColor("#05795a"),
	mcmap.BioSwampland:           gdk.NewColor("#6a7905"),
	mcmap.BioRiver:               gdk.NewColor("#196eff"),
	mcmap.BioHell:                gdk.NewColor("#d71900"),
	mcmap.BioSky:                 gdk.NewColor("#871eb3"),
	mcmap.BioFrozenOcean:         gdk.NewColor("#d6f0ff"),
	mcmap.BioFrozenRiver:         gdk.NewColor("#8fb6cd"),
	mcmap.BioIcePlains:           gdk.NewColor("#fbfbfb"),
	mcmap.BioIceMountains:        gdk.NewColor("#c6bfb1"),
	mcmap.BioMushroomIsland:      gdk.NewColor("#9776a4"),
	mcmap.BioMushroomIslandShore: gdk.NewColor("#9e8ebc"),
	mcmap.BioBeach:               gdk.NewColor("#fffdc9"),
	mcmap.BioDesertHills:         gdk.NewColor("#adb354"),
	mcmap.BioForestHills:         gdk.NewColor("#40694f"),
	mcmap.BioTaigaHills:          gdk.NewColor("#5b8578"),
	mcmap.BioExtremeHillsEdge:    gdk.NewColor("#a77748"),
	mcmap.BioJungle:              gdk.NewColor("#22db04"),
	mcmap.BioJungleHills:         gdk.NewColor("#63bf54"),
	mcmap.BioUncalculated:        gdk.NewColor("#333333"),
}
