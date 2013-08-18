# biomed

A Minecraft Biome Editor.

## Usage

Open your World. If biomed is able to find your savegame directory, you can use `File > Open World`, otherwise you can use `File > Open` and manually select the `region` directory of your map.

You can toggle between the biome view and a map view by clicking `Show Biomes` in the sidebar.

Select the tool and biome you want to use in the sidebar (Fill or Draw) and click on the map area to use the tool.

Use the middle mouse button to move around in your map.

If you want biomed to fix snow and ice (melt away for warm biomes, add for cold ones), click `Fix Snow/Ice` in the sidebar.

If you are done, click `File > Save`.

## WARNING

Although everythung seems to work, please make a backup of your maps, just in case.

## Dependencies / System Requirements

biomed works good on Linux. It also should work on Windows and Mac OSX. I currently try to compile it for Windows. Since I do not own a Mac, I can't compile / test for OSX.

You need a recent version of GTK 2.x. If you are on Linux, you probably already 

## Compiling

You will need Go (1.1 or greater) and the development files for GTK 2.x.

Then compiling is done by running: `go get github.com/kch42/biomed`
