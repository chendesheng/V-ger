package subtitles

import (
	"regexp"
	"strings"

	"github.com/peterbourgon/html"
)

func attr2Map(attrs []html.Attribute) (m map[string]string) {
	m = make(map[string]string, len(attrs))
	for _, a := range attrs {
		m[a.Key] = a.Val
	}
	return
}
func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" && strings.Index(a.Val, class) >= 0 {
			return true
		}
	}
	return false
}
func hasId(n *html.Node, id string) bool {
	return hasAttr(n, "id", id)
}
func hasAttr(n *html.Node, key string, val string) bool {
	for _, a := range n.Attr {
		if a.Key == key && a.Val == val {
			return true
		}
	}

	return false
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}

	return ""
}
func getText(n *html.Node) string {
	for _, c := range n.Child {
		if c.Type == html.TextNode {
			return strings.TrimSpace(c.Data)
		}
		if c.Type == html.ElementNode && c.Data == "a" {
			return getText(c)
		}
	}
	return ""
}
func getClass(n *html.Node, class string) []*html.Node {
	list := make([]*html.Node, 0)
	for _, c := range n.Child {
		if hasClass(c, class) {
			list = append(list, c)
		}
	}
	return list
}
func getClass1(n *html.Node, class string) *html.Node {
	for _, c := range n.Child {
		if hasClass(c, class) {
			return c
		}
	}
	return nil
}
func getId(n *html.Node, id string) *html.Node {
	for _, c := range n.Child {
		if hasId(c, id) {
			return c
		}
	}
	return nil
}

func getTag(n *html.Node, tag string) []*html.Node {
	list := make([]*html.Node, 0)
	for _, c := range n.Child {
		if c.Data == tag {
			list = append(list, c)
		}
	}
	return list
}
func getTag1(n *html.Node, tagName string) *html.Node {
	for _, c := range n.Child {
		if c.Data == tagName {
			return c
		}
	}
	return nil
}
func getRidOfTags(n *html.Node) (text string) {
	text = ""
	for _, c := range n.Child {
		if c.Type == html.TextNode {
			text = text + strings.TrimSpace(c.Data)
		}
		if c.Type == html.ElementNode {
			text += getText(c)
		}
	}
	return
}
func getRidOfSpace(text string) string {
	r, _ := regexp.Compile(`[\s\r\n]+`)
	r1, _ := regexp.Compile(`：[\s\r\n]+`)
	return r1.ReplaceAllString(r.ReplaceAllString(text, " "), "：")
}
