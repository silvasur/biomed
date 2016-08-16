package main

import (
	"errors"
	"fmt"
	"github.com/silvasur/gomcmap/mcmap"
)

type CachedRegion struct {
	Region      *mcmap.Region
	cacheChunks []*mcmap.Chunk
	cachePos    []XZPos
	cachesize   int
}

func NewCachedRegion(reg *mcmap.Region, cachesize int) *CachedRegion {
	if cachesize <= 0 {
		panic(errors.New("Cachesize must be >0"))
	}
	return &CachedRegion{
		Region:      reg,
		cacheChunks: make([]*mcmap.Chunk, cachesize),
		cachePos:    make([]XZPos, cachesize),
		cachesize:   cachesize,
	}
}

func (cr *CachedRegion) Chunk(x, z int) (*mcmap.Chunk, error) {
	pos := XZPos{x, z}

	for i, p := range cr.cachePos {
		if p == pos {
			if cr.cacheChunks[i] != nil {
				chunk := cr.cacheChunks[i]
				for j := i; j >= 1; j-- {
					cr.cacheChunks[j] = cr.cacheChunks[j-1]
					cr.cachePos[j] = cr.cachePos[j-1]
				}
				cr.cacheChunks[0] = chunk
				cr.cachePos[0] = pos
				return chunk, nil
			}
		}
	}

	chunk, err := cr.Region.Chunk(x, z)
	if err != nil {
		return nil, err
	}

	if cr.cacheChunks[cr.cachesize-1] != nil {
		if err := cr.cacheChunks[cr.cachesize-1].MarkUnused(); err != nil {
			return nil, fmt.Errorf("Could not remove oldest cache element: %s", err)
		}
	}

	for i := cr.cachesize - 1; i >= 1; i-- {
		cr.cacheChunks[i] = cr.cacheChunks[i-1]
		cr.cachePos[i] = cr.cachePos[i-1]
	}
	cr.cacheChunks[0] = chunk
	cr.cachePos[0] = pos

	return chunk, nil
}

func (cr *CachedRegion) Flush() error {
	for i, chunk := range cr.cacheChunks {
		if chunk == nil {
			continue
		}

		if err := chunk.MarkUnused(); err != nil {
			return err
		}
		cr.cacheChunks[i] = nil
	}

	return nil
}
