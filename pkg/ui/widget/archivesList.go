package widget

import (
	"gioui.org/layout"
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
}

func NewArchiveList(th *material.Theme) *ArchiveList {
	return &ArchiveList{
		rows:  nil,
		theme: th,
		list: widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
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

	if a.rows == nil {
		a.rows = make([]*ArchiveRow, 0)
	}

	for _, arch := range archives.Archives {
		a.rows = append(a.rows, NewArchiveRow(a.theme, arch))
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

func (a *ArchiveList) Layout(gtx layout.Context) layout.Dimensions {
	a.mutex.Lock()
	defer a.mutex.Unlock()
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
}
