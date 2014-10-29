package native

import (
	"io"
)

type NotificationClickCallback func(string, string)

type NativeAPI interface {
	SendNotification(title, infoText string)
	SetNotificationClickCallback(callbock NotificationClickCallback)
	ComputerShutdown(reason string) error
	MoveFileToTrash(dir, name string) error
	SetIcon(dir string, r io.Reader)
	GetIcon(dir string, w io.Writer) bool //return true means folder use custom icon. If folder not exists return false.
}

var DefaultNativeAPI NativeAPI

func SendNotification(title, infoText string) {
//	if DefaultNativeAPI != nil {
//		DefaultNativeAPI.SendNotification(title, infoText)
//	}
}

func Shutdown(reason string) error {
//	if DefaultNativeAPI != nil {
//		return DefaultNativeAPI.ComputerShutdown(reason)
//	}
	return nil
}

func MoveFileToTrash(dir, name string) error {
//	if DefaultNativeAPI != nil {
//		return DefaultNativeAPI.MoveFileToTrash(dir, name)
//	}
//
	return nil
}

func SetNotificationClickCallback(callback NotificationClickCallback) {
//	if DefaultNativeAPI != nil {
//		DefaultNativeAPI.SetNotificationClickCallback(callback)
//	}
}
