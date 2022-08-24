package gui

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/fasmide/gipass/store"

	"gioui.org/font/gofont"
)

var (
	th        *material.Theme
	mainColor color.NRGBA
)

var (
	header      material.LabelStyle
	searchInput = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}

	editor layout.Widget

	results    widget.List
	resultView layout.Widget
)

var (
	selected  int
	lastQuery []store.Result
)

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

	resultView = func(gtx layout.Context) layout.Dimensions {
		return material.List(th, &results).Layout(gtx, len(lastQuery), func(gtx layout.Context, i int) layout.Dimensions {
			m := material.Body1(th, fmt.Sprintf("%s @ %s", lastQuery[i].Username, lastQuery[i].URL))
			if selected == i {
				m.Color = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
			} else {
				m.Color = mainColor
			}

			return m.Layout(gtx)
		})
	}

	go func() {
		w := app.NewWindow(app.Title("Cyber Cracker"))
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
			handleQuery(&gtx)

			handleEscape(&gtx)
			handleNavigation(&gtx)

			paint.Fill(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(header.Layout),
				layout.Rigid(editor),
				layout.Flexed(1, resultView),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func selectResult(gtx *layout.Context) {
	pass, err := lastQuery[selected].CleartextPassword()
	if err != nil {
		panic(err)
	}
	log.Printf("%X != %X", "Wwr3NrNfrTTxHPh", pass)
	clipboard.WriteOp{Text: pass}.Add(gtx.Ops)
	op.InvalidateOp{}.Add(gtx.Ops)
}

func query(gtx *layout.Context, input string) {
	now := time.Now()

	var err error
	lastQuery, err = store.Logins.Query(input)
	if err != nil {
		panic(err)
	}

	log.Printf("%s: %d results, took %s", input, len(lastQuery), time.Now().Sub(now))

	// it should be okay to move the selected position up when a new result arrives
	selected = 0

	op.InvalidateOp{}.Add(gtx.Ops)
}

var submitTag = &struct{}{}

func handleQuery(gtx *layout.Context) {
	for _, e := range searchInput.Events() {

		switch e.(type) {
		case widget.ChangeEvent:
			query(gtx, searchInput.Text())

		case widget.SubmitEvent:
			selectResult(gtx)
		}
	}

	for _, i := range gtx.Events(&submitTag) {
		ki, ok := i.(key.Event)
		if !ok {
			continue
		}

		if ki.State == key.Press {
			selectResult(gtx)
		}
	}
	key.InputOp{Tag: &submitTag, Keys: key.NameReturn}.Add(gtx.Ops)
}

var escapeTag = &struct{}{}

func handleEscape(gtx *layout.Context) {
	for _, i := range gtx.Events(&escapeTag) {
		ki, ok := i.(key.Event)
		if !ok {
			continue
		}

		if ki.State == key.Press {
			os.Exit(0)
		}
	}
	key.InputOp{Tag: &escapeTag, Keys: key.NameEscape}.Add(gtx.Ops)
}

// The navigation is somewhat broken as the Editor widget eats some arrow keys
// depending on where the cursor is placed - there are some hacks we can explore
// but for now, we are just going to live with it and maybe the gioui will
// evolve enough to have a pleasent solution in the future
//
// see: https://gophers.slack.com/archives/CM87SNCGM/p1660297386051799
var navigationTag = &struct{}{}

func handleNavigation(gtx *layout.Context) {
	for _, i := range gtx.Events(&navigationTag) {
		ki, ok := i.(key.Event)
		if !ok {
			continue
		}

		// could be key.Release which we dont care about
		if ki.State != key.Press {
			continue
		}

		if ki.Name == key.NameUpArrow {
			selected = selected - 1

		}
		if ki.Name == key.NameDownArrow {
			selected = selected + 1

		}

		// make bounds wrap around
		if selected > len(lastQuery)-1 {
			selected = 0
		}
		if selected < 0 {
			selected = len(lastQuery) - 1
		}
	}
	key.InputOp{Tag: &navigationTag, Keys: "[↑,↓]"}.Add(gtx.Ops)

}
