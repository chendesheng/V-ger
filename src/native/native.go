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
	cmd := exec.Command("open", vgerHelper, "--args", WebSiteAddress, title, infoText)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
