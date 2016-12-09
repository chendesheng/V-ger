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

func onFullScreen(action int) {
	if appDelegate != nil {
		appDelegate.OnFullScreen(action)
	}
}

func onCloseOpenPanel(name string) {
	if appDelegate != nil {
		appDelegate.OnCloseOpenPanel(name)
	}
}

func onMouseWheel(deltaX, deltaY float64) {
	if appDelegate != nil {
		appDelegate.OnMouseWheel(deltaX, deltaY)
	}
}

func onWillSleep() {
  log.Print("onWillSleep");
  appDelegate.OnWillSleep()
}

func onDidWake() {
		appDelegate.OnDidWake()
}

type AppDelegate interface {
	OpenFile(string) bool
	OnOpenOpenPanel()
	OnCloseOpenPanel(filename string)
	WillTerminate()
	OnMenuClick(int, int) int
	OnMouseWheel(float64, float64)
	OnFullScreen(int)
  OnWillSleep()
  OnDidWake()
}

var appDelegate AppDelegate

func Run(d AppDelegate) {
	appDelegate = d

	run()
	log.Println("gui initialized")
}
