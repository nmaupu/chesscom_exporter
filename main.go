package main

import (
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
	"github.com/nmaupu/chesscom_exporter/pkg/api/chesscom"
	"github.com/nmaupu/chesscom_exporter/pkg/model"
	mywidget "github.com/nmaupu/chesscom_exporter/pkg/ui/widget"
	"golang.design/x/clipboard"
	"image/color"
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

	playerLineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	playerSubmitBtn = new(widget.Clickable)

	archivesLoading   bool
	archiveListWidget = mywidget.NewArchiveList(theme)
	archivesBorder    = &widget.Border{
		Color:        color.NRGBA{A: 0xff},
		CornerRadius: unit.Dp(8),
		Width:        unit.Px(2),
	}

	saveToClipboardBtn          = new(widget.Clickable)
	saveToClipboardInProgress   bool
	saveToClipboardStatus       string
	saveToClipboardProgress     float32
	saveToClipboardCancelChan   = make(chan bool, 1)
	saveToClipboardProgressChan = make(chan float32)
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
	playerLineEditor.Focus()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent: // repaint
				gtx := layout.NewContext(&ops, e)

				if playerSubmitBtn.Clicked() {
					go func() {
						archivesLoading = true
						defer func() {
							archivesLoading = false
						}()

						archiveListWidget.ResetList()

						var err error
						player := strings.Trim(playerLineEditor.Text(), " ")
						if player == "" {
							return
						}
						archives, err := chesscom.GetAllPlayerArchives(player)
						if err != nil {
							log.Printf("unable to get archives for %s, err=%v", player, err)
							return
						}

						archiveListWidget.AddRows(archives)
						w.Invalidate()
					}()
				}

				if saveToClipboardBtn.Clicked() {
					routine := func() { // Go routine to get all asked archives
						saveToClipboardInProgress = true
						saveToClipboardStatus = "In progress"
						defer func() {
							saveToClipboardInProgress = false
						}()

						// Resetting progress
						saveToClipboardProgressChan <- 0

						selectedArchives := archiveListWidget.GetSelectedArchives()
						total := len(selectedArchives.Archives)
						ch := make(chan struct {
							idx     int
							archive model.ChesscomArchive
						}, total)

						for i, archive := range selectedArchives.Archives {
							ch <- struct {
								idx     int
								archive model.ChesscomArchive
							}{idx: i, archive: archive}
						}
						close(ch)

						builder := strings.Builder{}
					loop:
						for {
							select {
							case <-saveToClipboardCancelChan:
								// Resetting progress
								saveToClipboardProgressChan <- 0
								saveToClipboardStatus = "Aborted."
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

								// Writing pgn games to a buffer
								for _, game := range res.Games {
									builder.Write([]byte(game.PGN + "\n"))
								}

								progress := float32(e.idx) / float32(total)
								saveToClipboardProgressChan <- progress
							}
						}

						clipboard.Write(clipboard.FmtText, []byte(builder.String()))
						saveToClipboardProgressChan <- 1
						saveToClipboardStatus = "Success !"
					}

					if saveToClipboardInProgress { // already in progress, canceling
						saveToClipboardCancelChan <- true
					} else {
						go routine()
					}
				}

				kitchen(gtx, theme)
				e.Frame(gtx.Ops)
			}
		case p := <-saveToClipboardProgressChan:
			if p <= 1 {
				saveToClipboardProgress = p
				log.Printf("Progress=%f", saveToClipboardProgress)
				w.Invalidate()
			}
		}
	}
}

func kitchen(gtx C, th *material.Theme) D {
	nickEditWidget := func(gtx C) D {
		return layout.Flex{
			Axis:      layout.Horizontal,
			WeightSum: 100,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Flexed(78, func(gtx C) D {
				return playerEditorLayout(gtx, th)
			}),
			layout.Flexed(22, func(gtx C) D {
				margins := layout.Inset{
					Left: unit.Dp(5),
				}

				if archivesLoading {
					gtx = gtx.Disabled()
				}

				btn := material.Button(th, playerSubmitBtn, "Get archives")
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

	saveToClipboardWidget := func(gtx C) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(material.ProgressBar(th, saveToClipboardProgress).Layout),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Alignment: layout.Middle,
					Axis:      layout.Horizontal,
					Spacing:   layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if saveToClipboardStatus != "" {
							lbl := material.Label(th, unit.Dp(16), saveToClipboardStatus)
							lbl.Font.Style = text.Italic
							return lbl.Layout(gtx)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						txt := "Export to file"
						btn := material.Button(th, saveToClipboardBtn, txt)
						gtx = gtx.Disabled()
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						txt := "Export to clipboard"
						if saveToClipboardInProgress {
							txt = "Cancel export to clipboard"
						}
						btn := material.Button(th, saveToClipboardBtn, txt)
						return btn.Layout(gtx)
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
				layout.Rigid(nickEditWidget),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Flexed(1, archivesListWidget),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(saveToClipboardWidget),
			)
		})
}

func playerEditorLayout(gtx C, th *material.Theme) D {
	e := material.Editor(th, playerLineEditor, "Enter player's name")
	e.Font.Style = text.Italic
	border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
	})
}
