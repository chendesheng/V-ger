package native

import (
	"cocoa"
	// "encoding/utf8"
	// "github.com/mkrautz/objc"
	// . "github.com/mkrautz/objc/AppKit"
	// . "github.com/mkrautz/objc/Foundation"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"util"
)

func init() {
	// c := objc.NewClass(NotificationCenterDelegate{})
	// c.AddMethod("didActivateNotification:", (*NotificationCenterDelegate).DidActivateNotification)

	// objc.RegisterClass(c)
}

// type NotificationCenterDelegate struct {
// 	objc.Object `objc:"GoNotificationCenterDelegate : NSObject"`
// }

// func (n *NotificationCenterDelegate) DidActivateNotification(notification objc.Object) {
// 	log.Print("DidActivateNotification")
// }
// var uiCh chan cocoa.UICommand

// func SetUIChan(ch chan cocoa.UICommand) {
// 	uiCh = ch
// }
func Start() {
	cocoa.Start()
}

func SendNotification(title, infoText string) error {
	// cocoa.SendNotification(title, infoText)

	// notification := NSUserNotification{objc.GetClass("NSUserNotification").Alloc().Init()}
	// notification.SetTitle(title)
	// notification.SetInformativeText(infoText)
	// notification.SetSoundName(NSUserNotificationDefaultSoundName)
	// notification.SetHasActionButton(true)
	// notification.SetActionButtonTitle("Open")

	// center := NSDefaultUserNotificationCenter()
	// center.DeliverNotification(notification)

	wd, _ := os.Getwd()
	vgerHelper := path.Join(wd, "vgerhelper.app")
	cmd := exec.Command("open", vgerHelper, "--args", "notification", util.ReadConfig("server"), title, infoText)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func MoveFileToTrash(dir, name string) error {
	f, err := os.Open(path.Join(dir, name))
	if err != nil {
		return err
	} else {
		f.Close()
	}

	log.Println("trash file ", name)

	u, err := user.Current()
	if err != nil {
		log.Println(err)
		return err
	}

	trashPath := path.Join(u.HomeDir, ".Trash")
	return os.Rename(path.Join(dir, name), path.Join(trashPath, name))
}

func Shutdown(reason string) error {
	wd, _ := os.Getwd()
	vgerHelper := path.Join(wd, "vgerhelper.app")
	cmd := exec.Command("open", vgerHelper, "--args", "shutdown", reason)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// func ConvertEncodingToUTF8(file string, srcEncoding string) {
// 	cmd := exec.Command("iconv", "-f", srcEncoding, "-t", "utf8", file)

// 	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
// 	if err != nil {
// 		return
// 	}
// 	cmd.Stdout = f
// 	cmd.Run()
// }
