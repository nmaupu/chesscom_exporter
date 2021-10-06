package widget

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/nmaupu/chesscom_exporter/pkg/model"
)

// ArchiveRow is a struct representing an archive with a checkbox.
// The entire widget is clickable and check the inner checkbox.
type ArchiveRow struct {
	widget.Clickable

	Archive  model.ChesscomArchive
	checkbox widget.Bool
	lblYear  widget.Label
	lblMonth widget.Label
	theme    *material.Theme
}

func NewArchiveRow(th *material.Theme, archive model.ChesscomArchive) *ArchiveRow {
	return &ArchiveRow{
		Archive: archive,
		theme:   th,
	}
}

func (a *ArchiveRow) Layout(gtx layout.Context) layout.Dimensions {
	return material.Clickable(gtx,
		&a.Clickable,
		func(gtx layout.Context) layout.Dimensions {
			return a.layoutRow(gtx)
		},
	)
}

func (a *ArchiveRow) layoutRow(gtx layout.Context) layout.Dimensions {
	lblYear := material.Label(a.theme, unit.Dp(16), fmt.Sprintf("%d", a.Archive.GetYear()))
	lblMonth := material.Label(a.theme, unit.Dp(16), fmt.Sprintf("%s", a.Archive.GetMonthAsString()))
	checkbox := material.CheckBox(a.theme, &a.checkbox, "")

	// if row is clicked, change the checkbox' state
	if a.Clicked() {
		checkbox.CheckBox.Value = !checkbox.CheckBox.Value
	}

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
