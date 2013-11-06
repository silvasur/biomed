package main

import (
	"errors"
	"fmt"
	"github.com/kch42/gonbt/nbt"
	"io"
)

func invalidLevelDat(err error) error {
	return fmt.Errorf("Invalid level.dat: %s", err)
}

func getMapCenter(leveldat io.Reader) (int, int, error) {
	lvl, _, err := nbt.ReadGzipdNamedTag(leveldat)
	if err != nil {
		return 0, 0, err
	}

	if lvl.Type != nbt.TAG_Compound {
		return 0, 0, invalidLevelDat(errors.New("Root tag has wrong type"))
	}
	root := lvl.Payload.(nbt.TagCompound)

	data, err := root.GetCompound("Data")
	if err != nil {
		return 0, 0, invalidLevelDat(err)
	}

	x, err := data.GetInt("SpawnX")
	if err != nil {
		return 0, 0, invalidLevelDat(err)
	}

	z, err := data.GetInt("SpawnZ")
	if err != nil {
		return 0, 0, invalidLevelDat(err)
	}

	return int(x), int(z), nil
}
