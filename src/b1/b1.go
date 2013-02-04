package b1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var Client *http.Client

func Extract(path string) []string {
	resp, err := postFile(path, "http://b1.org/rest/online/upload")
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//here is a bug, Client.Do() not set resp cookie
	url, _ := url.Parse("http://b1.org/")
	Client.Jar.SetCookies(url, resp.Cookies())

	text := readBody(resp)
	regJson, _ := regexp.Compile(`result.listing[ ]?=[ ]?({.*})`)
	jsonStr := regJson.FindStringSubmatch(text)[1]
	var filetree map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &filetree)

	files := make([]string, 0, 10)
	walkFileTree(filetree, "http://b1.org/rest/online/download", &files)
	return files
}

func walkFileTree(filetree map[string]interface{}, currentPath string, files *[]string) {
	children, ok := filetree["children"]
	if ok {
		for _, c := range children.([]interface{}) {
			walkFileTree(c.(map[string]interface{}),
				fmt.Sprintf("%s/%s", currentPath, filetree["name"].(string)),
				files)
		}
	} else {
		*files = append(*files, fmt.Sprintf("%s/%s", currentPath, filetree["name"].(string)))
	}
}
func getExt(path string) string {
	return path[strings.LastIndex(path, ".")+1:]
}
func readBody(resp *http.Response) string {
	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dumpBytes))

	text := string(bytes)
	return text
}

//copy from:  https://groups.google.com/forum/?fromgroups=#!msg/golang-nuts/Zjg5l4nKcQ0/jwKOZ7sb298J
//author: James
func postFile(filename string, target_url string) (*http.Response, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer 
	_, err := body_writer.CreateFormFile("upfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body 
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	// need to know the boundary to properly close the part myself. 
	boundary := body_writer.Boundary()
	// close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until writing to the socket buffer. 
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length 
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	return Client.Do(req)
}
