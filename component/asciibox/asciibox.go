package asciibox

import (
	ui "github.com/gizak/termui/v3"
	fl "github.com/mbndr/figlet4go"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/data"
	"image"
)

type AsciiBox struct {
	ui.Block
	data.Consumer
	*component.Alerter
	text    string
	ascii   string
	style   ui.Style
	render  *fl.AsciiRender
	options *fl.RenderOptions
}

const asciiFontExtension = ".flf"

func NewAsciiBox(c config.AsciiBoxConfig) *AsciiBox {

	consumer := data.NewConsumer()
	block := *ui.NewBlock()
	block.Title = c.Title

	options := fl.NewRenderOptions()
	options.FontName = string(*c.Font)

	fontStr, err := asset.Asset(options.FontName + asciiFontExtension)
	if err != nil {
		panic("Can't load the font: " + err.Error())
	}
	render := fl.NewAsciiRender()
	_ = render.LoadBindataFont(fontStr, options.FontName)

	box := AsciiBox{
		Block:    block,
		Consumer: consumer,
		Alerter:  component.NewAlerter(consumer.AlertChannel),
		style:    ui.NewStyle(*c.Color),
		render:   render,
		options:  options,
	}

	go box.consume()

	return &box
}

func (a *AsciiBox) consume() {
	for {
		select {
		case sample := <-a.SampleChannel:
			a.text = sample.Value
			a.ascii, _ = a.render.RenderOpts(sample.Value, a.options)
			//case alert := <-a.alertChannel:
			// TODO base alerting mechanism
		}
	}
}

func (a *AsciiBox) Draw(buffer *ui.Buffer) {

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(ui.ColorBlack)), a.GetRect())
	a.Block.Draw(buffer)

	point := a.Inner.Min
	cells := ui.ParseStyles(a.ascii, a.style)

	for i := 0; i < len(cells) && point.Y < a.Inner.Max.Y; i++ {
		if cells[i].Rune == '\n' {
			point = image.Pt(a.Inner.Min.X, point.Y+1)
		} else if point.In(a.Inner) {
			buffer.SetCell(cells[i], point)
			point = point.Add(image.Pt(1, 0))
		}
	}
}
