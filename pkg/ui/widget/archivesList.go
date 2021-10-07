package widget

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/nmaupu/chesscom_exporter/pkg/model"
	"sync"
)

// ArchiveList is a collection of ArchiveRow that can be layout
type ArchiveList struct {
	rows  []*ArchiveRow
	theme *material.Theme
	list  widget.List
	mutex sync.Mutex

	selectAll  widget.Clickable
	selectNone widget.Clickable
}

func NewArchiveList(th *material.Theme) *ArchiveList {
	return &ArchiveList{
		rows:  nil,
		theme: th,
		list: widget.List{
			List: layout.List{
				Axis:        layout.Vertical,
				Alignment:   layout.Middle,
				ScrollToEnd: false,
			},
		},
	}
}

// AddRows adds one or more rows to the list
func (a *ArchiveList) AddRows(archives *model.ChesscomArchives) {
	if archives == nil {
		return
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	lenArchives := len(archives.Archives)

	if a.rows == nil {
		a.rows = make([]*ArchiveRow, lenArchives)
	}

	// Revert order
	for i, arch := range archives.Archives {
		a.rows[len(archives.Archives)-i-1] = NewArchiveRow(a.theme, arch)
	}
}

func (a *ArchiveList) ResetList() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.rows = nil
}

func (a *ArchiveList) IsNil() bool {
	return a.rows == nil
}

func (a *ArchiveList) Size() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return len(a.rows)
}

// AtLeastOneSelected returns true if at least one element is selected, false otherwise
func (a *ArchiveList) AtLeastOneSelected() bool {
	for _, arch := range a.rows {
		if arch.checkbox.Value {
			return true
		}
	}
	return false
}

func (a *ArchiveList) GetSelectedArchives() model.ChesscomArchives {
	archives := model.ChesscomArchives{}
	for _, arch := range a.rows {
		if arch.checkbox.Value {
			archives.Archives = append(archives.Archives, arch.Archive)
		}
	}
	return archives
}

func (a *ArchiveList) Layout(gtx layout.Context) layout.Dimensions {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	selectAllClicked := a.selectAll.Clicked()
	selectNoneClicked := a.selectNone.Clicked()
	if selectAllClicked || selectNoneClicked {
		for _, arch := range a.rows {
			arch.checkbox.Value = selectAllClicked
		}
	}

	return layout.Flex{
		Spacing:   layout.SpaceEnd,
		Axis:      layout.Vertical,
		Alignment: layout.Start,
	}.Layout(gtx,
		layout.Rigid(a.layoutHeader),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(a.theme, &a.list).Layout(
				gtx,
				len(a.rows),
				func(gtx layout.Context, i int) layout.Dimensions {
					row := a.rows[i]
					if row == nil {
						return layout.Dimensions{}
					}
					return row.Layout(gtx)
				},
			)
		}),
	)
}

func (a *ArchiveList) layoutHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceStart,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &a.selectAll, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(a.theme, unit.Dp(16), "All")
				lbl.Font.Style = text.Italic
				if a.selectAll.Hovered() {
					lbl.Font.Weight = text.DemiBold
				}
				return lbl.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Rigid(material.Label(a.theme, unit.Dp(16), "|").Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &a.selectNone, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(a.theme, unit.Dp(16), "None")
				lbl.Font.Style = text.Italic
				if a.selectNone.Hovered() {
					lbl.Font.Weight = text.DemiBold
				}
				return lbl.Layout(gtx)
			})
		}),
	)
}
