package main

import (
	"github.com/rivo/tview"
)

type UI struct {
	App      *tview.Application
	Channels *tview.List
	Chat     *tview.TextView
	Members  *tview.List
	Input    *tview.InputField
}

func initUI() UI {
	app := tview.NewApplication()

	channels := tview.NewList()
	channels.SetTitle("Channels")
	channels.SetBorder(true)

	chat := tview.NewTextView()
	chat.SetTitle("Chat")
	chat.SetBorder(true)
	chat.SetScrollable(true)
	chat.SetChangedFunc(func() {
		chat.ScrollToEnd()
		app.QueueUpdateDraw(func() {})
	})
	chat.SetWordWrap(true)

	members := tview.NewList()
	members.SetTitle("Members")
	members.SetBorder(true)

	input := tview.NewInputField()
	input.SetLabel("> ")
	input.SetBorder(true)

	center := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chat, 0, 1, false).
		AddItem(input, 3, 0, true)

	layout := tview.NewFlex().
		AddItem(channels, 30, 0, false).
		AddItem(center, 0, 1, true).
		AddItem(members, 30, 0, false)

	app.SetRoot(layout, true).SetFocus(input)
	return UI{app, channels, chat, members, input}
}
