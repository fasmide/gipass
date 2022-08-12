package gui

import (
	"fmt"
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

var results widget.List
var result layout.Widget
var selected int

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

	searchInput.Focus()

	results = widget.List{List: layout.List{Axis: layout.Vertical}}

	result = func(gtx layout.Context) layout.Dimensions {
		return material.List(th, &results).Layout(gtx, 700, func(gtx layout.Context, i int) layout.Dimensions {
			m := material.Body1(th, fmt.Sprintf("index: %d", i))
			if selected == i {
				m.Color = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
			} else {
				m.Color = mainColor
			}

			return m.Layout(gtx)
		})
	}

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
			handleInput(&gtx)

			paint.Fill(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(header.Layout),
				layout.Rigid(editor),
				layout.Flexed(1, result),
			)

			e.Frame(gtx.Ops)
		}
	}
}

var inputTag *struct{}

func handleInput(gtx *layout.Context) {
	// TODO: initialize this elsewhere
	if inputTag == nil {
		inputTag = &struct{}{}
	}

	for idx, i := range gtx.Events(inputTag) {
		log.Printf("Input queue (%d): %T: %+v", idx, i, i)

		// Look for key.Event - change selected accordingly
		ki, ok := i.(key.Event)
		if !ok {
			continue
		}

		if ki.State == key.Press {
			if ki.Name == key.NameUpArrow {
				selected = selected - 1
			}

			if ki.Name == key.NameDownArrow {
				selected = selected + 1
			}
		}
	}

	key.InputOp{Tag: inputTag, Keys: "[←,→,↑,↓]"}.Add(gtx.Ops)

}
