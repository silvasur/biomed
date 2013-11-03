package main

import (
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"strconv"
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
	if name == "" {
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
	_, ok := frm.getBiomeInfo()
	return ok
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

func mkBioFFilters() (*gtk.FileFilter, *gtk.FileFilter) {
	f1 := gtk.NewFileFilter()
	f1.AddPattern("*.biomes")
	f1.SetName("Biome Infos (.biomes)")

	f2 := gtk.NewFileFilter()
	f2.AddPattern("*")
	f2.SetName("All files")

	return f1, f2
}

func errdlg(msg string, params ...interface{}) {
	dlg := gtk.NewMessageDialog(nil, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, msg, params...)
	dlg.Run()
	dlg.Destroy()
}

func (ed *BiomeInfoEditor) load() {
	f1, f2 := mkBioFFilters()
	dlg := gtk.NewFileChooserDialog("Load", nil, gtk.FILE_CHOOSER_ACTION_OPEN, "OK", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
	dlg.AddFilter(f1)
	dlg.AddFilter(f2)
	defer dlg.Destroy()
	if dlg.Run() == gtk.RESPONSE_OK {
		path := dlg.GetFilename()

		f, err := os.Open(path)
		if err != nil {
			errdlg("Could not load biome infos %s:\n%s", path, err.Error())
			return
		}
		defer f.Close()

		infos, err := ReadBiomeInfos(f)
		if err != nil {
			errdlg("Could not load biome infos %s:\n%s", path, err.Error())
			return
		}

		ed.biomes = infos
		// TODO: Update view
	}
}

func (ed *BiomeInfoEditor) save() {
	// TODO
}

func (ed *BiomeInfoEditor) Biomes() []BiomeInfo { return ed.biomes }
