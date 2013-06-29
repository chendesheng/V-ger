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
	c.AddMethod("applicationDidFinishLaunching:", (*AppDelegate).ApplicationDidFinishLaunching)
	c.AddMethod("buttonclick:", (*AppDelegate).IButtonClick)
	c.AddMethod("menuClick:", (*AppDelegate).MenuClick)

	objc.RegisterClass(c)
}

type AppDelegate struct {
	objc.Object `objc:"GOAppDelegate : NSObject"`
	Window      objc.Object `objc:"IBOutlet"`
}

func (delegate *AppDelegate) ApplicationDidFinishLaunching(notification objc.Object) {
	statusItem := NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
	statusItem.SetHighlightMode(true)
	statusItem.SetTarget(delegate.Pointer())
	statusItem.SetAction(objc.GetSelector("menuClick:"))

	go func() {
		// pool := NewNSAutoreleasePool()
		// defer pool.Release()

		i := 0
		for {
			t := time.Tick(time.Second)
			select {
			case <-t:
				i++
				if t, ok := task.GetDownloadingTask(); ok {
					statusItem.SetTitle(fmt.Sprintf("%s %.1f%%", task.CleanName(t.Name), float64(t.DownloadedSize)/float64(t.Size)*100.0))
					statusItem.SetToolTip(fmt.Sprintf("%.2f KB/s %s", t.Speed, t.Est))
				} else {
					statusItem.SetTitle("V'ger")
				}

				break
			}
		}
	}()
}

// func (delegate *AppDelegate) ApplicationDidFinishLaunching(notification objc.Object) {
// 	delegate.statusItem = NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
// 	delegate.statusItem.SetTitle("V'ger")

// 	ScheduledTimerWithTimeInterval(1, delegate.Pointer(), objc.GetSelector("timerTick:"), uintptr(0), true)
// }

// func (delegate *AppDelegate) TimerTick(sender objc.Object) {
// 	// log.Print("tick")

// 	// fmt.Printf("%v\n", delegate.statusItem)
// 	go func() {

// 		defer func() {
// 			if re := recover(); re != nil {
// 				log.Print(re)

// 				delegate.statusItem = NSStatusItem{NSSystemStatusBar().StatusItemWithLength(-1).Retain()}
// 				if t, ok := task.GetDownloadingTask(); ok {
// 					delegate.statusItem.SetTitle(fmt.Sprintf("%s %.1f%%", t.Name, float64(t.DownloadedSize)/float64(t.Size)*100.0))
// 				} else {
// 					delegate.statusItem.SetTitle("V'ger")
// 				}
// 			}
// 		}()

// 		if t, ok := task.GetDownloadingTask(); ok {
// 			delegate.statusItem.SetTitle(fmt.Sprintf("%s %.1f%%", t.Name, float64(t.DownloadedSize)/float64(t.Size)*100.0))
// 		} else {
// 			delegate.statusItem.SetTitle("V'ger")
// 		}
// 	}()
// }
//export IButtonClick
func (delegate *AppDelegate) IButtonClick(sender uintptr) {
	log.Print("clicked")
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

func Start() {
	runtime.LockOSThread()

	pool := NewNSAutoreleasePool()
	defer pool.Release()

	app := NSSharedApplication()
	app.SetDelegate(objc.GetClass("GOAppDelegate").Alloc().Init())
	app.Run()
}
