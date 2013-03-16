package native

import (
	"log"
	"os"
	"os/exec"
	"path"
)

var WebSiteAddress string

func init() {
	WebSiteAddress = "127.0.0.1:9527"
}

func SendNotification(title, infoText string) error {
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
