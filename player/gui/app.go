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

func onOpenOpenPanel() {
	if appDelegate != nil {
		appDelegate.OnOpenOpenPanel()
	}
}

func onMenuClick(typ int, tag int) int {
	if appDelegate != nil {
		return appDelegate.OnMenuClick(typ, tag)
	}
	return 0
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
	OnMenuClick(int, int) int
}

var appDelegate AppDelegate

func Run(d AppDelegate) {
	appDelegate = d

	run()
	log.Println("gui initialized")
}
