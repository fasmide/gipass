package gui

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

var header material.LabelStyle
var th *material.Theme
var searchInput = &widget.Editor{
	SingleLine: true,
	Submit:     true,
}

var mainColor color.NRGBA

func Run() {
	// initialize colors and theme
	mainColor = color.NRGBA{R: 127, G: 127, B: 87, A: 255}
	th = material.NewTheme(gofont.Collection())

	header = material.H3(th, "Password CyberCracker")
	header.Color = mainColor
	header.Alignment = text.Middle

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
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			//paint.Fill(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(header.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, searchInput, "Input Query...")
					e.Font.Style = text.Italic
					border := widget.Border{Color: mainColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
					return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
					})
				},
				),
				// layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				// 		return material.List(theme, &list).Layout(gtx, len(messages), func(gtx layout.Context, i int) layout.Dimensions {
				// 				m := material.Body1(theme, messages[i])
				// 				m.Color = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
				// 				return m.Layout(gtx)
				// 		})
				// }),
				// layout.Rigid(footer.Layout),
			)

			e.Frame(gtx.Ops)
		}
	}
}
