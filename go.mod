module github.com/nmaupu/chesscom_exporter

go 1.17

require (
	gioui.org v0.0.0-20211003134802-50476239f6a3
	gioui.org/x/explorer v0.0.0-20210929182633-199c05a62a31
	golang.design/x/clipboard v0.5.3
)

require (
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.4 // indirect
	git.wow.st/gmp/jni v0.0.0-20200827154156-014cd5c7c4c0 // indirect
	golang.org/x/exp v0.0.0-20210722180016-6781d3edade3 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/mobile v0.0.0-20210716004757-34ab1303b554 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace gioui.org/x/explorer => ../gio-x/explorer
