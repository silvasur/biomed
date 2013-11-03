package main

import (
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/kch42/kagus"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	"strconv"
	"unicode"
)

type biomeEditFrame struct {
	*gtk.Frame
	applyBtn                          *gtk.Button
	idInput, snowLineInput, nameInput *gtk.Entry
	colorInput                        *gtk.ColorButton
}

func newBiomeEditFrame() *biomeEditFrame {
	frm := &biomeEditFrame{
		Frame:         gtk.NewFrame("Edit Biome"),
		applyBtn:      gtk.NewButtonWithLabel("Apply"),
		idInput:       gtk.NewEntry(),
		snowLineInput: gtk.NewEntry(),
		nameInput:     gtk.NewEntry(),
		colorInput:    gtk.NewColorButton(),
	}

	frm.idInput.SetSizeRequest(40, -1)
	frm.snowLineInput.SetSizeRequest(40, -1)
	
	frm.idInput.Connect("changed", frm.unlockApply)
	frm.nameInput.Connect("changed", frm.unlockApply)
	frm.snowLineInput.Connect("changed", frm.unlockApply)
	frm.applyBtn.SetSensitive(false)

	vbox := gtk.NewVBox(false, 0)
	hbox := gtk.NewHBox(false, 0)

	hbox.PackStart(gtk.NewLabel("Color:"), false, false, 0)
	hbox.PackStart(frm.colorInput, false, false, 3)
	hbox.PackStart(gtk.NewLabel("ID:"), false, false, 0)
	hbox.PackStart(frm.idInput, false, false, 3)
	hbox.PackStart(gtk.NewLabel("Snowline:"), false, false, 0)
	hbox.PackStart(frm.snowLineInput, false, false, 3)
	hbox.PackStart(gtk.NewLabel("Name:"), false, false, 0)
	hbox.PackStart(frm.nameInput, true, true, 3)

	vbox.PackStart(hbox, false, false, 0)
	vbox.PackStart(frm.applyBtn, false, false, 3)
	frm.Add(vbox)

	return frm
}

func (frm *biomeEditFrame) setBiomeInfo(info BiomeInfo) {
	frm.colorInput.SetColor(gdk.NewColor(info.Color))
	frm.idInput.SetText(strconv.FormatInt(int64(info.ID), 10))
	frm.snowLineInput.SetText(strconv.FormatInt(int64(info.SnowLine), 10))
	frm.nameInput.SetText(info.Name)
}

func (frm *biomeEditFrame) getBiomeInfo() (BiomeInfo, bool) {
	id, err := strconv.ParseUint(frm.idInput.GetText(), 10, 8)
	if err != nil {
		return BiomeInfo{}, false
	}

	snow, err := strconv.ParseInt(frm.snowLineInput.GetText(), 10, 32)
	if err != nil {
		return BiomeInfo{}, false
	}
	if (snow > mcmap.ChunkSizeY) || (snow < 0) {
		snow = mcmap.ChunkSizeY
	}

	name := frm.nameInput.GetText()
	if name != "" {
		return BiomeInfo{}, false
	}

	col := frm.colorInput.GetColor()

	return BiomeInfo{
		ID:       mcmap.Biome(id),
		SnowLine: int(snow),
		Name:     name,
		Color:    fmt.Sprintf("#%02x%02x%02x", col.Red()<<8, col.Green()<<8, col.Blue()<<8),
	}, true
}

func (frm *biomeEditFrame) checkOK() bool {
	if id := frm.idInput.GetText(); (id == "") || (!kagus.StringConsistsOf(id, unicode.IsNumber)) {
		return false
	}

	if snow := frm.snowLineInput.GetText(); (snow == "") || (!kagus.StringConsistsOf(snow, unicode.IsNumber)) {
		return false
	}

	return frm.nameInput.GetText() != ""
}

func (frm *biomeEditFrame) unlockApply() {
	frm.applyBtn.SetSensitive(frm.checkOK())
}

type BiomeInfoEditor struct {
	*gtk.Dialog
	biomes []BiomeInfo
}

func NewBiomeInfoEditor(biomes []BiomeInfo) *BiomeInfoEditor {
	ed := &BiomeInfoEditor{
		Dialog: gtk.NewDialog(),
		biomes: biomes,
	}

	ed.SetModal(true)
	vbox := ed.GetVBox()

	btnHBox := gtk.NewHBox(true, 0)

	resetBtn := gtk.NewButtonWithLabel("Reset to defaults")
	resetBtn.Connect("clicked", ed.reset)
	loadBtn := gtk.NewButtonWithLabel("Load from file ...")
	loadBtn.Connect("clicked", ed.load)
	saveBtn := gtk.NewButtonWithLabel("Save to file ...")
	saveBtn.Connect("clicked", ed.save)

	btnHBox.PackStart(resetBtn, true, true, 3)
	btnHBox.PackStart(loadBtn, true, true, 3)
	btnHBox.PackStart(saveBtn, true, true, 3)
	vbox.PackStart(btnHBox, false, false, 3)

	editFrame := newBiomeEditFrame()
	vbox.PackStart(editFrame, false, false, 3)

	ed.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	ed.AddButton("OK", gtk.RESPONSE_OK)
	ed.ShowAll()
	return ed
}

func (ed *BiomeInfoEditor) reset() {
	ed.biomes = ReadDefaultBiomes()
	// TODO: Update view
}

func (ed *BiomeInfoEditor) load() {
	// TODO
}

func (ed *BiomeInfoEditor) save() {
	// TODO
}

func (ed *BiomeInfoEditor) Biomes() []BiomeInfo { return ed.biomes }
