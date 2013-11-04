package main

import (
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"strconv"
)

type biomeEditFrame struct {
	*gtk.Frame
	applyBtn                          *gtk.Button
	idInput, snowLineInput, nameInput *gtk.Entry
	colorInput                        *gtk.ColorButton
	bList                             *biomeList
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

	frm.applyBtn.Connect("clicked", frm.doApply)

	return frm
}

func (frm *biomeEditFrame) setBiomeInfo(info BiomeInfo) {
	frm.colorInput.SetColor(gdk.NewColor(info.Color))
	frm.idInput.SetText(strconv.FormatInt(int64(info.ID), 10))
	frm.snowLineInput.SetText(strconv.FormatInt(int64(info.SnowLine), 10))
	frm.nameInput.SetText(info.Name)
}

func (frm *biomeEditFrame) doApply() {
	biome, ok := frm.getBiomeInfo()
	if !ok {
		return
	}

	frm.bList.setCurrentBiome(biome)
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
		Color:    fmt.Sprintf("#%02x%02x%02x", col.Red()>>8, col.Green()>>8, col.Blue()>>8),
	}, true
}

func (frm *biomeEditFrame) checkOK() bool {
	_, ok := frm.getBiomeInfo()
	return ok
}

func (frm *biomeEditFrame) unlockApply() {
	frm.applyBtn.SetSensitive(frm.checkOK())
}

type biomeList struct {
	*gtk.HBox
	treeview *gtk.TreeView
	lStore   *gtk.ListStore
	biomes   []BiomeInfo
	editfrm  *biomeEditFrame
}

func newBiomeList() *biomeList {
	bl := &biomeList{
		HBox:     gtk.NewHBox(false, 0),
		treeview: gtk.NewTreeView(),
		lStore:   gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING),
	}

	scroll := gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scroll.Add(bl.treeview)
	bl.PackStart(scroll, true, true, 3)

	bl.treeview.SetModel(bl.lStore)
	bl.treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Color", gtk.NewCellRendererText(), "background", 0))
	bl.treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("ID", gtk.NewCellRendererText(), "text", 1))
	bl.treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Snowline", gtk.NewCellRendererText(), "text", 2))
	bl.treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Name", gtk.NewCellRendererText(), "text", 3))

	bl.treeview.GetSelection().SetMode(gtk.SELECTION_SINGLE)
	bl.treeview.Connect("cursor-changed", bl.onCursorChanged)

	vbox := gtk.NewVBox(false, 0)

	addBtn := gtk.NewButton()
	addBtn.Add(gtk.NewImageFromStock(gtk.STOCK_ADD, gtk.ICON_SIZE_SMALL_TOOLBAR))
	delBtn := gtk.NewButton()
	delBtn.Add(gtk.NewImageFromStock(gtk.STOCK_DELETE, gtk.ICON_SIZE_SMALL_TOOLBAR))
	upBtn := gtk.NewButton()
	upBtn.Add(gtk.NewImageFromStock(gtk.STOCK_GO_UP, gtk.ICON_SIZE_SMALL_TOOLBAR))
	downBtn := gtk.NewButton()
	downBtn.Add(gtk.NewImageFromStock(gtk.STOCK_GO_DOWN, gtk.ICON_SIZE_SMALL_TOOLBAR))

	addBtn.Connect("clicked", bl.onAdd)
	delBtn.Connect("clicked", bl.onDel)
	upBtn.Connect("clicked", bl.onUp)
	downBtn.Connect("clicked", bl.onDown)

	vbox.PackStart(addBtn, false, false, 3)
	vbox.PackStart(delBtn, false, false, 3)
	vbox.PackStart(upBtn, false, false, 3)
	vbox.PackStart(downBtn, false, false, 3)

	bl.PackStart(vbox, false, false, 0)

	return bl
}

func (bl *biomeList) setBiome(iter *gtk.TreeIter, biome BiomeInfo) {
	bl.lStore.Set(iter, biome.Color, strconv.FormatInt(int64(biome.ID), 10), strconv.FormatInt(int64(biome.SnowLine), 10), biome.Name)
}

func (bl *biomeList) setCurrentBiome(biome BiomeInfo) {
	idx, iter := bl.treeviewIdx()
	if idx < 0 {
		return
	}
	bl.biomes[idx] = biome
	bl.setBiome(iter, biome)
}

func (bl *biomeList) SetBiomes(biomes []BiomeInfo) {
	bl.biomes = biomes

	bl.lStore.Clear()
	var iter gtk.TreeIter
	for _, bio := range biomes {
		bl.lStore.Append(&iter)
		bl.setBiome(&iter, bio)
	}
}

