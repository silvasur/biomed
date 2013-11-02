package main

import (
	"bufio"
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"io"
	"strconv"
	"strings"
)

type BiomeInfo struct {
	ID       mcmap.Biome
	SnowLine int
	Color    string
	Name     string
}

func ReadBiomeInfos(r io.Reader) ([]BiomeInfo, error) {
	var biomes []BiomeInfo

	sc := bufio.NewScanner(r)
	for i := 1; sc.Scan(); i++ {
		line := sc.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)

		if len(parts) != 4 {
			return nil, fmt.Errorf("Line %d corrupted: Not enough parts", i)
		}

		id, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("Line %d corrupted: %s", i, err)
		}

		snow, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Line %d corrupted: %s", i, err)
		}
		if (snow >= mcmap.ChunkSizeY) || (snow < 0) {
			snow = mcmap.ChunkSizeY
		}

		info := BiomeInfo{
			ID:       mcmap.Biome(id),
			SnowLine: int(snow),
			Color:    parts[2],
			Name:     parts[3],
		}
		biomes = append(biomes, info)
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return biomes, nil
}

func WriteBiomeInfos(biomes []BiomeInfo, w io.Writer) error {
	for _, bio := range biomes {
		if _, err := fmt.Fprintf(w, "%d\t%d\t%s\t%s\n", bio.ID, bio.SnowLine, bio.Color, bio.Name); err != nil {
			return err
		}
	}

	return nil
}

type BiomeLookup map[mcmap.Biome]BiomeInfo

var colBlack = gdk.NewColor("#000000")
var colBuf = make(ColorBuffer)

func (bl BiomeLookup) Color(bio mcmap.Biome) *gdk.Color {
	if info, ok := bl[bio]; ok {
		return colBuf.Color(info.Color)
	}

	return colBlack
}

func (bl BiomeLookup) SnowLine(bio mcmap.Biome) int {
	if info, ok := bl[bio]; ok {
		return info.SnowLine
	}

	return mcmap.ChunkSizeY
}

func (bl BiomeLookup) Name(bio mcmap.Biome) string {
	if info, ok := bl[bio]; ok {
		return info.Name
	}

	return "?"
}

func MkBiomeLookup(biomes []BiomeInfo) BiomeLookup {
	lookup := make(BiomeLookup)
	for _, biome := range biomes {
		lookup[biome.ID] = biome
	}
	return lookup
}

var defaultBiomes = `0	-1	#0000ff	Ocean
1	-1	#9fe804	Plains
2	-1	#f5ff58	Desert
3	95	#a75300	Extreme Hills
4	-1	#006f2a	Forest
5	-1	#05795a	Taiga
6	-1	#6a7905	Swampland
7	-1	#196eff	River
8	-1	#d71900	Hell
9	-1	#871eb3	Sky
10	0	#d6f0ff	Frozen Ocean
11	0	#8fb6cd	Frozen River
12	0	#fbfbfb	Ice Plains
13	0	#c6bfb1	Ice Mountains
14	-1	#9776a4	Mushroom Island
15	-1	#9e8ebc	Mushroom Island Shore
16	-1	#fffdc9	Beach
17	-1	#adb354	Desert Hills
18	-1	#40694f	Forest Hills
19	0	#5b8578	Taiga Hills
20	95	#a77748	Extreme Hills Edge
21	-1	#22db04	Jungle
22	-1	#63bf54	Jungle Hills
23	-1	#40ba2c	Jungle Edge
24	-1	#0000b3	Deep Ocean
25	-1	#9292a6	Stone Beach
26	0	#c7c7e8	Cold Beach
27	-1	#1d964b	Birch Forest
28	-1	#498045	Birch Forest Hills
29	-1	#075a26	Roofed Forest
30	0	#1b948e	Cold Taiga
31	0	#1d7a76	Cold Taiga Hills
32	-1	#1f8f68	Mega Taiga
33	-1	#217a5c	Mega Taiga Hills
34	95	#d76a00	Extreme Hills+
35	-1	#b2bc0f	Savanna
36	-1	#aba60e	Savanna Plateau
37	-1	#ff6c00	Mesa
38	-1	#d9691e	Mesa Plateau F
39	-1	#d95b07	Mesa Plateau
40	-1	#ffd504	Sunflower Plains
41	-1	#f4ff3f	Desert M
42	95	#8c4500	Extreme Hills M
43	-1	#e02f4a	Flower Forest
44	-1	#0a6148	Taiga M
45	-1	#58630e	Swampland M
46	0	#ace8e8	Ice Plains Spikes
47	0	#91cccc	Ice Mountains Spikes
48	-1	#30ba07	Jungle M
49	-1	#3e9130	JungleEdge M
50	-1	#228548	Birch Forest M
51	-1	#2b7547	Birch Forest Hills M
52	-1	#1a5428	Roofed Forest M
53	0	#0f706b	Cold Taiga M
54	-1	#198058	Mega Spruce Taiga
55	-1	#156e4c	Mega Spruce Taiga Hills
56	95	#ba5c00	Extreme Hills+ M
57	-1	#858111	Savanna M
58	-1	#87830b	Savanna Plateau M
59	-1	#ff5100	Mesa (Bryce)
60	-1	#ba5a1a	Mesa Plateau F M
61	-1	#ba4e06	Mesa Plateau M
255	-1	#333333	(Uncalculated)
`

func ReadDefaultBiomes() []BiomeInfo {
	r := strings.NewReader(defaultBiomes)
	biomes, err := ReadBiomeInfos(r)
	if err != nil {
		panic(err)
	}
	return biomes
}
