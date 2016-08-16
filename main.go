package main

import (
	"fmt"
	"github.com/silvasur/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"path"
)

type GUI struct {
	window    *gtk.Window
	statusbar *gtk.Statusbar

	accel *gtk.AccelGroup

	showbiomes *gtk.CheckButton
	fixSnowIce *gtk.CheckButton

	menuitemSave *gtk.ImageMenuItem

	statusContext uint
	lastStatus    string

	biomes      []BiomeInfo
	bioVBox     *gtk.VBox
	bioVBoxWrap *gtk.VBox

	mapw *MapWidget
}

func (g *GUI) openWorldDlg() {
	dlg := gtk.NewFileChooserDialog("Open World (level.dat)", g.window, gtk.FILE_CHOOSER_ACTION_OPEN, "Open", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
	filter := gtk.NewFileFilter()
	filter.AddPattern("level.dat")
	filter.SetName("level.dat")
	dlg.AddFilter(filter)
	if dlg.Run() == gtk.RESPONSE_OK {
		g.openWorld(dlg.GetFilename())
	}
	dlg.Destroy()
}

func readWorld(p string) (*mcmap.Region, int, int, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, 0, 0, err
	}
	defer f.Close()

	x, z, err := getMapCenter(f)
	if err != nil {
		return nil, 0, 0, err
	}

	dir, _ := path.Split(p)
	region, err := mcmap.OpenRegion(path.Join(dir, "region"), false)
	return region, x, z, err
}

func (g *GUI) openWorld(p string) {
	region, centerX, centerZ, err := readWorld(p)
	if err != nil {
		dlg := gtk.NewMessageDialog(g.window, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "Could not load world %s:\n%s", p, err.Error())
		dlg.Run()
		dlg.Destroy()
	}

	g.menuitemSave.SetSensitive(true)

	g.mapw.SetCenter(centerX, centerZ)
	g.mapw.SetRegion(region)
}

func (g *GUI) aboutDlg() {
	dlg := gtk.NewAboutDialog()
	dlg.SetProgramName("biomed")
	dlg.SetComments("A Minecraft Biome Editor")
	dlg.SetVersion("β")
	dlg.SetCopyright("© 2013 by Laria Carolin Chabowski")
	dlg.SetAuthors([]string{"Laria Carolin Chabowski <laria@laria.me>"})
	dlg.SetLicense(`Copyright (c) 2013 Laria Carolin Chabowski

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`)
	dlg.Run()
	dlg.Destroy()
}

func (g *GUI) biomeEditor() {
	ed := NewBiomeInfoEditor(g.biomes)
	rv := ed.Run()
	ed.Destroy()
	if rv == gtk.RESPONSE_OK {
		g.biomes = ed.Biomes()
		g.updateBiomeInfo()
	}
}

func (g *GUI) mkMenuBar() *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	fileMenu := gtk.NewMenu()

	open := gtk.NewImageMenuItemFromStock(gtk.STOCK_OPEN, g.accel)
	open.Connect("activate", g.openWorldDlg)
	fileMenu.Append(open)

	if quickopen, ok := g.mkQuickOpen(); ok {
		quickopenItem := gtk.NewMenuItemWithLabel("Open Map")
		quickopenItem.SetSubmenu(quickopen)
		fileMenu.Append(quickopenItem)
	}

	g.menuitemSave = gtk.NewImageMenuItemFromStock(gtk.STOCK_SAVE, g.accel)
	g.menuitemSave.Connect("activate", g.save)
	g.menuitemSave.SetSensitive(false)
	fileMenu.Append(g.menuitemSave)

	quit := gtk.NewImageMenuItemFromStock(gtk.STOCK_QUIT, g.accel)
	quit.Connect("activate", g.exitApp)
	fileMenu.Append(quit)

	fileMenuItem := gtk.NewMenuItemWithLabel("File")
	fileMenuItem.SetSubmenu(fileMenu)
	menubar.Append(fileMenuItem)

	helpMenu := gtk.NewMenu()

	controls := gtk.NewMenuItemWithLabel("Controls")
	controls.Connect("activate", func() {
		dlg := gtk.NewMessageDialog(g.window, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "Click to use selected tool.\nMiddle mouse button to move around.")
		dlg.Run()
		dlg.Destroy()
	})
	helpMenu.Append(controls)

	about := gtk.NewImageMenuItemFromStock(gtk.STOCK_ABOUT, g.accel)
	about.Connect("activate", g.aboutDlg)
	helpMenu.Append(about)

	helpMenuItem := gtk.NewMenuItemWithLabel("Help")
	helpMenuItem.SetSubmenu(helpMenu)
	menubar.Append(helpMenuItem)

	return menubar
}

