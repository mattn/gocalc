package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

var (
	display   *gtk.GtkEntry // where values are displayed
	inputMode = true
	nums      = "789/456x123-0.=+"
	operators = "/x-+="
)

// End the program
func Quit() {
	gtk.MainQuit()
}

// Action to be performed by each button, returns a handler function
func Input(b *gtk.GtkButton) func() {
	return func() {
		if strings.Index(operators, b.GetLabel()) != -1 {
			val, _ := strconv.ParseFloat(display.GetText(), 32)
			Calculation(val, b.GetLabel())
			display.SetText(GetResult())
			inputMode = false
		} else {
			if inputMode {
				display.SetText(display.GetText() + b.GetLabel())
			} else {
				display.SetText(b.GetLabel())
				inputMode = true
				if result.operator == "=" {
					Reset()
				}
			}
		}
	}
}

func main() {
	gtk.Init(&os.Args)
	display = gtk.Entry()
	display.SetAlignment(1.0)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("Simple Go Calculator")
	window.Connect("destroy", Quit)

	// Vertical box containing all components
	vbox := gtk.VBox(false, 1)

	// Menu bar
	menubar := gtk.MenuBar()
	vbox.PackStart(menubar, false, false, 0)

	// Add calculator display to vertical box
	display.SetCanFocus(false) // disable focus on calcuator display
	vbox.Add(display)

	// Menu items
	filemenu := gtk.MenuItemWithMnemonic("_File")
	menubar.Append(filemenu)
	filesubmenu := gtk.Menu()
	filemenu.SetSubmenu(filesubmenu)

	aboutmenuitem := gtk.MenuItemWithMnemonic("_About")
	aboutmenuitem.Connect("activate", func() {
		messagedialog := gtk.MessageDialog(
			window.GetTopLevelAsWindow(),
			gtk.GTK_DIALOG_MODAL,
			gtk.GTK_MESSAGE_INFO,
			gtk.GTK_BUTTONS_OK,
			"Simple Go Calculator")
		messagedialog.Response(func() {})
		messagedialog.Run()
		messagedialog.Destroy()
	},
		nil)
	filesubmenu.Append(aboutmenuitem)

	resetmenuitem := gtk.MenuItemWithMnemonic("_Reset")
	resetmenuitem.Connect("activate", func() { Reset(); display.SetText("0") })
	filesubmenu.Append(resetmenuitem)

	exitmenuitem := gtk.MenuItemWithMnemonic("E_xit")
	exitmenuitem.Connect("activate", Quit)
	filesubmenu.Append(exitmenuitem)

	// Vertical box containing all buttons
	buttons := gtk.VBox(false, 5)

	bmap := map[string]*gtk.GtkButton{}

	for i := 0; i < 4; i++ {
		hbox := gtk.HBox(false, 5) // a horizontal box for each 4 buttons
		for j := 0; j < 4; j++ {
			b := gtk.ButtonWithLabel(string(nums[i*4+j]))
			b.Clicked(Input(b)) //add click event
			hbox.Add(b)
			bmap[string(nums[i*4+j])] = b
		}
		buttons.Add(hbox) // add horizonatal box to vertical buttons' box
	}

	vbox.Add(buttons)

	window.Connect("key-press-event", func(ctx *glib.CallbackContext) bool {
		arg := ctx.Args(0)
		kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		c := (string(uint8(kev.Keyval % 0xff)))
		if kev.Keyval == gdk.GDK_KEY_Return {
			c = "="
			return true
		}
		if b, ok := bmap[c]; ok {
			Input(b)()
			b.GrabFocus()
		} else if kev.Keyval == gdk.GDK_KEY_Delete {
			Reset()
			display.SetText("0")
			return true
		}
		return false
	})

	window.Add(vbox)
	window.SetSizeRequest(250, 250)
	window.ShowAll()
	gtk.Main()
}
