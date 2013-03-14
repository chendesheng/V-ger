package main

import (
	// "cocoa"
	// "runtime"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"website"
)

func init() {
	// runtime.LockOSThread()
}
func main() {
	website.Run()

	// cocoa.NSAppRun()
	// count()
}

type word struct {
	str   string
	times int
}
type wordList []word

func (p wordList) Len() int           { return len(p) }
func (p wordList) Less(i, j int) bool { return p[i].times > p[j].times }
func (p wordList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func count() {
	f, _ := os.Open("a.srt")
	data, _ := ioutil.ReadAll(f)
	text := string(data)
	fmt.Println(text)

	reg := regexp.MustCompile("[^A-Za-z \n]")
	text = string(reg.ReplaceAll([]byte(text), []byte("")))
	// fmt.Println(text)
	text = strings.Replace(text, "\n", " ", -1)
	words := strings.Split(text, " ")
	// fmt.Printf("%v", words)
	wordsMap := make(map[string]int)

	for _, w := range words {
		w = strings.ToLower(w)
		if i, ok := wordsMap[w]; ok {
			wordsMap[w] = i + 1
		} else {
			wordsMap[w] = 1
		}
	}
	// fmt.Printf("%v", wordsMap)
	res := make([]word, 0)
	for k, v := range wordsMap {
		res = append(res, word{k, v})
	}
	fmt.Printf("%v", res)

	sort.Sort(wordList(res))

	out, _ := os.OpenFile("w.txt", os.O_CREATE|os.O_WRONLY, 0666)
	for _, w := range res {
		out.WriteString(fmt.Sprintln(w.str, " ", w.times))
	}
	out.Close()
}
