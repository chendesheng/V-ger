package cocoa

import (
	"fmt"
	"github.com/mkrautz/objc"
	. "github.com/mkrautz/objc/AppKit"
	. "github.com/mkrautz/objc/Foundation"
	// "log"
	"runtime"
	"task"
	// "unsafe"
	"util"
)

func NewSpeedMenuItem(speed int, target objc.Object) NSMenuItem {
	item := NewNSMenuItem(fmt.Sprintf("Up to %dk/s", speed), objc.GetSelector("speedClick:"), "")
	item.SetTarget(target)
	if speed == util.ReadIntConfig("max-speed") {
		item.SetState(NSOnState)
	}

	return item
}
func NewSimultaneousItem(cnt int, target objc.Object) NSMenuItem {
	item := NewNSMenuItem(fmt.Sprintf("Up to %d task(s)", cnt), objc.GetSelector("simultaneousClick:"), "")
	item.SetTarget(target)
	if cnt == util.ReadIntConfig("simultaneous-downloads") {
		item.SetState(NSOnState)
	}

	return item
}
func NewNoLimitMenuItem(target objc.Object) NSMenuItem {
	item := NewNSMenuItem("No Limit", objc.GetSelector("speedClick:"), "")
	item.SetTarget(target)
	if 0 == util.ReadIntConfig("max-speed") {
		item.SetState(NSOnState)
	}
	return item
}
func watchTaskChange(chTaskChange chan *task.Task) {
	watcher := make(chan *task.Task)
	task.WatchChange(watcher)

	for t := range watcher {
		go func() {
			chTaskChange <- t
		}()
	}
}
func Start() {
	runtime.LockOSThread()

	pool := NewNSAutoreleasePool()

	app := NSSharedApplication()

	delegate := &AppDelegate{Object: objc.NewInstance("GOAppDelegate")}
	app.SetDelegate(delegate)

	delegate.ApplicationDidFinishLaunching(objc.NilObject())

	chTaskChange := make(chan *task.Task)
	go watchTaskChange(chTaskChange)
	for {
		pool.Release()
		pool = NewNSAutoreleasePool()

		event := app.NextEventMatchingMask(0xffffff, NSDateWithTimeIntervalSinceNow(1),
			"kCFRunLoopDefaultMode", true)

		app.SendEvent(event)

		select {
		case t := <-chTaskChange:
			delegate.updateStatusBar(t)
			break
		default:
			break
		}
	}
}
func SendNotification(title, infoText string) error {
	pool := NewNSAutoreleasePool()

	notification := NSUserNotification{objc.GetClass("NSUserNotification").Alloc().Init()}
	notification.SetTitle(title)
	notification.SetInformativeText(infoText)
	notification.SetSoundName(NSUserNotificationDefaultSoundName)
	notification.SetHasActionButton(true)
	notification.SetActionButtonTitle("Open")

	center := NSDefaultUserNotificationCenter()
	center.DeliverNotification(notification)

	pool.Release()
	return nil
}
