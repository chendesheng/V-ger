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
	"task"
	"time"
)

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
		cmd := exec.Command("open", path.Join(task.BaseDir, t.Name))
		cmd.Start()
	} else {
		cmd := exec.Command("open", "/Applications/V'ger.app")
		cmd.Start()
	}
}

// func (delegate *AppDelegate) DidActivateNotification(notification objc.Object) {
// 	log.Print("DidActivateNotification")
// }

type uiCommand struct {
	name      string
	arguments interface{}
}

func goAppStarted(chUI chan uiCommand) {
	go func(chUI chan uiCommand) {
		for {
			t := time.Tick(time.Second)
			select {
			case <-t:
				var properties []string
				if t, ok := task.GetDownloadingTask(); ok {
					properties = []string{fmt.Sprintf("%s %.1f%%", task.CleanName(t.Name), float64(t.DownloadedSize)/float64(t.Size)*100.0), fmt.Sprintf("%.2f KB/s %s", t.Speed, t.Est)}
				} else {
					properties = []string{"V'ger"}
				}

				chUI <- uiCommand{"statusItem", properties}

				break
			}
		}
	}(chUI)
}

var chUI chan uiCommand

// func SendNotification(title string, infoText string) {
// 	chUI <- uiCommand{"sendNotification", []string{title, infoText}}
// }
func TrashFile(dir string, name string) {
	chUI <- uiCommand{"trashFile", []string{dir, name}}
}

func Start() {
	runtime.LockOSThread()

	pool := NewNSAutoreleasePool()

	InstallNSBundleHook()

	delegate := objc.GetClass("GOAppDelegate").Alloc().Init()

	app := NSSharedApplication()

	NSDefaultUserNotificationCenter().SetDelegate(delegate)

	statusItem := NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	statusItem.SetHighlightMode(true)
	statusItem.SetTarget(delegate.Pointer())
	statusItem.SetAction(objc.GetSelector("menuClick:"))
	statusItem.SetTitle("V'ger")

	chUI = make(chan uiCommand)

	goAppStarted(chUI)

	for {
		pool.Release()
		pool = NewNSAutoreleasePool()

		event := app.NextEventMatchingMask(0xffffff, NSDateWithTimeIntervalSinceNow(0.05), "kCFRunLoopDefaultMode", true)

		app.SendEvent(event)
		app.UpdateWindows()

		t := time.After(time.Millisecond * 5)
		select {
		case cmd := <-chUI:
			switch cmd.name {
			case "statusItem":
				prop := cmd.arguments.([]string)

				statusItem.SetTitle(prop[0])

				if len(prop) > 1 {
					statusItem.SetToolTip(prop[1])
				}
				break
			// case "sendNotification":
			// 	args := cmd.arguments.([]string)
			// 	title := args[0]
			// 	infoText := args[1]

			// 	notification := NSUserNotification{objc.GetClass("NSUserNotification").Alloc().Init()}
			// 	notification.SetTitle(title)
			// 	notification.SetInformativeText(infoText)
			// 	notification.SetSoundName(NSUserNotificationDefaultSoundName)
			// 	notification.SetHasActionButton(true)
			// 	notification.SetActionButtonTitle("Open")

			// 	center := NSDefaultUserNotificationCenter()
			// 	center.DeliverNotification(notification)

			// 	break
			case "trashFile":
				prop := cmd.arguments.([]string)
				NSTrashFile(prop[0], prop[1])
			default:
				log.Printf("unknown cmd %v", cmd)
				break
			}
			break
		case <-t:
			break
		}
	}

	pool.Release()
}