func (bl *biomeList) Biomes() []BiomeInfo { return bl.biomes }

func (bl *biomeList) treeviewIdx() (int, *gtk.TreeIter) {
	var path *gtk.TreePath
	var column *gtk.TreeViewColumn
	bl.treeview.GetCursor(&path, &column)

	idxs := path.GetIndices()
	if len(idxs) != 1 {
		return -1, nil
	}
	var iter gtk.TreeIter
	bl.lStore.GetIter(&iter, path)

	return idxs[0], &iter
}

func (bl *biomeList) onCursorChanged() {
	idx, _ := bl.treeviewIdx()
	if idx < 0 {
		return
	}

	bl.editfrm.setBiomeInfo(bl.biomes[idx])
}

func (bl *biomeList) onAdd() {
	bio := BiomeInfo{
		Color:    "#000000",
		ID:       0,
		SnowLine: 255,
		Name:     "(new)",
	}
	bl.biomes = append(bl.biomes, bio)

	var iter gtk.TreeIter
	bl.lStore.Append(&iter)
	bl.setBiome(&iter, bio)
	path := gtk.NewTreePath()
	path.AppendIndex(len(bl.biomes) - 1)
	bl.treeview.SetCursor(path, nil, false)
}

func (bl *biomeList) onDel() {
	idx, iter := bl.treeviewIdx()
	if idx < 0 {
		return
	}

	copy(bl.biomes[idx:], bl.biomes[idx+1:])
	bl.biomes = bl.biomes[:len(bl.biomes)-1]

	bl.lStore.Remove(iter)
}
func (bl *biomeList) onUp()   {} // TODO
func (bl *biomeList) onDown() {} // TODO

func connectBiomeListEditFrame(bl *biomeList, frm *biomeEditFrame) {
	bl.editfrm = frm
	frm.bList = bl
}

type BiomeInfoEditor struct {
	*gtk.Dialog
	biolist *biomeList
}

func NewBiomeInfoEditor(biomes []BiomeInfo) *BiomeInfoEditor {
	ed := &BiomeInfoEditor{
		Dialog:  gtk.NewDialog(),
		biolist: newBiomeList(),
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

	ed.biolist.SetBiomes(biomes)
	vbox.PackStart(ed.biolist, true, true, 3)

	editFrame := newBiomeEditFrame()
	connectBiomeListEditFrame(ed.biolist, editFrame)
	vbox.PackStart(editFrame, false, false, 3)

	ed.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	ed.AddButton("OK", gtk.RESPONSE_OK)
	ed.ShowAll()
	return ed
}

func (ed *BiomeInfoEditor) reset() {
	ed.biolist.SetBiomes(ReadDefaultBiomes())
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
	dlg := gtk.NewFileChooserDialog("Load", nil, gtk.FILE_CHOOSER_ACTION_OPEN, "OK", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
	defer dlg.Destroy()
askFilename:
	if dlg.Run() == gtk.RESPONSE_OK {
		path := dlg.GetFilename()

		f, err := os.Open(path)
		if err != nil {
			errdlg("Could not load biome infos %s:\n%s", path, err.Error())
			goto askFilename
		}
		defer f.Close()

		infos, err := ReadBiomeInfos(f)
		if err != nil {
			errdlg("Could not load biome infos %s:\n%s", path, err.Error())
			goto askFilename
		}

		ed.biolist.SetBiomes(infos)
	}
}

func (ed *BiomeInfoEditor) save() {
	dlg := gtk.NewFileChooserDialog("Save", nil, gtk.FILE_CHOOSER_ACTION_SAVE, "OK", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
	defer dlg.Destroy()
askFilename:
	if dlg.Run() == gtk.RESPONSE_OK {
		path := dlg.GetFilename()

		if _, err := os.Stat(path); err == nil {
			qdlg := gtk.NewMessageDialog(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "File %s already exists. Overwrite?", path)
			resp := qdlg.Run()
			qdlg.Destroy()

			if resp != gtk.RESPONSE_YES {
				goto askFilename
			}
		}

		f, err := os.Create(path)
		if err != nil {
			errdlg("Could not save biome infos %s:\n%s", path, err.Error())
			goto askFilename
		}
		defer f.Close()

		if err := WriteBiomeInfos(ed.biolist.Biomes(), f); err != nil {
			errdlg("Could not save biome infos %s:\n%s", path, err.Error())
			goto askFilename
		}
	}
}

func (ed *BiomeInfoEditor) Biomes() []BiomeInfo { return ed.biolist.Biomes() }