func (g *GUI) betaWarning() {
	dlg := gtk.NewMessageDialog(g.window, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "This software is currently in beta.\nAlthough everything seems to work, you should make a backup of your maps, just in case!")
	dlg.Run()
	dlg.Destroy()
}

func (g *GUI) save() {
	g.setBusy(true)
	g.mapw.Save()
	g.setBusy(false)

	dlg := gtk.NewMessageDialog(g.window, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "Map saved!")
	dlg.Run()
	dlg.Destroy()
}

func (g *GUI) mkQuickOpen() (*gtk.Menu, bool) {
	maps := allMaps()
	if (maps == nil) || (len(maps) == 0) {
		return nil, false
	}

	menu := gtk.NewMenu()
	for name, p := range maps {
		mitem := gtk.NewMenuItemWithLabel(name)
		p2 := p
		mitem.Connect("activate", func() { g.openWorld(p2) })
		menu.Append(mitem)
	}

	return menu, true
}

func labelCustomFont(text, font string) *gtk.Label {
	label := gtk.NewLabel(text)
	label.ModifyFontEasy(font)
	return label
}

func (g *GUI) mkSidebar() *gtk.ScrolledWindow {
	sbVBox := gtk.NewVBox(false, 0)

	sbVBox.PackStart(labelCustomFont("Tools", "Sans Bold 14"), false, false, 3)

	g.showbiomes = gtk.NewCheckButtonWithLabel("Show Biomes")
	g.showbiomes.SetActive(true)
	g.showbiomes.Connect("toggled", g.showbiomesToggled)
	sbVBox.PackStart(g.showbiomes, false, false, 3)

	g.fixSnowIce = gtk.NewCheckButtonWithLabel("Fix Snow/Ice")
	g.fixSnowIce.SetTooltipText("Add Snow/Ice for Taiga/Ice Plains. Remove Snow/Ice for other biomes.")
	g.fixSnowIce.Connect("toggled", g.fixSnowIceToggled)
	sbVBox.PackStart(g.fixSnowIce, false, false, 3)

	fill := gtk.NewRadioButtonWithLabel(nil, "Fill")
	fill.SetActive(true)
	fill.Connect("toggled", g.mkUpdateToolFx(fill, NewFillTool()))

	draw := gtk.NewRadioButtonWithLabel(fill.GetGroup(), "Draw")
	drawRadius := gtk.NewSpinButtonWithRange(1, 20, 1)
	drawHBox := gtk.NewHBox(false, 0)
	drawHBox.PackStart(draw, true, true, 0)
	drawHBox.PackStart(gtk.NewLabel("Radius:"), false, false, 3)
	drawHBox.PackEnd(drawRadius, false, false, 3)
	draw.Connect("toggled", g.mkUpdateToolFx(draw, NewDrawTool(func() int { return drawRadius.GetValueAsInt() })))

	sbVBox.PackStart(fill, false, false, 3)
	sbVBox.PackStart(drawHBox, false, false, 3)

	sbVBox.PackStart(gtk.NewHSeparator(), false, false, 3)
	bioHeaderHBox := gtk.NewHBox(false, 0)
	bioHeaderHBox.PackStart(labelCustomFont("Biomes", "Sans Bold 14"), true, false, 0)
	editBiomesBtn := gtk.NewButton()
	editBiomesBtn.Add(gtk.NewImageFromStock(gtk.STOCK_EDIT, gtk.ICON_SIZE_SMALL_TOOLBAR))
	editBiomesBtn.Connect("clicked", g.biomeEditor)
	editBiomesBtn.SetTooltipText("Configure Biomes")
	bioHeaderHBox.PackStart(editBiomesBtn, false, false, 0)
	sbVBox.PackStart(bioHeaderHBox, false, false, 3)

	g.bioVBoxWrap = gtk.NewVBox(false, 0)
	g.bioVBox = gtk.NewVBox(false, 0)
	g.bioVBoxWrap.PackStart(g.bioVBox, false, false, 0)
	sbVBox.PackStart(g.bioVBoxWrap, false, false, 3)
	g.updateBiomeInfo()

	scrolled := gtk.NewScrolledWindow(nil, nil)
	scrolled.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scrolled.AddWithViewPort(sbVBox)
	return scrolled
}

