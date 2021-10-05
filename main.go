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
	"image/color"
	"log"
	"os"
	"strings"
	"sync"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	mutexArchives sync.Mutex
	archives      *model.ChesscomArchives
)

var (
	listMain = &widget.List{
		List: layout.List{Axis: layout.Vertical},
	}
	playerLineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	playerSubmitBtn = new(widget.Clickable)

	loadingArchives bool
	listArchives    = &widget.List{
		List: layout.List{Axis: layout.Vertical},
	}
	borderArchives = &widget.Border{
		Color:        color.NRGBA{A: 0xff},
		CornerRadius: unit.Dp(8),
		Width:        unit.Px(2),
	}
	checkboxesArchives = make(map[model.ChesscomArchive]*widget.Bool)
	clickableArchives  = make(map[model.ChesscomArchive]*widget.Clickable)
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Chesscom exporter"),
			app.Size(unit.Dp(800), unit.Dp(600)),
			app.MinSize(unit.Dp(600), unit.Dp(400)),
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

	th := material.NewTheme(gofont.Collection())

	// Focusing player's text edit by default
	playerLineEditor.Focus()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent: // repaint
				gtx := layout.NewContext(&ops, e)

				for k, v := range clickableArchives {
					if v.Clicked() {
						checkbox, ok := checkboxesArchives[k]
						if !ok {
							continue
						}
						checkbox.Value = !checkbox.Value
						w.Invalidate()
						break
					}

					if v.Hovered() {
						// TODO change background
					}
				}

				if playerSubmitBtn.Clicked() {
					go func() {
						loadingArchives = true
						defer func() {
							loadingArchives = false
						}()

						//time.Sleep(5 * time.Second)

						var err error
						player := strings.Trim(playerLineEditor.Text(), " ")
						if player == "" {
							return
						}
						mutexArchives.Lock()
						defer mutexArchives.Unlock()
						archives, err = chesscom.GetAllPlayerArchives(player)
						if err != nil {
							log.Printf("unable to get archives for %s, err=%v", player, err)
							return
						}
						w.Invalidate()
					}()
				}

				kitchen(gtx, th)
				e.Frame(gtx.Ops)
			case system.DestroyEvent:
				return e.Err
			}
		}
	}
}

func kitchen(gtx C, th *material.Theme) D {
	widgets := []layout.Widget{
		func(gtx C) D {
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
						Left:  unit.Dp(5),
						Right: unit.Dp(5),
					}

					if loadingArchives {
						gtx = gtx.Disabled()
					}

					btn := material.Button(th, playerSubmitBtn, "Get archives")
					return margins.Layout(gtx, btn.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if loadingArchives {
						return material.Loader(th).Layout(gtx)
					}
					return D{}
				}),
			)
		},
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(400))
			return borderArchives.Layout(gtx, func(gtx C) D {
				insets := layout.UniformInset(unit.Dp(10))
				return insets.Layout(gtx, func(gtx C) D {
					mutexArchives.Lock()
					defer mutexArchives.Unlock()

					lbl := material.Label(th, unit.Dp(16), "")
					txt := "Enter a player's name to display archives"
					if loadingArchives {
						txt = "Loading archives..."
					}
					if archives == nil || len(archives.Archives) == 0 {
						lbl.Text = txt
						lbl.Alignment = text.Middle
						lbl.MaxLines = 1
						lbl.Font.Style = text.Italic
						return lbl.Layout(gtx)
					}

					rowWidget := func(gtx C, i int) D {
						archive := archives.Archives[i]
						lblYear := material.Label(th, unit.Dp(16), fmt.Sprintf("%d", archive.GetYear()))
						lblMonth := material.Label(th, unit.Dp(16), fmt.Sprintf("%s", archive.GetMonthAsString()))
						if checkboxesArchives[archive] == nil {
							checkboxesArchives[archive] = &widget.Bool{}
						}

						checkbox := material.CheckBox(th, checkboxesArchives[archive], "")

						return layout.Flex{
							Axis:      layout.Horizontal,
							Spacing:   layout.SpaceEnd,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(checkbox.Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
							layout.Rigid(lblYear.Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Flexed(1, lblMonth.Layout),
						)
					}

					return material.List(th, listArchives).Layout(gtx, len(archives.Archives), func(gtx C, i int) D {
						archive := archives.Archives[i]
						if clickableArchives[archive] == nil {
							clickableArchives[archive] = &widget.Clickable{}
						}
						clickableRow := clickableArchives[archive]
						return material.Clickable(gtx, clickableRow, func(gtx C) D {
							/*clip.Rect{
								Max: image.Point{
									X: gtx.Constraints.Min.X,
									Y: gtx.Constraints.Min.Y,
								},
							}.Add(gtx.Ops)
							paint.Fill(gtx.Ops, color.NRGBA{R: 0xff})*/

							return rowWidget(gtx, i)
						})
					})
				})
			})
		},
	}

	return material.List(th, listMain).Layout(
		gtx, len(widgets), func(gtx C, i int) D {
			return layout.UniformInset(unit.Dp(5)).Layout(gtx, widgets[i])
		},
	)
}

func playerEditorLayout(gtx C, th *material.Theme) D {
	e := material.Editor(th, playerLineEditor, "Enter player's name")
	e.Font.Style = text.Italic
	border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
	})
}
