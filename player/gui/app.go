package gui

import "log"

func onOpenFile(filename string) bool {
	if appDelegate != nil {
		return appDelegate.OpenFile(filename)
	} else {
		return false
	}
}

func onWillTerminate() {
	if appDelegate != nil {
		appDelegate.WillTerminate()
	}
}

func onSearchSubtitleMenuItemClick() {
	if appDelegate != nil {
		appDelegate.ToggleSearchSubtitle()
	}
}

func onOpenOpenPanel() {
	if appDelegate != nil {
		appDelegate.OnOpenOpenPanel()
	}
}

func onCloseOpenPanel(name string) {
	if appDelegate != nil {
		appDelegate.OnCloseOpenPanel(name)
	}
}

type AppDelegate interface {
	OpenFile(string) bool
	OnOpenOpenPanel()
	OnCloseOpenPanel(filename string)
	WillTerminate()
	ToggleSearchSubtitle()
	OnMenuClick(int)
}

var appDelegate AppDelegate

func Run(d AppDelegate) {
	appDelegate = d

	run()
	log.Println("gui initialized")
}
