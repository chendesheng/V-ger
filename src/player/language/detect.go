package language

//only support Chinese & English for now  (can't detect German)
func DetectLanguages(str string) (string, string) {
	var chsCnt, chtCnt, asciiCnt, anyCnt float64
	for _, r := range str {
		if (r >= '0' && r <= '9') || r == '\r' || r == '\n' || r == '\t' || r == ' ' ||
			r == '[' || r == ']' || r == ';' || r == '-' || r == '>' || r == ':' {
			continue //ignore all numbers and space
		}

		if r < 128 && r > -128 {
			asciiCnt++
		} else if r >= '\u4E00' && r <= '\u9FA5' {
			if _, ok := tables[r]; ok {
				chtCnt++
			} else {
				chsCnt++
			}
		}

		anyCnt++
	}
	asciiCnt /= 2
	anyCnt -= asciiCnt
	println("en:", asciiCnt, "chs:", chsCnt, "cht:", chtCnt, "all:", anyCnt)

	if asciiCnt/anyCnt > 0.7 {
		return "en", ""
	}

	if (chsCnt+chtCnt)/anyCnt > 0.7 {
		if chtCnt/anyCnt > 0.2 {
			return "cht", ""
		} else {
			return "chs", ""
		}
	}

	if (asciiCnt+chsCnt+chtCnt)/anyCnt > 0.8 {
		if chtCnt/(chsCnt+chtCnt) > 0.1 {
			return "en", "cht"
		} else {
			return "en", "chs"
		}
	}

	return "", ""
}
