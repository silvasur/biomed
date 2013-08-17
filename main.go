package main

import (
	"fmt"
	"github.com/kch42/gomcmap/mcmap"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
)

type GUI struct {
	window    *gtk.Window
	statusbar *gtk.Statusbar

	showbiomes *gtk.CheckButton
	fixsnowice *gtk.CheckButton

	statusContext uint
	lastStatus    string

	mapw *MapWidget
}

func (g *GUI) openWorldDlg() {
	dlg := gtk.NewFileChooserDialog("Open World (region directory)", g.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER, "Open Region dir", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
	if dlg.Run() == gtk.RESPONSE_OK {
		g.openWorld(dlg.GetFilename())
	}
	dlg.Destroy()
}

func (g *GUI) openWorld(path string) {
	region, err := mcmap.OpenRegion(path, false)
	if err != nil {
		dlg := gtk.NewMessageDialog(g.window, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "Could not load world %s:\n%s", path, err.Error())
		dlg.Run()
		dlg.Destroy()
	}

	go g.mapw.setRegion(region)
}

func (g *GUI) aboutDlg() {
	dlg := gtk.NewAboutDialog()
	dlg.SetName("biome-editor")
	dlg.SetVersion("α")
	dlg.SetCopyright("© 2013 by Kevin Chabowski")
	dlg.SetAuthors([]string{"Kevin Chabowski <kevin@kch42.de>"})
	dlg.SetLicense(`Copyright (c) 2013 Kevin Chabowski

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

func (g *GUI) mkMenuBar() *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	fileMenu := gtk.NewMenu()

	open := gtk.NewMenuItemWithLabel("Open")
	open.Connect("activate", g.openWorldDlg)
	fileMenu.Append(open)

	if quickopen, ok := g.mkQuickOpen(); ok {
		quickopenItem := gtk.NewMenuItemWithLabel("Open Map")
		quickopenItem.SetSubmenu(quickopen)
		fileMenu.Append(quickopenItem)
	}

	save := gtk.NewMenuItemWithLabel("Save")
	save.Connect("activate", g.save)
	fileMenu.Append(save)

	quit := gtk.NewMenuItemWithLabel("Quit")
	quit.Connect("activate", g.exitApp)
	fileMenu.Append(quit)

	fileMenuItem := gtk.NewMenuItemWithLabel("File")
	fileMenuItem.SetSubmenu(fileMenu)
	menubar.Append(fileMenuItem)

	/*editMenu := gtk.NewMenu()

	undo := gtk.NewMenuItemWithLabel("Undo")
	undo.Connect("activate", g.undo)
	editMenu.Append(undo)

	editMenuItem := gtk.NewMenuItemWithLabel("Edit")
	editMenuItem.SetSubmenu(editMenu)
	menubar.Append(editMenuItem)*/

	helpMenu := gtk.NewMenu()

	about := gtk.NewMenuItemWithLabel("About")
	about.Connect("activate", g.aboutDlg)
	helpMenu.Append(about)

	helpMenuItem := gtk.NewMenuItemWithLabel("Help")
	helpMenuItem.SetSubmenu(helpMenu)
	menubar.Append(helpMenuItem)

	return menubar
}

func (g *GUI) save() {
	g.mapw.Save()
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
	vbox := gtk.NewVBox(false, 0)

	vbox.PackStart(labelCustomFont("Tools", "Sans Bold 14"), false, false, 3)

	g.showbiomes = gtk.NewCheckButtonWithLabel("Show Biomes")
	g.showbiomes.SetActive(true)
	g.showbiomes.Connect("toggled", g.showbiomesToggled)
	vbox.PackStart(g.showbiomes, false, false, 3)

	g.fixsnowice = gtk.NewCheckButtonWithLabel("Fix Snow/Ice")
	g.fixsnowice.SetTooltipText("Add Snow/Ice for Taiga/Ice Plains. Remove Snow/Ice for other biomes.")
	g.fixsnowice.Connect("toggled", g.fixsnowiceToggled)
	vbox.PackStart(g.fixsnowice, false, false, 3)

	fill := gtk.NewRadioButtonWithLabel(nil, "Fill")
	fill.SetActive(true)
	fill.Connect("toggled", g.mkUpdateToolFx(fill, NewFillTool()))

	draw := gtk.NewRadioButtonWithLabel(fill.GetGroup(), "Draw")
	drawRadius := gtk.NewSpinButtonWithRange(1, 20, 1)
	drawHBox := gtk.NewHBox(false, 0)
	drawHBox.PackStart(draw, true, true, 0)
	drawHBox.PackStart(gtk.NewLabel("Radius:"), false, false, 3)
	drawHBox.PackEnd(drawRadius, false, false, 3)
	draw.Connect("toggled", g.mkUpdateToolFx(draw, NewDrawTool(func() int { return int(drawRadius.GetValue()) })))

	vbox.PackStart(fill, false, false, 3)
	vbox.PackStart(drawHBox, false, false, 3)

	vbox.PackStart(gtk.NewHSeparator(), false, false, 3)
	vbox.PackStart(labelCustomFont("Biomes", "Sans Bold 14"), false, false, 3)

	var grp *glib.SList
	for _, bio := range bioList {
		biohbox := gtk.NewHBox(false, 0)
		cbox := colorBox(bioColors[bio])
		cbox.SetSizeRequest(20, 20)
		biohbox.PackStart(cbox, false, false, 3)
		rbutton := gtk.NewRadioButtonWithLabel(grp, bio.String())
		grp = rbutton.GetGroup()
		rbutton.Connect("toggled", g.mkUpdateBiomeFx(rbutton, bio))
		biohbox.PackEnd(rbutton, true, true, 3)
		vbox.PackStart(biohbox, false, false, 3)
	}

	scrolled := gtk.NewScrolledWindow(nil, nil)
	scrolled.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scrolled.AddWithViewPort(vbox)
	return scrolled
}

func (g *GUI) Init() {
	g.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	g.window.SetTitle("Biome Editor")

	menubar := g.mkMenuBar()
	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(menubar, false, false, 0)

	hbox := gtk.NewHBox(false, 0)

	g.mapw = NewMapWidget(g.reportError, g.updateInfo, g.setBusy)
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

func (g *GUI) updateInfo(x, z int, bio mcmap.Biome) {
	g.lastStatus = fmt.Sprintf("X:%d, Z:%d, Biome:%s", x, z, bio)
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

func (g *GUI) fixsnowiceToggled() {
	g.mapw.SetFixSnowIce(g.fixsnowice.GetActive())
}

/*func (g *GUI) undo() {
	fmt.Println("Undo")
}*/

func (g *GUI) Show() {
	g.window.ShowAll()
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
