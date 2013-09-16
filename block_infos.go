package main

import (
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
)

var blockColors = map[mcmap.BlockID]*gdk.Color{
	mcmap.BlkStone:                      gdk.NewColor("#666666"),
	mcmap.BlkGrassBlock:                 gdk.NewColor("#00aa00"),
	mcmap.BlkDirt:                       gdk.NewColor("#644804"),
	mcmap.BlkCobblestone:                gdk.NewColor("#7a7a7a"),
	mcmap.BlkWoodPlanks:                 gdk.NewColor("#a4721c"),
	mcmap.BlkBedrock:                    gdk.NewColor("#111111"),
	mcmap.BlkWater:                      gdk.NewColor("#0000ff"),
	mcmap.BlkStationaryWater:            gdk.NewColor("#0000ff"),
	mcmap.BlkLava:                       gdk.NewColor("#ff4400"),
	mcmap.BlkStationaryLava:             gdk.NewColor("#ff4400"),
	mcmap.BlkSand:                       gdk.NewColor("#f1ee85"),
	mcmap.BlkGravel:                     gdk.NewColor("#9ba3a9"),
	mcmap.BlkGoldOre:                    gdk.NewColor("#ffa200"),
	mcmap.BlkIronOre:                    gdk.NewColor("#e1e1e1"),
	mcmap.BlkCoalOre:                    gdk.NewColor("#333333"),
	mcmap.BlkWood:                       gdk.NewColor("#a4721c"),
	mcmap.BlkLeaves:                     gdk.NewColor("#57a100"),
	mcmap.BlkGlass:                      gdk.NewColor("#eeeeff"),
	mcmap.BlkLapisLazuliOre:             gdk.NewColor("#3114e3"),
	mcmap.BlkLapisLazuliBlock:           gdk.NewColor("#3114e3"),
	mcmap.BlkDispenser:                  gdk.NewColor("#7a7a7a"),
	mcmap.BlkSandstone:                  gdk.NewColor("#f1ee85"),
	mcmap.BlkNoteBlock:                  gdk.NewColor("#a4721c"),
	mcmap.BlkBed:                        gdk.NewColor("#a00000"),
	mcmap.BlkPoweredRail:                gdk.NewColor("#ff0000"),
	mcmap.BlkDetectorRail:               gdk.NewColor("#ff0000"),
	mcmap.BlkStickyPiston:               gdk.NewColor("#91ba12"),
	mcmap.BlkCobweb:                     gdk.NewColor("#dddddd"),
	mcmap.BlkGrass:                      gdk.NewColor("#a0f618"),
	mcmap.BlkPiston:                     gdk.NewColor("#a4721c"),
	mcmap.BlkPistonExtension:            gdk.NewColor("#a4721c"),
	mcmap.BlkWool:                       gdk.NewColor("#ffffff"),
	mcmap.BlkBlockOfGold:                gdk.NewColor("#ffa200"),
	mcmap.BlkBlockOfIron:                gdk.NewColor("#e1e1e1"),
	mcmap.BlkTNT:                        gdk.NewColor("#a20022"),
	mcmap.BlkBookshelf:                  gdk.NewColor("#a4721c"),
	mcmap.BlkMossStone:                  gdk.NewColor("#589b71"),
	mcmap.BlkObsidian:                   gdk.NewColor("#111144"),
	mcmap.BlkTorch:                      gdk.NewColor("#ffcc00"),
	mcmap.BlkFire:                       gdk.NewColor("#ffcc00"),
	mcmap.BlkMonsterSpawner:             gdk.NewColor("#344e6a"),
	mcmap.BlkOakWoodStairs:              gdk.NewColor("#a4721c"),
	mcmap.BlkChest:                      gdk.NewColor("#a4721c"),
	mcmap.BlkRedstoneWire:               gdk.NewColor("#ff0000"),
	mcmap.BlkDiamondOre:                 gdk.NewColor("#00fff6"),
	mcmap.BlkBlockOfDiamond:             gdk.NewColor("#00fff6"),
	mcmap.BlkCraftingTable:              gdk.NewColor("#a4721c"),
	mcmap.BlkWheat:                      gdk.NewColor("#e7ae00"),
	mcmap.BlkFarmland:                   gdk.NewColor("#644804"),
	mcmap.BlkFurnace:                    gdk.NewColor("#7a7a7a"),
	mcmap.BlkBurningFurnace:             gdk.NewColor("#7a7a7a"),
	mcmap.BlkSignPost:                   gdk.NewColor("#a4721c"),
	mcmap.BlkWoodenDoor:                 gdk.NewColor("#a4721c"),
	mcmap.BlkLadders:                    gdk.NewColor("#a4721c"),
	mcmap.BlkRail:                       gdk.NewColor("#dbdbdb"),
	mcmap.BlkCobblestoneStairs:          gdk.NewColor("#7a7a7a"),
	mcmap.BlkWallSign:                   gdk.NewColor("#a4721c"),
	mcmap.BlkLever:                      gdk.NewColor("#a4721c"),
	mcmap.BlkStonePressurePlate:         gdk.NewColor("#666666"),
	mcmap.BlkIronDoor:                   gdk.NewColor("#e1e1e1"),
	mcmap.BlkWoodenPressurePlate:        gdk.NewColor("#a4721c"),
	mcmap.BlkRedstoneOre:                gdk.NewColor("#a00000"),
	mcmap.BlkGlowingRedstoneOre:         gdk.NewColor("#ff0000"),
	mcmap.BlkRedstoneTorchInactive:      gdk.NewColor("#ff0000"),
	mcmap.BlkRedstoneTorchActive:        gdk.NewColor("#ff0000"),
	mcmap.BlkStoneButton:                gdk.NewColor("#666666"),
	mcmap.BlkSnow:                       gdk.NewColor("#e5fffe"),
	mcmap.BlkIce:                        gdk.NewColor("#9fdcff"),
	mcmap.BlkSnowBlock:                  gdk.NewColor("#e5fffe"),
	mcmap.BlkCactus:                     gdk.NewColor("#01bc3a"),
	mcmap.BlkClay:                       gdk.NewColor("#767a82"),
	mcmap.BlkSugarCane:                  gdk.NewColor("#12db50"),
	mcmap.BlkJukebox:                    gdk.NewColor("#a4721c"),
	mcmap.BlkFence:                      gdk.NewColor("#a4721c"),
	mcmap.BlkPumpkin:                    gdk.NewColor("#ff7000"),
	mcmap.BlkNetherrack:                 gdk.NewColor("#851c2d"),
	mcmap.BlkSoulSand:                   gdk.NewColor("#796a59"),
	mcmap.BlkGlowstone:                  gdk.NewColor("#ffff00"),
	mcmap.BlkNetherPortal:               gdk.NewColor("#ff00ff"),
	mcmap.BlkJackOLantern:               gdk.NewColor("#ff7000"),
	mcmap.BlkRedstoneRepeaterInactive:   gdk.NewColor("#ff0000"),
	mcmap.BlkRedstoneRepeaterActive:     gdk.NewColor("#ff0000"),
	mcmap.BlkTrapdoor:                   gdk.NewColor("#a4721c"),
	mcmap.BlkStoneBricks:                gdk.NewColor("#666666"),
	mcmap.BlkHugeBrownMushroom:          gdk.NewColor("#b07859"),
	mcmap.BlkHugeRedMushroom:            gdk.NewColor("#dd0000"),
	mcmap.BlkIronBars:                   gdk.NewColor("#e1e1e1"),
	mcmap.BlkGlassPane:                  gdk.NewColor("#eeeeff"),
	mcmap.BlkMelon:                      gdk.NewColor("#9ac615"),
	mcmap.BlkVines:                      gdk.NewColor("#50720d"),
	mcmap.BlkFenceGate:                  gdk.NewColor("#a4721c"),
	mcmap.BlkBrickStairs:                gdk.NewColor("#c42500"),
	mcmap.BlkStoneBrickStairs:           gdk.NewColor("#666666"),
	mcmap.BlkMycelium:                   gdk.NewColor("#7c668c"),
	mcmap.BlkLilyPad:                    gdk.NewColor("#50720d"),
	mcmap.BlkNetherBrick:                gdk.NewColor("#c42500"),
	mcmap.BlkNetherBrickFence:           gdk.NewColor("#c42500"),
	mcmap.BlkNetherBrickStairs:          gdk.NewColor("#c42500"),
	mcmap.BlkEnchantmentTable:           gdk.NewColor("#222244"),
	mcmap.BlkBrewingStand:               gdk.NewColor("#666666"),
	mcmap.BlkCauldron:                   gdk.NewColor("#666666"),
	mcmap.BlkEndPortal:                  gdk.NewColor("#000000"),
	mcmap.BlkEndPortalBlock:             gdk.NewColor("#e0dbce"),
	mcmap.BlkEndStone:                   gdk.NewColor("#e0dbce"),
	mcmap.BlkRedstoneLampInactive:       gdk.NewColor("#ffff00"),
	mcmap.BlkRedstoneLampActive:         gdk.NewColor("#ffff00"),
	mcmap.BlkSandstoneStairs:            gdk.NewColor("#f1ee85"),
	mcmap.BlkEmeraldOre:                 gdk.NewColor("#00c140"),
	mcmap.BlkEnderChest:                 gdk.NewColor("#222244"),
	mcmap.BlkBlockOfEmerald:             gdk.NewColor("#00c140"),
	mcmap.BlkSpruceWoodStairs:           gdk.NewColor("#a4721c"),
	mcmap.BlkBirchWoodStairs:            gdk.NewColor("#a4721c"),
	mcmap.BlkJungleWoodStairs:           gdk.NewColor("#a4721c"),
	mcmap.BlkCommandBlock:               gdk.NewColor("#e8ec78"),
	mcmap.BlkBeacon:                     gdk.NewColor("#00fff6"),
	mcmap.BlkCobblestoneWall:            gdk.NewColor("#7a7a7a"),
	mcmap.BlkCarrots:                    gdk.NewColor("#ff6000"),
	mcmap.BlkPotatoes:                   gdk.NewColor("#c6cd0c"),
	mcmap.BlkWoodenButton:               gdk.NewColor("#a4721c"),
	mcmap.BlkAnvil:                      gdk.NewColor("#444444"),
	mcmap.BlkTrappedChest:               gdk.NewColor("#a4721c"),
	mcmap.BlkRedstoneComparatorInactive: gdk.NewColor("#ff0000"),
	mcmap.BlkRedstoneComparatorActive:   gdk.NewColor("#ff0000"),
	mcmap.BlkBlockOfRedstone:            gdk.NewColor("#ff0000"),
	mcmap.BlkNetherQuartzOre:            gdk.NewColor("#e7e7e7"),
	mcmap.BlkHopper:                     gdk.NewColor("#444444"),
	mcmap.BlkBlockOfQuartz:              gdk.NewColor("#e7e7e7"),
	mcmap.BlkQuartzStairs:               gdk.NewColor("#e7e7e7"),
	mcmap.BlkActivatorRail:              gdk.NewColor("#ff0000"),
	mcmap.BlkDropper:                    gdk.NewColor("#444444"),
	mcmap.BlkStainedClay:                gdk.NewColor("#767a82"),
	mcmap.BlkHayBlock:                   gdk.NewColor("#e7ae00"),
	mcmap.BlkCarpet:                     gdk.NewColor("#ffffff"),
	mcmap.BlkHardenedClay:               gdk.NewColor("#767a82"),
	mcmap.BlkBlockOfCoal:                gdk.NewColor("#333333"),
	mcmap.BlkPackedIce:                  gdk.NewColor("#63bff4"),
}

