package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

func allMaps() map[string]string {
	savesDir := ""
	switch runtime.GOOS {
	case "linux":
		savesDir = fmt.Sprintf("%s/.minecraft/saves", os.Getenv("HOME"))
	case "darwin":
		savesDir = fmt.Sprintf("%s/Library/Application Support/minecraft/saves", os.Getenv("HOME"))
	case "windows":
		savesDir = fmt.Sprintf(`%s\.minecraft\saves`, os.Getenv("appdata"))
	default:
		return nil
	}

	f, err := os.Open(savesDir)
	if err != nil {
		return nil
	}
	defer f.Close()
	fi, err := f.Stat()
	if (err != nil) || (!fi.IsDir()) {
		return nil
	}

	infos, err := f.Readdir(-1)
	if err != nil {
		return nil
	}

	maps := make(map[string]string)
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		p := path.Join(savesDir, info.Name())

		leveldat := path.Join(p, "level.dat")
		fi, err := os.Stat(leveldat)
		if (err != nil) || (!fi.Mode().IsRegular()) {
			continue
		}

		maps[info.Name()] = leveldat
	}

	return maps
}
