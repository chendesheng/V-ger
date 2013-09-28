package cocoa

import (
	"download"
	"fmt"
	"github.com/mkrautz/objc"
	. "github.com/mkrautz/objc/AppKit"
	. "github.com/mkrautz/objc/Foundation"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	// "path"
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
	c.AddMethod("speedClick:", (*AppDelegate).SpeedClick)
	c.AddMethod("simultaneousClick:", (*AppDelegate).SimultaneousClick)
	// c.AddMethod("didActivateNotification:", (*AppDelegate).DidActivateNotification)

	objc.RegisterClass(c)
}

type AppDelegate struct {
	objc.Object `objc:"GOAppDelegate : NSObject"`
}

func (delegate *AppDelegate) OpenClick(sender uintptr) {
	// if t, ok := task.GetDownloadingTask(); ok {
	// 	cmd := exec.Command("open", path.Join(util.ReadConfig("dir"), t.Name))
	// 	cmd.Start()
	// } else {
	cmd := exec.Command("open", "/Applications/V'ger.app")
	cmd.Start()
	// }
}

func (delegate *AppDelegate) ShutdownAfterFinishClick(sender uintptr) {
	item := NSMenuItem{objc.NewObject(sender)}

	if util.ToggleBoolConfig("shutdown-after-finish") {
		item.SetState(NSOnState)
	} else {
		item.SetState(NSOffState)
	}
}
func (delegate *AppDelegate) SpeedClick(sender uintptr) {
	item := NSMenuItem{objc.NewObject(sender)}
	title := item.Title()
	speed := 0
	if title != "No Limit" {
		speedReg := regexp.MustCompile("Up to (\\d+)")
		speedStr := speedReg.FindStringSubmatch(title)[1]
		speed, _ = strconv.Atoi(speedStr)
	}

	util.SaveConfig("max-speed", fmt.Sprint(speed))
	download.LimitSpeed(int64(speed))

	for _, item := range speedMenuItems {
		item.SetState(NSOffState)
	}

	item.SetState(NSOnState)
}
func (delegate *AppDelegate) SimultaneousClick(sender uintptr) {
	item := NSMenuItem{objc.NewObject(sender)}
	title := item.Title()
	cntReg := regexp.MustCompile("Up to (\\d+)")
	cntStr := cntReg.FindStringSubmatch(title)[1]
	cnt, _ := strconv.Atoi(cntStr)

	downloadingCnt := task.NumOfDownloadingTasks()

	for i := cnt; i < downloadingCnt; i++ {
		task.QueueDownloadingTask()
	}

	for i := downloadingCnt; i < cnt; i++ {
		task.ResumeNextTask()
	}

	util.SaveConfig("simultaneous-downloads", fmt.Sprint(cnt))

	for _, item := range simultaneousMenuItems {
		item.SetState(NSOffState)
	}

	item.SetState(NSOnState)
}

type statusItemData struct {
	sync.RWMutex
	title   string
	tooltip string
}

type taskStatusItem struct {
	title      string
	tooltip    string
	statusItem NSStatusItem
}
type taskStatusItems struct {
	sync.RWMutex
	items map[string]*taskStatusItem
}

var downloadingTaskStatusItems taskStatusItems = taskStatusItems{sync.RWMutex{}, make(map[string]*taskStatusItem)}

var speedMenuItems []NSMenuItem = make([]NSMenuItem, 0)
var simultaneousMenuItems []NSMenuItem = make([]NSMenuItem, 0)

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
				&taskStatusItem{title, tooltip, NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}}
		}
		downloadingTaskStatusItems.Unlock()
	}
}

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

func Start() {
	runtime.LockOSThread()

	pool := NewNSAutoreleasePool()

	delegate := objc.GetClass("GOAppDelegate").Alloc().Init()

	app := NSSharedApplication()

	NSDefaultUserNotificationCenter().SetDelegate(delegate)

	mainItem := NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	mainItem.SetHighlightMode(true)
	mainItem.SetTarget(delegate)
	mainItem.SetToolTip("V'ger")

	img := NewNSImageWithContentOfFile("assets/icon.png")
	img.SetTemplate(true)
	mainItem.SetImage(img)

	menu := NewNSMenuWithTitle("V'ger")
	mainItem.SetMenu(menu)

	shutdownAfterFinishItem := NewNSMenuItem("Shutdown after finish", objc.GetSelector("shutdownAfterFinishClick:"), "")
	shutdownAfterFinishItem.SetTarget(delegate)
	if util.ReadBoolConfig("shutdown-after-finish") {
		shutdownAfterFinishItem.SetState(NSOnState)
	} else {
		shutdownAfterFinishItem.SetState(NSOffState)
	}

	openItem := NewNSMenuItem("Open tasks panel", objc.GetSelector("openClick:"), "")
	openItem.SetTarget(delegate)

	speedItem := NewNSMenuItem("Speed", objc.GetSelector(""), "")
	simultaneousItem := NewNSMenuItem("Simultaneous", objc.GetSelector(""), "")

	speedMenuItems = append(speedMenuItems,
		NewNoLimitMenuItem(delegate),
		NewSpeedMenuItem(50, delegate),
		NewSpeedMenuItem(100, delegate),
		NewSpeedMenuItem(200, delegate),
		NewSpeedMenuItem(300, delegate))

	simultaneousMenuItems = append(simultaneousMenuItems,
		NewSimultaneousItem(1, delegate),
		NewSimultaneousItem(2, delegate),
		NewSimultaneousItem(3, delegate),
		NewSimultaneousItem(4, delegate))

	menu.AddItem(openItem)
	menu.AddItem(NSSeparatorMenuItem())
	menu.AddItem(speedItem)

	for _, item := range speedMenuItems {
		menu.AddItem(item)
	}

	menu.AddItem(NSSeparatorMenuItem())

	menu.AddItem(simultaneousItem)
	for _, item := range simultaneousMenuItems {
		menu.AddItem(item)
	}

	menu.AddItem(NSSeparatorMenuItem())
	menu.AddItem(shutdownAfterFinishItem)

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
