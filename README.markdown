# biomed

A Minecraft Biome Editor.

## Usage

Open your World. If biomed is able to find your savegame directory, you can use `File > Open World`, otherwise you can use `File > Open` and manually select the `region` directory of your map.

You can toggle between the biome view and a map view by clicking `Show Biomes` in the sidebar.

Select the tool and biome you want to use in the sidebar (Fill or Draw) and click on the map area to use the tool.

Use the middle mouse button to move around in your map.

If you want biomed to fix snow and ice (melt away for warm biomes, add for cold ones), click `Fix Snow/Ice` in the sidebar.

If you are done, click `File > Save`.

## Minecraft 1.7

The first snapshots for Minecraft 1.7 are out and introduce a ton of new biomes. biomed can handle these, but currently not with the "official" version.

This Git repository contains a `minecraft-1.7` branch that contains these changes. If you want to try these out, you'll need to checkout the `minecraft-1.7` branch and (re-)compile biomed, which is done by running `go build` in the biomed directory.

I will not merge these changes into master before Minecraft 1.7 is released, since biomed versions from master should work with the currently official Minecraft release and I don't know how Minecraft will handle unknown biomes.

## WARNING

Although everything seems to work, please make a backup of your maps, just in case.

## Precompiled versions

* Linux (64bit): [biomed-linux64.tar.bz2](http://kch42.de/progs/biomed/biomed-linux64.tar.bz2)
* Windows (32bit, tested with Win8): [biomed-win32.zip](biomed-win32.zip)

## Dependencies / System Requirements

You need a recent version of GTK 2.x. If you are on Linux, you probably already have it installed (if not, you should be able to find it in your package manager). For Windows you need the [GTK all in one bundle](http://ftp.gnome.org/pub/gnome/binaries/win32/gtk+/2.24/gtk+-bundle_2.24.10-20120208_win32.zip), unpack all `*.dll` files from the `bin` directory in the directory that contains `biomed.exe`.

biomed works on Linux and Windows. It also should work on Mac OSX, but since I do not own a Mac, I can't test it.

## Compiling

### Linux

You will need Go (1.1 or greater) and the development files for GTK 2.x. Both should be available from your package manager.

Create a directory that contains the directories `bin`, `src` and `pkg` and set the `GOPATH` environment variable to that directory.

Then compiling is done by running: `go get github.com/kch42/biomed`

### Windows

Install Go (1.1 or greater) from [golang.org](http://golang.org/doc/install).

Create a directory that contains the directories `bin`, `src` and `pkg` and set the `GOPATH` environment variable to that directory.

Follow these instructions for installing go-gtk: [http://stackoverflow.com/questions/16999498/how-to-set-up-gtk-for-go/17042596#17042596](http://stackoverflow.com/questions/16999498/how-to-set-up-gtk-for-go/17042596#17042596)

Then compiling is done by running: `go get github.com/kch42/biomed`
