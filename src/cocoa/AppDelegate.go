package cocoa

import (
	"download"
	"fmt"
	"github.com/mkrautz/objc"
	. "github.com/mkrautz/objc/AppKit"
	. "github.com/mkrautz/objc/Foundation"
	"log"
	"native"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"thunder"
	// "path"
	// "runtime"
	// "sync"
	"task"
	"util"
)

// type statusItemData struct {
// 	sync.RWMutex
// 	title   string
// 	tooltip string
// }

// type taskStatusItem struct {
// 	title      string
// 	tooltip    string
// 	statusItem NSStatusItem
// }
// type taskStatusItems struct {
// 	sync.RWMutex
// 	items map[string]*taskStatusItem
// }

// var downloadingTaskStatusItems map[string]NSStatusItem = make(map[string]NSStatusItem) //taskStatusItems{sync.RWMutex{}, make(map[string]*taskStatusItem)}

func init() {
	c := objc.NewClass(AppDelegate{})
	c.AddMethod("openClick:", (*AppDelegate).OpenClick)
	c.AddMethod("shutdownAfterFinishClick:", (*AppDelegate).ShutdownAfterFinishClick)
	c.AddMethod("speedClick:", (*AppDelegate).SpeedClick)
	c.AddMethod("simultaneousClick:", (*AppDelegate).SimultaneousClick)
	c.AddMethod("newTaskFromPasteboardClick:", (*AppDelegate).NewTaskFromPasteboardClick)

	c.AddMethod("applicationDidFinishLaunching:", (*AppDelegate).ApplicationDidFinishLaunching)
	c.AddMethod("userNotificationCenter:didActivateNotification:", (*AppDelegate).DidActivateNotification)
	c.AddMethod("userNotificationCenter:shouldPresentNotification:", (*AppDelegate).ShouldPresentNotification)

	objc.RegisterClass(c)

	native.DefaultNativeAPI = new(cocoaNativeAPI)
}

type AppDelegate struct {
	objc.Object `objc:"GOAppDelegate : NSObject"`
	// Window      objc.Object `objc:"IBOutlet"`

	simultaneousMenuItems []NSMenuItem
	speedMenuItems        []NSMenuItem
	window                NSWindow

	downloadingTaskStatusItems map[string]NSStatusItem
}

var appDelegate *AppDelegate

func (delegate *AppDelegate) ShouldPresentNotification(center objc.Object) bool {
	return true
}
func (delegate *AppDelegate) DidActivateNotification(center objc.Object, notification objc.Object) {
	noti := NSUserNotification{notification}
	title := noti.Title()
	if strings.Contains(title, "Sleep") {
		if quit := native.DefaultNativeAPI.(*cocoaNativeAPI).shutdownQuit; quit != nil {
			close(quit)
		}
	} else if strings.Contains(title, "Finish") {
		name := noti.InformativeText()
		cmd := exec.Command("open", path.Join(util.ReadConfig("dir"), name))
		cmd.Start()
	} else {
		delegate.OpenClick(0)
	}
}

