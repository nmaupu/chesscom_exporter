package main

import (
	"bytes"
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	"github.com/nmaupu/chesscom_exporter/pkg/api/chesscom"
	"github.com/nmaupu/chesscom_exporter/pkg/model"
	mywidget "github.com/nmaupu/chesscom_exporter/pkg/ui/widget"
	"golang.design/x/clipboard"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	theme = material.NewTheme(gofont.Collection())

	usernameLineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	usernameSubmitBtn = new(widget.Clickable)

	archivesLoading   bool
	archiveListWidget = mywidget.NewArchiveList(theme)
	archivesBorder    = &widget.Border{
		Color:        color.NRGBA{A: 0xff},
		CornerRadius: unit.Dp(8),
		Width:        unit.Px(2),
	}

	saveToFileBtn      = new(widget.Clickable)
	saveToClipboardBtn = new(widget.Clickable)
	saveCancelBtn      = new(widget.Clickable)

	saveInProgress         bool
	saveStatus             string
	saveProgress           float32
	saveProgressCancelChan = make(chan bool, 1)
	saveProgressChan       = make(chan float32)
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Chesscom exporter"),
			app.Size(unit.Dp(800), unit.Dp(600)),
			app.MinSize(unit.Dp(400), unit.Dp(400)),
		)

		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)

	}()
	app.Main()
}

func draw(w *app.Window) error {
	var ops op.Ops

	// Focusing player's text edit by default
	usernameLineEditor.Focus()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent: // repaint
				gtx := layout.NewContext(&ops, e)

				if usernameSubmitBtn.Clicked() {
					go func() {
						archivesLoading = true
						defer func() {
							archivesLoading = false
						}()

						archiveListWidget.ResetList()

						var err error
						username := strings.Trim(usernameLineEditor.Text(), " ")
						if username == "" {
							return
						}
						archives, err := chesscom.GetAllPlayerArchives(username)
						if err != nil {
							log.Printf("unable to get archives for %s, err=%v", username, err)
							return
						}

						archiveListWidget.AddRows(archives)
						w.Invalidate()
					}()
				}

				if saveToClipboardBtn.Clicked() {
					go func() { // Go routine to get all checked archives
						loadAndSaveSelectedArchivesInMemory(func(r io.Reader) error {
							data, err := ioutil.ReadAll(r)
							if err != nil {
								return err
							}
							clipboard.Write(clipboard.FmtText, data)
							return nil
						})
					}()

				}

				if saveToFileBtn.Clicked() {
					username := strings.Trim(usernameLineEditor.Text(), " ")
					fileWriter, err := explorer.WriteFile(fmt.Sprintf("chesscom-export-%s.pgn", username))
					if err != nil {
						saveStatus = "Not supported, sorry :/"
					} else {
						go func() {
							loadAndSaveSelectedArchivesInMemory(func(r io.Reader) error {
								defer fileWriter.Close()
								_, err := io.Copy(fileWriter, r)
								return err
							})
						}()
					}
				}

				if saveCancelBtn.Clicked() {
					if saveInProgress { // button is normally disabled when not in progress though
						saveProgressCancelChan <- true
					}
				}

				kitchen(gtx, theme)
				e.Frame(gtx.Ops)
			}
		case p := <-saveProgressChan:
			if p <= 1 {
				saveProgress = p
				log.Printf("Progress=%f", saveProgress)
				w.Invalidate()
			}
		}
	}
}

