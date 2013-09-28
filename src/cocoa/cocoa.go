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
	c.AddMethod("openClick:", (*AppDelegate).OpenClick)
	c.AddMethod("shutdownAfterFinishClick:", (*AppDelegate).ShutdownAfterFinishClick)
	// c.AddMethod("didActivateNotification:", (*AppDelegate).DidActivateNotification)

	objc.RegisterClass(c)
}

type AppDelegate struct {
	objc.Object `objc:"GOAppDelegate : NSObject"`
}

func (delegate *AppDelegate) OpenClick(sender uintptr) {
	if t, ok := task.GetDownloadingTask(); ok {
		cmd := exec.Command("open", path.Join(util.ReadConfig("dir"), t.Name))
		cmd.Start()
	} else {
		cmd := exec.Command("open", "/Applications/V'ger.app")
		cmd.Start()
	}
}

func (delegate *AppDelegate) ShutdownAfterFinishClick(sender uintptr) {
	item := NSMenuItem{objc.NewObject(sender)}

	if util.ToggleBoolConfig("shutdown-after-finish") {
		item.SetState(NSOnState)
	} else {
		item.SetState(NSOffState)
	}
}

type statusItemData struct {
	sync.RWMutex
	title   string
	tooltip string
}

type taskStatusItem struct {
	name       string
	title      string
	tooltip    string
	statusItem NSStatusItem
}
type taskStatusItems struct {
	sync.RWMutex
	items map[string]*taskStatusItem
}

var downloadingTaskStatusItems taskStatusItems = taskStatusItems{sync.RWMutex{}, make(map[string]*taskStatusItem)}

func timerStart() {
	watch := make(chan *task.Task)

	log.Println("status bar watch task change: ", watch)
	task.WatchChange(watch)

	for t := range watch {
		var title string
		var tooltip string
		if t.Status == "Downloading" {
			title = fmt.Sprintf("%s %.1f%%", util.CleanMovieName(t.Name)[:15],
				float64(t.DownloadedSize)/float64(t.Size)*100.0)
			tooltip = fmt.Sprintf("%.2f KB/s %s", t.Speed, t.Est)
		} else if t.Status == "Playing" {
			title = fmt.Sprintf("%s %.1f KB/s", util.CleanMovieName(t.Name)[:15], t.Speed)
			tooltip = ""
		} else {
			title = ""
			tooltip = ""
		}

		downloadingTaskStatusItems.Lock()
		if item, ok := downloadingTaskStatusItems.items[t.Name]; ok {
			item.title = title
			item.tooltip = tooltip
		} else {
			downloadingTaskStatusItems.items[t.Name] =
				&taskStatusItem{t.Name, title, tooltip, NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}}
		}
		downloadingTaskStatusItems.Unlock()
	}
}

func Start() {
	runtime.LockOSThread()

	pool := NewNSAutoreleasePool()

	delegate := objc.GetClass("GOAppDelegate").Alloc().Init()

	app := NSSharedApplication()

	NSDefaultUserNotificationCenter().SetDelegate(delegate)

	mainItem := NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	mainItem.SetHighlightMode(true)
	mainItem.SetTarget(delegate.Pointer())
	mainItem.SetTitle("V'ger")

	menu := NewNSMenuWithTitle("V'ger")
	mainItem.SetMenu(menu)

	itemShutdownAfterFinish := NewNSMenuItem("Shutdown after finish", objc.GetSelector("shutdownAfterFinishClick:"), "")
	itemShutdownAfterFinish.SetTarget(delegate)
	if util.ReadBoolConfig("shutdown-after-finish") {
		itemShutdownAfterFinish.SetState(NSOnState)
	} else {
		itemShutdownAfterFinish.SetState(NSOffState)
	}
	menu.AddItem(itemShutdownAfterFinish)

	go timerStart()

	for {
		pool.Release()
		pool = NewNSAutoreleasePool()

		event := app.NextEventMatchingMask(0xffffff, NSDateWithTimeIntervalSinceNow(1),
			"kCFRunLoopDefaultMode", true)

		app.SendEvent(event)

		downloadingTaskStatusItems.Lock()
		for _, item := range downloadingTaskStatusItems.items {
			item.statusItem.SetTitle(item.title)
			item.statusItem.SetToolTip(item.tooltip)
		}
		downloadingTaskStatusItems.Unlock()
	}

	mainItem.Release()
	downloadingTaskStatusItems.Lock()
	for _, item := range downloadingTaskStatusItems.items {
		item.statusItem.Release()
	}
	downloadingTaskStatusItems.items = nil
	downloadingTaskStatusItems.Unlock()
	pool.Release()
}
