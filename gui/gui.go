package gui

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

var th *material.Theme
var mainColor color.NRGBA

var header material.LabelStyle

var searchInput = &widget.Editor{
	SingleLine: true,
	Submit:     true,
}
var editor layout.Widget

func Run() {
	// initialize colors and theme
	mainColor = color.NRGBA{R: 127, G: 127, B: 87, A: 255}
	th = material.NewTheme(gofont.Collection())

	header = material.H3(th, "Cyber Cracker")
	header.Color = mainColor
	header.Alignment = text.Middle

	q := material.Editor(th, searchInput, "Query...")
	q.Color = mainColor
	q.HintColor = mainColor
	q.Font.Style = text.Italic

	border := widget.Border{Color: mainColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
	editor = func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, q.Layout)
		})
	}

	//searchInput.Focus()

	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case key.Event:
			log.Printf("Event %+v", e)
			if e.Modifiers.Contain(key.ModCtrl) {
				if e.Name == "c" {
					return nil
				}
			}
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			paint.Fill(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(header.Layout),
				layout.Rigid(editor),
				// layout.Rigid(header.Layout),
			)

			e.Frame(gtx.Ops)
		}
	}
}
