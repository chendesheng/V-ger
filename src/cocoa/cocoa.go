package cocoa

import (
	"fmt"
	"github.com/mkrautz/objc"
	. "github.com/mkrautz/objc/AppKit"
	. "github.com/mkrautz/objc/Foundation"
	"log"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"task"
	// "time"
	"util"
)

var config map[string]string

func init() {
	c := objc.NewClass(AppDelegate{})
	c.AddMethod("menuClick:", (*AppDelegate).MenuClick)
	// c.AddMethod("didActivateNotification:", (*AppDelegate).DidActivateNotification)

	objc.RegisterClass(c)
}

type AppDelegate struct {
	objc.Object `objc:"GOAppDelegate : NSObject"`
}

func (delegate *AppDelegate) MenuClick(sender uintptr) {
	if t, ok := task.GetDownloadingTask(); ok {
		cmd := exec.Command("open", path.Join(util.ReadConfig("dir"), t.Name))
		cmd.Start()
	} else {
		cmd := exec.Command("open", "/Applications/V'ger.app")
		cmd.Start()
	}
}

type statusItemData struct {
	sync.RWMutex
	title   string
	tooltip string
}

var currentStatusItem statusItemData = statusItemData{sync.RWMutex{}, "V'ger", "Speed is fun!"}

func timerStart() {
	watch := make(chan *task.Task)

	log.Println("status bar watch task change: ", watch)
	task.WatchChange(watch)

	for t := range watch {
		var title string
		var tooltip string
		if t.Status == "Downloading" {
			title = fmt.Sprintf("%s %.1f%%", util.CleanMovieName(t.Name),
				float64(t.DownloadedSize)/float64(t.Size)*100.0)
			tooltip = fmt.Sprintf("%.2f KB/s %s", t.Speed, t.Est)
		} else if t.Status == "Playing" {
			title = fmt.Sprintf("%s %.1f KB/s", util.CleanMovieName(t.Name), t.Speed)
			tooltip = ""
		} else {
			if !task.HasDownloadingOrPlaying() {
				title = "V'ger"
				tooltip = "Speed is fun!"
			}
		}

		currentStatusItem.Lock()
		currentStatusItem.title = title
		currentStatusItem.tooltip = tooltip
		currentStatusItem.Unlock()
	}
}

func Start() {
	runtime.LockOSThread()

	pool := NewNSAutoreleasePool()

	delegate := objc.GetClass("GOAppDelegate").Alloc().Init()

	app := NSSharedApplication()

	NSDefaultUserNotificationCenter().SetDelegate(delegate)

	statusItem := NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	statusItem.SetHighlightMode(true)
	statusItem.SetTarget(delegate.Pointer())
	statusItem.SetAction(objc.GetSelector("menuClick:"))
	statusItem.SetTitle(currentStatusItem.title)
	statusItem.SetToolTip(currentStatusItem.tooltip)

	go timerStart()

	for {
		pool.Release()
		pool = NewNSAutoreleasePool()

		event := app.NextEventMatchingMask(0xffffff, NSDateWithTimeIntervalSinceNow(1),
			"kCFRunLoopDefaultMode", true)

		app.SendEvent(event)

		currentStatusItem.RLock()
		title := currentStatusItem.title
		tooltip := currentStatusItem.tooltip
		currentStatusItem.RUnlock()

		statusItem.SetTitle(title)
		statusItem.SetToolTip(tooltip)
	}

	statusItem.Release()
	pool.Release()
}