func kitchen(gtx C, th *material.Theme) D {
	usernameEditWidget := func(gtx C) D {
		return layout.Flex{
			Axis:      layout.Horizontal,
			WeightSum: 100,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Flexed(78, func(gtx C) D {
				return usernameEditorLayout(gtx, th)
			}),
			layout.Flexed(22, func(gtx C) D {
				margins := layout.Inset{
					Left: unit.Dp(5),
				}

				if archivesLoading {
					gtx = gtx.Disabled()
				}

				btn := material.Button(th, usernameSubmitBtn, "Get archives")
				return margins.Layout(gtx, btn.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if archivesLoading {
					return material.Loader(th).Layout(gtx)
				}
				return D{}
			}),
		)
	}

	archivesListWidget := func(gtx C) D {
		return archivesBorder.Layout(gtx, func(gtx C) D {
			insets := layout.UniformInset(unit.Dp(10))
			return insets.Layout(gtx, func(gtx C) D {

				if archiveListWidget.IsNil() || archiveListWidget.Size() == 0 {
					txt := "Enter a player's name to display archives"
					if !archiveListWidget.IsNil() {
						txt = fmt.Sprintf("No archives available for the selected user")
					}

					if archivesLoading {
						txt = "Loading archives..."
					}

					lbl := material.Label(th, unit.Dp(16), "")
					lbl.Text = txt
					lbl.Alignment = text.Middle
					lbl.MaxLines = 1
					lbl.Font.Style = text.Italic
					return lbl.Layout(gtx)
				}

				return archiveListWidget.Layout(gtx)
			})
		})
	}

	saveWidget := func(gtx C) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(material.ProgressBar(th, saveProgress).Layout),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Alignment: layout.Middle,
					Axis:      layout.Horizontal,
					Spacing:   layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if saveStatus != "" {
							lbl := material.Label(th, unit.Dp(16), saveStatus)
							lbl.Font.Style = text.Italic
							return lbl.Layout(gtx)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx C) D {
						if saveInProgress || !archiveListWidget.AtLeastOneSelected() {
							gtx = gtx.Disabled()
						}
						return material.Button(th, saveToFileBtn, "Export to a file").Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx C) D {
						if saveInProgress || !archiveListWidget.AtLeastOneSelected() {
							gtx = gtx.Disabled()
						}
						return material.Button(th, saveToClipboardBtn, "Export to clipboard").Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx C) D {
						if !saveInProgress {
							gtx = gtx.Disabled()
						}
						return material.Button(th, saveCancelBtn, "Cancel").Layout(gtx)
					}),
				)
			}),
		)
	}

	outerInset := layout.UniformInset(unit.Dp(5))
	return outerInset.Layout(gtx,
		func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(usernameEditWidget),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Flexed(1, archivesListWidget),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(saveWidget),
			)
		})
}

func usernameEditorLayout(gtx C, th *material.Theme) D {
	e := material.Editor(th, usernameLineEditor, "Enter player's name")
	e.Font.Style = text.Italic
	border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
	})
}

func loadAndSaveSelectedArchivesInMemory(saver func(o io.Reader) error) {
	saveInProgress = true
	defer func() {
		saveInProgress = false
	}()
	saveStatus = "In progress"
	saveProgressChan <- 0 // resetting progress

	selectedArchives := archiveListWidget.GetSelectedArchives()
	total := len(selectedArchives.Archives)
	ch := make(chan struct {
		idx     int
		archive model.ChesscomArchive
	}, total)

	// Pushing jobs in channel
	for i, archive := range selectedArchives.Archives {
		ch <- struct {
			idx     int
			archive model.ChesscomArchive
		}{idx: i, archive: archive}
	}
	close(ch)

	buf := bytes.Buffer{}
loop:
	for {
		select {
		case <-saveProgressCancelChan:
			// Resetting progress
			saveProgressChan <- 0
			saveStatus = "Aborted."
			return
		case e, ok := <-ch:
			if !ok { // no more jobs in channel
				break loop
			}
			res, err := chesscom.GetPlayerMonthlyArchivesByURL(e.archive.GetURL())
			if err != nil {
				log.Printf("an error occurred trying to get archive %s", e.archive.GetURL())
				continue
			}

			// Writing pgn games to the buffer
			for _, game := range res.Games {
				buf.Write([]byte(game.PGN + "\n"))
			}

			// update progress
			saveProgressChan <- float32(e.idx) / float32(total)
		}
	}

	if err := saver(&buf); err != nil {
		saveStatus = fmt.Sprintf("Error: %v", err)
		saveProgressChan <- 0 // resetting progress
		return
	}
	saveProgressChan <- 1
	saveStatus = "Success !"
}
