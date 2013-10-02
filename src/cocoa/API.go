package cocoa

import (
	// "github.com/mkrautz/objc"
	// . "github.com/mkrautz/objc/AppKit"
	. "github.com/mkrautz/objc/Foundation"
	"log"
	"native"
	"os"
	// "os/exec"
	"os/user"
	"path"
	"strings"
	"time"
	// "util"
)

type cocoaNativeAPI struct {
	shutdownQuit chan bool
}

var notificationClickCallback native.NotificationClickCallback

func (api *cocoaNativeAPI) SendNotification(title, infoText string) {
	SendNotification(title, infoText)
}
func (api *cocoaNativeAPI) SetNotificationClickCallback(callbock native.NotificationClickCallback) {
	notificationClickCallback = callbock
}

func (api *cocoaNativeAPI) ComputerShutdown(reason string) error {
	SendNotification("Shutdown after 60 seconds", reason)

	t := time.NewTimer(time.Second * 60)
	api.shutdownQuit = make(chan bool)
	func() {
		select {
		case <-api.shutdownQuit:
			log.Print("shutdown stop")
			t.Stop()
		case <-t.C:
			log.Print("shutdown now")
			/*    NSAppleScript* script = [[NSAppleScript alloc] initWithSource:
			                            @"Tell application \"System Events\" to shut down"];
			if (script != NULL)
			{
			    NSDictionary* errDict = NULL;
			    // execution of the following line ends with EXC
			    if (YES == [script compileAndReturnError: &errDict])
			    {
			        NSLog(@"compiled the script");
			        [script executeAndReturnError: &errDict];
			    }
			    [script release];
			}
			[NSApp terminate: nil];*/
			s := NewNSAppleScript()
			s.InitWithSource("Tell application \"System Events\" to shut down")
			s.CompileAndReturnError()
			s.ExecuteAndReturnError()
		}
		api.shutdownQuit = nil
	}()

	return nil
}
func (api *cocoaNativeAPI) MoveFileToTrash(dir, name string) error {
	f, err := os.Open(path.Join(dir, name))
	if err != nil {
		return err
	} else {
		f.Close()
	}

	// log.Println("trash file ", name)

	u, err := user.Current()
	if err != nil {
		log.Println(err)
		return err
	}

	// print(u.Uid)
	trashPath := path.Join(u.HomeDir, ".Trash")
	if strings.HasPrefix(dir, "/Volumes") {
		strs := strings.SplitN(dir, "/", 4)
		if len(strs) >= 3 {
			trashPath = "/" + path.Join(strs[1], strs[2], ".Trashes", u.Uid)
		} else {
			log.Println("Error external volumes directory.")
		}
	}

	println(path.Join(dir, name))
	println(path.Join(trashPath, name))
	return os.Rename(path.Join(dir, name), path.Join(trashPath, name))
}
