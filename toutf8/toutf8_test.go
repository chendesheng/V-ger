package toutf8

import (
	"download"
	"strings"
	// "fmt"
	"io/ioutil"
	"os"
	// "strings"
	"testing"
)

// func TestGB18030ToUTF8(t *testing.T) {
// 	// data, err := ioutil.ReadFile("gb18030.txt")
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// 	return
// 	// }

// 	res, err := ConverToUTF8("gb18030.txt")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	println(res)
// }

func TestGuessEncoding(t *testing.T) {
	_, name, size, err := download.GetDownloadInfo("http://gdl.lixian.vip.xunlei.com/download?fid=Ss7ITuBtzz5UOYEvjdYqZymfcAd9g64bAAAAABMxNDRgilgayPzxKh8NUdzh0exi&mid=666&threshold=150&tid=4C9C93B69AE817D8EF955724415F360B&srcid=4&verno=1&g=13313434608A581AC8FCF12A1F0D51DCE1D1EC62&scn=t18&i=E17C7338EF7371ADF576B4AAC0A20E060094096E&t=6&ui=119888259&ti=515163683489345&s=464421757&m=0&n=01D69929BDC7E9D5E64F658C3A2E47756113558D3E6E2E53305074D46E2E4368693E748A382E445644135894714143332E5603D4FEC13430304F49D669342E42594F62AB1147424F48341F893476000000&ih=E17C7338EF7371ADF576B4AAC0A20E060094096E&fi=0&pi=515163683423745&ff=0&co=4C9A1ECCE428F5014879142E685583B7&cm=1")
	if err != nil {
		t.Error(err)
	}

	utf8name, encoding, err := ConverToUTF8(strings.NewReader(name))
	if err != nil {
		t.Error(err)
	}

	println(size)
	println(utf8name)
	println(encoding)

	// os.OpenFile("utf16le.txt", flag, perm)
	// data, err := ioutil.ReadFile("gb18030.txt")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// res, err := ConverToUTF8("utf16le.txt")
	// if err != nil {
	// 	t.Error(err)
	// }

	// println(res)
	// println(len(res))

	// infoes, _ := ioutil.ReadDir("/Volumes/Data/Downloads/Video/Girls")
	// for _, f := range infoes {
	// 	println(f.Name())
	// 	if strings.HasPrefix(f.Name(), "Icon") {
	// 		if f.IsDir() {
	// 			println("yes")
	// 		}
	// 		println("size:", f.Size())
	// r, err := os.Open("/Volumes/Data/Downloads/Video/Girls/Icon\r/..namedfork/rsrc")

	// if err != nil {
	// 	println(err)
	// }

	// bytes, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	println(err)
	// }
	// println(len(bytes))
	// r.Close()
	bytes, _ := ioutil.ReadFile("/Volumes/Data/Downloads/Video/Rake/Icon\r/..namedfork/rsrc")
	// fmt.Printf("%v", bytes)
	// bytes, _ := ioutil.ReadFile("/Volumes/Data/Downloads/Video/Rake/b.jpg")
	// ioutil.WriteFile("/Volumes/Data/Downloads/Video/Rake/Icon\r/..namedfork/rsrc", bytes, os.ModeDevice)
	ioutil.WriteFile("/Volumes/Data/Downloads/Video/Rake/a.jpg", bytes, os.ModeType)
	// f, err := os.OpenFile("/Volumes/Data/Downloads/Video/Rake/Icon\r/..namedfork/rsrc", os.O_RDWR)
	// f.Write(bytes)
	// 	}
	// }
}
