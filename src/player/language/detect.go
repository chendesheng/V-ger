package language

//only support Chinese & English for now  (can't detect German)
func DetectLanguages(str string) (string, string) {
	var chineseCnt, asciiCnt, anyCnt float64
	for _, r := range str {
		if (r >= '0' && r <= '9') || r == '\r' || r == '\n' || r == '\t' || r == ' ' ||
			r == '[' || r == ']' || r == ';' || r == '-' || r == '>' || r == ':' {
			continue //ignore all numbers and space
		}

		if r < 128 && r > -128 {
			asciiCnt++
		} else if r >= '\u4E00' && r <= '\u9FA5' {
			chineseCnt++
		}

		anyCnt++
	}
	asciiCnt /= 2
	anyCnt -= asciiCnt
	println("en:", asciiCnt, "cn:", chineseCnt, "all:", anyCnt)

	if asciiCnt/anyCnt > 0.7 {
		return "en", ""
	}

	if chineseCnt/anyCnt > 0.7 {
		return "zh", ""
	}

	if (asciiCnt+chineseCnt)/anyCnt > 0.8 {
		return "en", "zh"
	}

	return "", ""
}
