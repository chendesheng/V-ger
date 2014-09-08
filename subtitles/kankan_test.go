package subtitles

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/cookiejar"
// 	"testing"
// 	"thunder"
// )

// func TestSearchSubtitle(t *testing.T) {
// 	if http.DefaultClient.Jar == nil {
// 		jar, _ := cookiejar.New(nil)
// 		http.DefaultClient.Jar = jar
// 	}

// 	_, _, err := thunder.Login2("5120E7CE422D1E3F34D7ED1501A1C86A", "129697884", "057764593828")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	result := make(chan Subtitle)

// 	go func() {
// 		for s := range result {
// 			fmt.Printf("%#v\n", s)
// 		}
// 	}()

// 	err = kankanSearch(
// 		userid,
// 		"http://gdl.lixian.vip.xunlei.com/download?fid=HK3nXfTE4FGnDdZXX2mmlvBpiQzPnbmRAAAAABSVG7CmQDqOzp7H8lD3rHQz4LZZ&mid=666&threshold=150&tid=A82C68BEB7FD18A62FEE79C5D7B3BF68&srcid=4&verno=1&g=14951BB0A6403A8ECE9EC7F250F7AC7433E0B659&scn=t11&i=815FB379C4C26B89094CE763C5105F83&t=4&ui=119888259&ti=247620789601025&s=2444860879&m=0&n=013543913A2E426C6F0E55CA0C30354530521FD36D30702E482565B271783236344C78A912455253454F5C8F2900000000&ff=0&co=22CCD462C9C00B7DFAEEC5F86DC9126B&cm=1",
// 		result)

// 	close(result)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	return
// }
