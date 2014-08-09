package subtitles

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"thunder"
)

func kankanSearch(fileurl string, result chan Subtitle, quit chan struct{}) error {
	log.Print("kankan search:", fileurl)
	regFid := regexp.MustCompile("[?].*fid=([^&]+)")
	if matches := regFid.FindStringSubmatch(fileurl); len(matches) > 0 {
		fid := matches[1]
		cid, _, gcid := thunder.ParseFid(fid)
		userid := thunder.GetUserId()

		sourceUrl := fmt.Sprintf("http://i.vod.xunlei.com/subtitle/list?gcid=%s&cid=%s&userid=%s", gcid, cid, userid)

		content, err := sendGet(sourceUrl, nil)
		if err != nil {
			return err
		}

		println(content)

		v := make(map[string]interface{})
		if err := json.Unmarshal([]byte(content), &v); err != nil {
			return err
		}

		for _, s := range v["sublist"].([]interface{}) {
			m := s.(map[string]interface{})

			sub := Subtitle{}
			sub.URL = m["surl"].(string)
			sub.Description = m["sname"].(string)
			sub.Source = "Kankan"

			select {
			case result <- sub:
				break
			case <-quit:
				return nil
			}
		}

		return nil
	} else {
		return fmt.Errorf("no fid in url")
	}
}
