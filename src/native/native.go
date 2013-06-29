package native

import (
	// "encoding/utf8"
	// "github.com/mkrautz/objc"
	// . "github.com/mkrautz/objc/AppKit"
	// . "github.com/mkrautz/objc/Foundation"
	"log"
	"os"
	"os/exec"
	"path"
)

var WebSiteAddress string

func init() {
	WebSiteAddress = "127.0.0.1:9527"

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

func SendNotification(title, infoText string) error {
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
	cmd := exec.Command("open", vgerHelper, "--args", "notification", WebSiteAddress, title, infoText)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func MoveFileToTrash(dir, name string) error {
	print("trash file ", name)
	wd, _ := os.Getwd()
	vgerHelper := path.Join(wd, "vgerhelper.app")
	cmd := exec.Command("open", vgerHelper, "--args", "trash", dir, name)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
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

func ConvertEncodingToUTF8(file string, srcEncoding string) {
	cmd := exec.Command("iconv", "-f", srcEncoding, "-t", "utf8", file)

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	cmd.Stdout = f
	cmd.Run()
}
