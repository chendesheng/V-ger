package cocoa

import (
	"github.com/mkrautz/objc"
	. "github.com/mkrautz/objc/AppKit"
	. "github.com/mkrautz/objc/Foundation"
	"io"
	"log"
	"native"
	"os"
	// "os/exec"
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
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
	SendNotification("Sleep after 60 seconds", reason)

	t := time.NewTimer(time.Second * 60)
	api.shutdownQuit = make(chan bool)
	func() {
		select {
		case <-api.shutdownQuit:
			log.Print("sleep stop")
			t.Stop()
		case <-t.C:
			log.Print("sleep now")
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
			s.InitWithSource("Tell application \"Finder\" to sleep")
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

func (api *cocoaNativeAPI) SetIcon(dir string, r io.Reader) {
	pool := NewNSAutoreleasePool()
	defer pool.Drain()

	goimg, err := jpeg.Decode(r)
	if err != nil {
		println(err.Error())
	}
	if goimg == nil {
		println("nil")
	}
	sz := goimg.Bounds().Size()
	maxw := sz.X
	if sz.X < sz.Y {
		maxw = sz.Y
	}
	newimg := image.NewRGBA(image.Rect(0, 0, maxw, maxw))
	println((maxw - sz.X) / 2)
	rt := goimg.Bounds()
	rt.Min.X = (maxw - sz.X) / 2
	rt.Min.Y = (maxw - sz.Y) / 2
	rt.Max.X += (maxw - sz.X) / 2
	rt.Max.Y += (maxw - sz.Y) / 2
	// rt.Add(image.Point{(maxw - sz.X) / 2, (maxw - sz.Y) / 2})

	draw.Draw(newimg, rt, goimg, image.Point{0, 0}, draw.Src)

	buf := bytes.NewBuffer(nil)
	png.Encode(buf, newimg)
	b, _ := ioutil.ReadAll(buf)
	data := DataWithBytes(b)

	// data.WriteToFile("/Volumes/Data/Downloads/Video/Girls/c.jpg", true)

	img := NSImage{objc.GetClass("NSImage").Alloc()}
	img.AutoRelease()

	img.InitWithData(data)

	// bmp1 := NSBitmapImageRep{img.Representations().ObjectAtIndex(0)}
	// data1 := bmp1.RepresentationUsingType(NSPNGFileType, NSDictionary{objc.NilObject()})
	// // data1.WriteToFile("/Volumes/Data/Downloads/Video/Girls/d.jpg", true)
	// bf := bytes.NewBuffer(data1.Bytes())
	// _, err = png.Decode(bf)
	// if err != nil {
	// 	log.Print(err)
	// }

	w := NSSharedWorkspace()
	w.SetIcon(img, dir, NSExcludeQuickDrawElementsIconCreationOption)
}

func (api *cocoaNativeAPI) GetIcon(dir string, w io.Writer) bool {
	if _, err := os.Stat(path.Join(dir, "Icon\r")); os.IsNotExist(err) {
		return false
	}

	pool := NewNSAutoreleasePool()
	defer pool.Drain()

	ws := NSSharedWorkspace()
	img := ws.IconForFile(dir)
	bmp := NSBitmapImageRep{img.Representations().ObjectAtIndex(0)}
	data := bmp.RepresentationUsingType(NSPNGFileType, NSDictionary{objc.NilObject()})
	w.Write(data.Bytes())
	return true
}