func (g *GUI) updateBiomeInfo() {
	vbox := gtk.NewVBox(false, 0)
	var grp *glib.SList

	for _, biome := range g.biomes {
		biohbox := gtk.NewHBox(false, 0)
		cbox := colorBox(gdk.NewColor(biome.Color))
		cbox.SetSizeRequest(20, 20)
		biohbox.PackStart(cbox, false, false, 3)
		rbutton := gtk.NewRadioButtonWithLabel(grp, biome.Name)
		grp = rbutton.GetGroup()
		rbutton.Connect("toggled", g.mkUpdateBiomeFx(rbutton, biome.ID))
		biohbox.PackEnd(rbutton, true, true, 3)
		vbox.PackStart(biohbox, false, false, 3)
	}

	g.bioVBoxWrap.Remove(g.bioVBox)
	g.bioVBoxWrap.PackStart(vbox, false, false, 3)
	vbox.ShowAll()
	g.bioVBox = vbox

	g.mapw.updateBioLookup(MkBiomeLookup(g.biomes))
}

func (g *GUI) Init() {
	g.biomes = ReadDefaultBiomes()

	g.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	g.window.SetTitle("biomed")

	g.accel = gtk.NewAccelGroup()
	g.window.AddAccelGroup(g.accel)

	menubar := g.mkMenuBar()
	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(menubar, false, false, 0)

	hbox := gtk.NewHBox(false, 0)

	g.mapw = NewMapWidget(GUICallbacks{g.reportError, g.updateInfo, g.setBusy}, MkBiomeLookup(g.biomes))
	hbox.PackStart(g.mapw.DArea(), true, true, 3)

	sidebar := g.mkSidebar()
	hbox.PackEnd(sidebar, false, false, 3)

	vbox.PackStart(hbox, true, true, 0)

	g.statusbar = gtk.NewStatusbar()
	g.statusContext = g.statusbar.GetContextId("mapinfo")
	vbox.PackEnd(g.statusbar, false, false, 0)

	g.window.Add(vbox)
	g.window.SetDefaultSize(800, 600)

	g.window.Connect("destroy", g.exitApp)

	g.setTool(NewFillTool())
}

func (g *GUI) setBusy(b bool) {
	g.window.SetSensitive(!b)
	g.statusbar.Pop(g.statusContext)
	if b {
		g.statusbar.Push(g.statusContext, "!!! PLEASE WAIT !!!")

	} else {
		g.statusbar.Push(g.statusContext, g.lastStatus)
	}
}

func (g *GUI) reportError(msg string) {
	dlg := gtk.NewMessageDialog(g.window, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, msg)
	dlg.Run()
	dlg.Destroy()
	os.Exit(1)
}

func (g *GUI) updateInfo(x, z int, bio mcmap.Biome, name string) {
	g.lastStatus = fmt.Sprintf("X:%d, Z:%d, Biome: %s(%d)", x, z, name, bio)
	g.statusbar.Pop(g.statusContext)
	g.statusbar.Push(g.statusContext, g.lastStatus)
}

func (g *GUI) mkUpdateToolFx(rb *gtk.RadioButton, t Tool) func() {
	return func() {
		if rb.GetActive() {
			g.setTool(t)
		}
	}
}

func (g *GUI) mkUpdateBiomeFx(rb *gtk.RadioButton, bio mcmap.Biome) func() {
	return func() {
		if rb.GetActive() {
			g.setBiome(bio)
		}
	}
}

func (g *GUI) setTool(t Tool) {
	g.mapw.SetTool(t)
}

func (g *GUI) setBiome(bio mcmap.Biome) {
	g.mapw.SetBiome(bio)
}

func (g *GUI) showbiomesToggled() {
	g.mapw.SetShowBiomes(g.showbiomes.GetActive())
}

func (g *GUI) fixSnowIceToggled() {
	g.mapw.SetFixSnowIce(g.fixSnowIce.GetActive())
}

func (g *GUI) Show() {
	g.window.ShowAll()
	g.betaWarning()
}

func (g *GUI) exitApp() {
	gtk.MainQuit()
}

func main() {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	gui := new(GUI)
	gui.Init()
	gui.Show()

	gtk.Main()
	gdk.ThreadsLeave()
}