var blockCanSnowIn = map[mcmap.BlockID]bool{
	mcmap.BlkAir:                        false,
	mcmap.BlkStone:                      true,
	mcmap.BlkGrassBlock:                 true,
	mcmap.BlkDirt:                       true,
	mcmap.BlkCobblestone:                true,
	mcmap.BlkWoodPlanks:                 true,
	mcmap.BlkSaplings:                   false,
	mcmap.BlkBedrock:                    true,
	mcmap.BlkWater:                      false,
	mcmap.BlkStationaryWater:            false,
	mcmap.BlkLava:                       false,
	mcmap.BlkStationaryLava:             false,
	mcmap.BlkSand:                       true,
	mcmap.BlkGravel:                     true,
	mcmap.BlkGoldOre:                    true,
	mcmap.BlkIronOre:                    true,
	mcmap.BlkCoalOre:                    true,
	mcmap.BlkWood:                       true,
	mcmap.BlkLeaves:                     true,
	mcmap.BlkSponge:                     true,
	mcmap.BlkGlass:                      false,
	mcmap.BlkLapisLazuliOre:             true,
	mcmap.BlkLapisLazuliBlock:           true,
	mcmap.BlkDispenser:                  true,
	mcmap.BlkSandstone:                  true,
	mcmap.BlkNoteBlock:                  true,
	mcmap.BlkBed:                        false,
	mcmap.BlkPoweredRail:                false,
	mcmap.BlkDetectorRail:               false,
	mcmap.BlkStickyPiston:               true,
	mcmap.BlkCobweb:                     false,
	mcmap.BlkGrass:                      false,
	mcmap.BlkDeadBush:                   false,
	mcmap.BlkPiston:                     true,
	mcmap.BlkPistonExtension:            false,
	mcmap.BlkWool:                       true,
	mcmap.BlkBlockMovedByPiston:         false,
	mcmap.BlkDandelion:                  false,
	mcmap.BlkFlower:                     false,
	mcmap.BlkBrownMushroom:              false,
	mcmap.BlkRedMushroom:                false,
	mcmap.BlkBlockOfGold:                true,
	mcmap.BlkBlockOfIron:                true,
	mcmap.BlkDoubleSlabs:                true,
	mcmap.BlkSlabs:                      false,
	mcmap.BlkBricks:                     true,
	mcmap.BlkTNT:                        true,
	mcmap.BlkBookshelf:                  true,
	mcmap.BlkMossStone:                  true,
	mcmap.BlkObsidian:                   true,
	mcmap.BlkTorch:                      false,
	mcmap.BlkFire:                       false,
	mcmap.BlkMonsterSpawner:             true,
	mcmap.BlkOakWoodStairs:              true,
	mcmap.BlkChest:                      false,
	mcmap.BlkRedstoneWire:               false,
	mcmap.BlkDiamondOre:                 true,
	mcmap.BlkBlockOfDiamond:             true,
	mcmap.BlkCraftingTable:              true,
	mcmap.BlkWheat:                      false,
	mcmap.BlkFarmland:                   false,
	mcmap.BlkFurnace:                    false,
	mcmap.BlkBurningFurnace:             false,
	mcmap.BlkSignPost:                   false,
	mcmap.BlkWoodenDoor:                 false,
	mcmap.BlkLadders:                    false,
	mcmap.BlkRail:                       false,
	mcmap.BlkCobblestoneStairs:          false,
	mcmap.BlkWallSign:                   false,
	mcmap.BlkLever:                      false,
	mcmap.BlkStonePressurePlate:         false,
	mcmap.BlkIronDoor:                   false,
	mcmap.BlkWoodenPressurePlate:        false,
	mcmap.BlkRedstoneOre:                true,
	mcmap.BlkGlowingRedstoneOre:         true,
	mcmap.BlkRedstoneTorchInactive:      false,
	mcmap.BlkRedstoneTorchActive:        false,
	mcmap.BlkStoneButton:                false,
	mcmap.BlkSnow:                       false,
	mcmap.BlkIce:                        false,
	mcmap.BlkSnowBlock:                  true,
	mcmap.BlkCactus:                     false,
	mcmap.BlkClay:                       true,
	mcmap.BlkSugarCane:                  false,
	mcmap.BlkJukebox:                    true,
	mcmap.BlkFence:                      false,
	mcmap.BlkPumpkin:                    true,
	mcmap.BlkNetherrack:                 true,
	mcmap.BlkSoulSand:                   true,
	mcmap.BlkGlowstone:                  true,
	mcmap.BlkNetherPortal:               false,
	mcmap.BlkJackOLantern:               true,
	mcmap.BlkCakeBlock:                  false,
	mcmap.BlkRedstoneRepeaterInactive:   false,
	mcmap.BlkRedstoneRepeaterActive:     false,
	mcmap.BlkLockedChest:                false,
	mcmap.BlkTrapdoor:                   false,
	mcmap.BlkMonsterEgg:                 false,
	mcmap.BlkStoneBricks:                true,
	mcmap.BlkHugeBrownMushroom:          true,
	mcmap.BlkHugeRedMushroom:            true,
	mcmap.BlkIronBars:                   false,
	mcmap.BlkGlassPane:                  false,
	mcmap.BlkMelon:                      true,
	mcmap.BlkPumpkinStem:                false,
	mcmap.BlkMelonStem:                  false,
	mcmap.BlkVines:                      false,
	mcmap.BlkFenceGate:                  false,
	mcmap.BlkBrickStairs:                false,
	mcmap.BlkStoneBrickStairs:           false,
	mcmap.BlkMycelium:                   true,
	mcmap.BlkLilyPad:                    false,
	mcmap.BlkNetherBrick:                true,
	mcmap.BlkNetherBrickFence:           false,
	mcmap.BlkNetherBrickStairs:          false,
	mcmap.BlkNetherWart:                 false,
	mcmap.BlkEnchantmentTable:           false,
	mcmap.BlkBrewingStand:               false,
	mcmap.BlkCauldron:                   false,
	mcmap.BlkEndPortal:                  false,
	mcmap.BlkEndPortalBlock:             false,
	mcmap.BlkEndStone:                   true,
	mcmap.BlkDragonEgg:                  false,
	mcmap.BlkRedstoneLampInactive:       true,
	mcmap.BlkRedstoneLampActive:         true,
	mcmap.BlkWoodenDoubleSlab:           true,
	mcmap.BlkWoodenSlab:                 false,
	mcmap.BlkCocoa:                      false,
	mcmap.BlkSandstoneStairs:            false,
	mcmap.BlkEmeraldOre:                 true,
	mcmap.BlkEnderChest:                 false,
	mcmap.BlkTripwireHook:               false,
	mcmap.BlkTripwire:                   false,
	mcmap.BlkBlockOfEmerald:             true,
	mcmap.BlkSpruceWoodStairs:           false,
	mcmap.BlkBirchWoodStairs:            false,
	mcmap.BlkJungleWoodStairs:           false,
	mcmap.BlkCommandBlock:               true,
	mcmap.BlkBeacon:                     false,
	mcmap.BlkCobblestoneWall:            false,
	mcmap.BlkFlowerPot:                  false,
	mcmap.BlkCarrots:                    false,
	mcmap.BlkPotatoes:                   false,
	mcmap.BlkWoodenButton:               false,
	mcmap.BlkMobHead:                    false,
	mcmap.BlkAnvil:                      false,
	mcmap.BlkTrappedChest:               false,
	mcmap.BlkWeightedPressurePlateLight: false,
	mcmap.BlkWeightedPressurePlateHeavy: false,
	mcmap.BlkRedstoneComparatorInactive: false,
	mcmap.BlkRedstoneComparatorActive:   false,
	mcmap.BlkDaylightSensor:             false,
	mcmap.BlkBlockOfRedstone:            true,
	mcmap.BlkNetherQuartzOre:            true,
	mcmap.BlkHopper:                     false,
	mcmap.BlkBlockOfQuartz:              true,
	mcmap.BlkQuartzStairs:               false,
	mcmap.BlkActivatorRail:              false,
	mcmap.BlkDropper:                    false,
	mcmap.BlkStainedClay:                true,
	mcmap.BlkHayBlock:                   true,
	mcmap.BlkCarpet:                     false,
	mcmap.BlkHardenedClay:               true,
	mcmap.BlkBlockOfCoal:                true,
	mcmap.BlkPackedIce:                  false,
	mcmap.BlkLargeFlower:                false,
}
