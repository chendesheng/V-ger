package subtitles

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"vger/httpex"
	"vger/thunder"
)

type kankanSearch struct {
	url  string
	quit chan struct{}
}

func (k *kankanSearch) search(result chan Subtitle) error {
	log.Print("kankan search:", k.url)

	err := thunder.Login(k.quit)
	if err != nil {
		return err
	}

	regFid := regexp.MustCompile("[?].*fid=([^&]+)")
	if matches := regFid.FindStringSubmatch(k.url); len(matches) > 0 {
		fid := matches[1]
		cid, _, gcid := thunder.ParseFid(fid)
		userid := thunder.GetUserId()

		sourceUrl := fmt.Sprintf("http://i.vod.xunlei.com/subtitle/list?gcid=%s&cid=%s&userid=%s", gcid, cid, userid)

		content, err := httpex.GetStringResp(sourceUrl, nil, k.quit)
		if err != nil {
			return err
		}

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
			case <-k.quit:
				return fmt.Errorf("quit")
			case result <- sub:
			}
		}

		return nil
	} else {
		return fmt.Errorf("no fid in url")
	}
}
