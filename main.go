package main

import (
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type config struct{
	EditWidget *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cfg config

func main(){
	// create a fyne app
	a:=app.New()
	// create a window for the app
	win:=a.NewWindow("Markdown")
	// get the user interface
	edit,preview:= cfg.makeUI()

	cfg.createMenuItems(win)
	// set the content of the window
	win.SetContent(container.NewHSplit(edit,preview))
	// show window and run app
	win.Resize(fyne.Size{Width:800, Height:800})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry,*widget.RichText){
	edit:=widget.NewMultiLineEntry()
	preview:=widget.NewRichTextFromMarkdown("")
	app.EditWidget=edit
	app.PreviewWidget=preview

	edit.OnChanged = preview.ParseMarkdown
	return edit,preview
}

func (app *config) createMenuItems(win fyne.Window){
	openMenuItem:=fyne.NewMenuItem("Open...",app.openFunc(win))
	saveMenuItem:=fyne.NewMenuItem("Save",func(){})
	app.SaveMenuItem= saveMenuItem
	app.SaveMenuItem.Disabled=true
	saveAsMenuItem:=fyne.NewMenuItem("Save as...",app.saveAsFunc(win))

	fileMenu:=fyne.NewMenu("File",openMenuItem,saveMenuItem,saveAsMenuItem)

	menu:=fyne.NewMainMenu(fileMenu)
	win.SetMainMenu(menu)
}
func (app *config) openFunc(win fyne.Window) func(){
	return func(){
		openDialog:=dialog.NewFileOpen(func (read fyne.URIReadCloser, err error){
			if err!=nil{
				dialog.ShowError(err, win)
				return
			}

			if read==nil{
				return
			}

			defer read.Close()
			data,err:=ioutil.ReadAll(read)
			if err!=nil{
				dialog.ShowError(err,win)
				return
			}

			app.EditWidget.SetText(string(data))
			app.CurrentFile = read.URI()
			win.SetTitle(win.Title()+" - " + read.URI().Name())
			app.SaveMenuItem.Disabled=false
		}, win)
		openDialog.Show()
	}
}

func (app *config) saveAsFunc(win fyne.Window) func(){
	return func(){
		saveDialog := dialog.NewFileSave(func (write fyne.URIWriteCloser, err error){
			if err!=nil{
				dialog.ShowError(err,win)
				return
			}

			if write==nil{
				// user cancelled
				return
			}

			// save the file
			write.Write([]byte(app.EditWidget.Text))
			log.Println(string([]byte(app.EditWidget.Text)))
			app.CurrentFile=write.URI()
			log.Println("write.URI()=",app.CurrentFile)
			defer write.Close()
			win.SetTitle(win.Title()+" - " + write.URI().Name())
			log.Println( "write.URI().Name()=",write.URI().Name())
			app.SaveMenuItem.Disabled=false
		},win)
		saveDialog.Show()
	}
}