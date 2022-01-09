package src

import (
	"github.com/vaiktorg/grimoire/log/logger"
)

type Table struct {
	app.Compo
	Msgs []logger.Log
}

func (t *Table) Render() app.UI {
	return app.Table().Class("table table-hover").Body(
		app.THead().Body(
			app.Tr().Body(
				//TODO Dynamically generate these.
				app.Th().Scope("col").Text("Expires"),
				app.Th().Scope("col").Text("Level"),
				app.Th().Scope("col").Text("Message"),
				app.Th().Scope("col").Text("data"),
			),
		),
		app.TBody().Body(
			app.Range(t.Msgs).Slice(func(s int) app.UI {
				msg := t.Msgs[s]
				return app.Tr().OnClick(func(ctx app.Context, e app.Event) {
					//TODO WTF SHOULD THEY DO WHEN CLICKED
				}).Body(
					app.Th().Scope("row").Text(msg.Timestamp),
					app.Td().Text(msg.Level),
					app.Td().Text(msg.Msg),
					app.Td().Text(msg.Data),
				)
			}),
		),
	)
}

// TODO: Check for duplicity when updating content from array
func (t *Table) AddMsgs(msgs ...logger.Log) {
	t.Msgs = append(t.Msgs, msgs...)
	t.Update()
}