func (delegate *AppDelegate) ApplicationDidFinishLaunching(obj objc.Object) {
	appDelegate = delegate

	mainItem := NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	mainItem.SetHighlightMode(true)
	mainItem.SetTarget(delegate)
	mainItem.SetToolTip("V'ger")

	img := NewNSImageWithContentOfFile("assets/icon.png")
	img.SetTemplate(true)
	mainItem.SetImage(img)

	shutdownAfterFinishItem := NewNSMenuItem("Sleep after finish", objc.GetSelector("shutdownAfterFinishClick:"), "")
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

	newTaskFromPasteboardItem := NewNSMenuItem("New task from pasteboard", objc.GetSelector("newTaskFromPasteboardClick:"), "")
	newTaskFromPasteboardItem.SetTarget(delegate)

	menu := NewNSMenuWithTitle("V'ger")
	mainItem.SetMenu(menu)
	delegate.speedMenuItems = append(delegate.speedMenuItems,
		NewNoLimitMenuItem(delegate),
		NewSpeedMenuItem(50, delegate),
		NewSpeedMenuItem(100, delegate),
		NewSpeedMenuItem(200, delegate),
		NewSpeedMenuItem(300, delegate))

	delegate.simultaneousMenuItems = append(delegate.simultaneousMenuItems,
		NewSimultaneousItem(1, delegate),
		NewSimultaneousItem(2, delegate),
		NewSimultaneousItem(3, delegate),
		NewSimultaneousItem(4, delegate))

	menu.AddItem(openItem)
	menu.AddItem(newTaskFromPasteboardItem)
	menu.AddItem(NSSeparatorMenuItem())
	menu.AddItem(speedItem)

	for _, item := range delegate.speedMenuItems {
		menu.AddItem(item)
	}

	menu.AddItem(NSSeparatorMenuItem())

	menu.AddItem(simultaneousItem)
	for _, item := range delegate.simultaneousMenuItems {
		menu.AddItem(item)
	}

	menu.AddItem(NSSeparatorMenuItem())
	menu.AddItem(shutdownAfterFinishItem)

	// ScheduledTimerWithTimeInterval(1, delegate, objc.GetSelector("timerTick:"), objc.NilObject(), true)

	center := NSDefaultUserNotificationCenter()
	center.SetDelegate(delegate.Object)
}
func (delegate *AppDelegate) updateStatusBar(t *task.Task) {
	if delegate.downloadingTaskStatusItems == nil {
		delegate.downloadingTaskStatusItems = make(map[string]NSStatusItem)
	}

	downloadingTaskStatusItems := delegate.downloadingTaskStatusItems

	var title string
	var tooltip string
	if t.Status == "Downloading" {
		title = fmt.Sprintf("%s %.1f%%", util.CleanMovieNameWithMaxLen(t.Name, 15),
			float64(t.DownloadedSize)/float64(t.Size)*100.0)
		est := ""
		if t.Est > 0 {
			est = fmt.Sprintf(" %s", t.Est)
		}
		tooltip = fmt.Sprintf("%.2f KB/s%s", t.Speed, est)
	} else if t.Status == "Playing" {
		title = fmt.Sprintf("%s %.1f KB/s", util.CleanMovieNameWithMaxLen(t.Name, 15), t.Speed)
		tooltip = ""
	} else {
		title = ""
		tooltip = ""
	}
	if item, ok := downloadingTaskStatusItems[t.Name]; ok {
		if title == "" {
			delete(downloadingTaskStatusItems, t.Name)
			NSSystemStatusBar().RemoveStatusItem(item)
			item.Release()
		} else {
			item.SetTitle(title)
			item.SetToolTip(tooltip)
		}
	} else {
		downloadingTaskStatusItems[t.Name] = NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	}
}
func (delegate *AppDelegate) OpenClick(sender uintptr) {
	var openUrl string
	if util.IsPathExists(util.ReadConfig("client-app")) {
		openUrl = util.ReadConfig("client-app")
	} else {
		openUrl = "http://" + util.ReadConfig("server")
	}
	cmd := exec.Command("open", openUrl)
	cmd.Start()
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
	delegate = appDelegate

	item := NSMenuItem{objc.NewObject(sender)}
	title := item.Title()
	speed := 0
	if title != "No Limit" {
		speedReg := regexp.MustCompile("Up to (\\d+)")
		speedStr := speedReg.FindStringSubmatch(title)[1]
		speed, _ = strconv.Atoi(speedStr)
	}

	util.SaveConfig("max-speed", fmt.Sprint(speed))
	download.LimitSpeed(speed)

	println(delegate.speedMenuItems)
	println(delegate)
	for _, item := range delegate.speedMenuItems {
		println(item.Title())
		item.SetState(NSOffState)
	}

	item.SetState(NSOnState)
}

func (delegate *AppDelegate) SimultaneousClick(sender uintptr) {
	delegate = appDelegate

	senderItem := NSMenuItem{objc.NewObject(sender)}
	title := senderItem.Title()
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

	for _, item := range delegate.simultaneousMenuItems {
		item.SetState(NSOffState)
	}

	senderItem.SetState(NSOnState)
}
func (delegate *AppDelegate) NewTaskFromPasteboardClick(sender uintptr) {
	pasteboard := NSGeneralPasteboard()
	str := pasteboard.StringForType(NSPasteboardTypeString)

	go func() {
		if len(str) > 0 {
			log.Print("From pasteboard: ", str)

			url := ""
			name := ""
			regFormat := regexp.MustCompile(".(mkv|avi|mp4|rmvb|rm|wmv)$")
			files, err := thunder.NewTask(str, "")
			if err != nil {
				SendNotification("Download failed", err.Error())
				return
			}
			for _, f := range files {
				if f.Percent == 100 && regFormat.Match([]byte(f.Name)) {
					url = f.DownloadURL
					name = f.Name
					break
				}
			}

			if len(url) == 0 {
				SendNotification("Download failed", "No file's ready")
				return
			}

			_, name2, size, err := download.GetDownloadInfo(url)
			if err != nil {
				_, name2, size, err = download.GetDownloadInfo(url)

				if err != nil {
					SendNotification("Download failed", err.Error())
				}
				return
			}

			if name == "" {
				name = name2
			}

			fmt.Printf("add download \"%s\".\nname: %s\n", url, name)

			if t, err := task.GetTask(name); err == nil {
				if t.Status == "Finished" {

				} else {
					log.Print("task already exists")
					task.ResumeTask(name)
				}
			} else if err := task.StartNewTask(name, url, size); err != nil {
				SendNotification("Download failed", err.Error())
			} else {
				SendNotification("V'ger add task", name)
			}

		}
	}()
}
