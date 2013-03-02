package main

import (
	"code.google.com/p/cookiejar"
	"download"
	"fmt"
	"html/template"
	// "io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	// "regexp"
	// "io"
	"runtime"
	// "strings"
	// "encoding/json"
	"b1"
	"log"
	"net/http"
	"net/url"
	"os"
	"subtitles"
	"thunder"
	"time"
)

func init() {
	f, err := os.OpenFile("vger.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	config = readConfig()

	client := &http.Client{
		Jar: cookiejar.NewJar(true),
	}
	cookie := http.Cookie{
		Name:    "gdriveid",
		Value:   config["gdriveid"],
		Domain:  "xunlei.com",
		Expires: time.Now().AddDate(100, 0, 0),
	}
	cookies := []*http.Cookie{&cookie}
	url, _ := url.Parse("http://vip.lixian.xunlei.com")
	client.Jar.SetCookies(url, cookies)

	download.DownloadClient = client
	thunder.Client = client
	subtitles.Client = client
	b1.Client = client

	download.BaseDir = config["dir"]

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func pick(arr []string, emptyMessage string) (int, string) {
	if len(arr) == 0 {
		if emptyMessage != "" {
			fmt.Println(emptyMessage)
		}
		return -1, ""
	}

	for i, item := range arr {
		fmt.Printf("[%d] %s\n", i+1, item)
	}

	next := ""
	i := 0
	fmt.Scanf("%d%s", &i, &next)
	i--
	if i >= 0 && i < len(arr) {
		return i, next
	}
	fmt.Println("pick wrong number.")
	return -1, ""
}
func checkIfSubtitle(input string) bool {
	return !(strings.Contains(input, "://") || strings.HasSuffix(input, ".torrent") || strings.HasPrefix(input, "magnet:"))
}
func checkIfSpeed(input string) (int64, bool) {
	num, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, false
	}
	if num > 10*1024*1024 {
		num = 10 * 1024 * 1024
	}
	return int64(num), true
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("tasks").Parse(`<html>
<head>
	<title>V'ger</title>
	<link href="/assets/style.css" rel="stylesheet" type="text/css"> 
	<script src="http://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
</head>
<body>
	<script type="text/javascript">
		$(document).ready(function() {
			function init() {
				$('.action-play').on('click', function() {
					$.get('/play/' + $(this).data('name'), function() {})
				});

				$('.action-resume-download').on('click', function() {
					$.get('/resume/' + $(this).data('name'), function(resp) {
						get_progress()
						resp && alert(resp);
					})
				});

				$('.action-stop-download').on('click', function() {
					$.get('/stop/' + $(this).data('name'), function(resp) {
						get_progress()
						resp && alert(resp);
					})
				});
				
				$('.action-limit').on('change', function() {
					$.post('/limit/'+$(this).data('name'), {'limit': $(this).val() }, function(resp) {
						get_progress()
						resp && alert(resp);
					})
				})
			}
			init();
			function get_progress() {
				$.get('/progress', function(resp) {
					$('#tasks').html(resp);
					init();
				});
			}
			setInterval(get_progress, 2000)

			$('#new-task').on('click', function() {
				$.post('/new', {'url': $('#new-url').val()}, function(resp) {
					get_progress()
					resp && alert(resp);
				})
			})
			$('#refresh-tasks').on('click', function() {
				get_progress();
			})
		});
	</script>
	<!--for debug-->
	<input type="button" id="refresh-tasks" value="Refresh" style="display:none;" />
	<h1>V'ger</h1>
	<h2>Speed is Fun!</h2>
	<div id="tasks">
		<ul>
		{{range .}}
	       <li>
	       		{{.}}
	       		<div>
	       			<input class="action-play" type="button" value="Open" data-name="{{.Name}}"/>
	       			<input class="action-resume-download" type="button" value="Resume" data-name="{{.Name}}"/>
	       			<input class="action-stop-download" type="button" value="Stop" data-name="{{.Name}}"/>
	       			<select id="limit-{{.NameHash}}" class="action-limit" data-name="{{.Name}}">
	       				<option value="0">No limit</option>
	       				<option value="50">Up to 50K</option>
	       				<option value="100">Up to 100K</option>
	       				<option value="150">Up to 150K</option>
	       				<option value="200">Up to 200K</option>
	       				<option value="300">Up to 300K</option>
	       			</select>
		   			<script type="text/javascript">
		   				$("#limit-{{.NameHash}}").val({{.LimitSpeed}})
		   			</script>
	       		<div>
	       </li>
	    {{end}}
		</ul>
	</div>
	<div id="start-download">
		<div><input type="text" id="new-url" placeHolder="Input anything you want"></div>
		<div><input type="button" id="new-task" value="Start Download"/></div>
	</div>
</body>
</html>`))
	tasks := download.GetTasks()
	download.SortTasksByCreateTime(tasks)
	t.Execute(w, tasks)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("play \"%s\".\n", name)
	cmd := exec.Command("open", fmt.Sprintf("%s%c%s", download.BaseDir, os.PathSeparator, name))
	cmd.Start()

	w.Write([]byte(``))
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[8:])
	fmt.Printf("resume download \"%s\".\n", name)

	w.Write([]byte(download.ResumeDownload(name)))
}
func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")

	if strings.Contains(url, "lixian.vip.xunlei.com") {
		fmt.Printf("add download \"%s\".\n", url)

		w.Write([]byte(download.NewDownload(url)))
		return
	}

	thunder.Login(config["thunder-user"], config["thunder-password"])
	tasks := thunder.NewTask(url)

	arr := make([]string, len(tasks))
	for i, s := range tasks {
		arr[i] = s.String()
	}
	i, next := pick(arr, "")
	if i != -1 {
		selectedTask := tasks[i]
		if selectedTask.Percent < 100 {
			fmt.Println("the task is not ready.")
			return
		}
		if next == "" {
			getMovieSub(selectedTask.Name)
		}

		fmt.Printf("add download \"%s\".\n", selectedTask.DownloadURL)

		w.Write([]byte(download.NewDownload(selectedTask.DownloadURL)))
	}
}

// func newThunderTask(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	url := r.FormValue("url")

// 	thunder.NewTask(url)
// }
func stopHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[6:])
	fmt.Printf("stop download \"%s\".\n", name)

	w.Write([]byte(download.StopDownload(name)))
}
func limitHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := url.QueryUnescape(r.URL.String()[7:])
	r.ParseForm()
	speed := r.FormValue("limit")

	fmt.Printf("download \"%s\" limit speed %dK.\n", name, speed)

	w.Write([]byte(download.LimitSpeed(name, speed)))
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {

}
func progressHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("tasks").Parse(`<ul>
	{{range .}}
       <li>
       		{{.}}
       		<div>
       			<input class="action-play" type="button" value="Open" data-name="{{.Name}}"/>
       			<input class="action-resume-download" type="button" value="Resume" data-name="{{.Name}}"/>
       			<input class="action-stop-download" type="button" value="Stop" data-name="{{.Name}}"/>
	   			<select id="limit-{{.NameHash}}" class="action-limit" data-name="{{.Name}}">
	   				<option value="0">No limit</option>
	   				<option value="50">Up to 50K</option>
	   				<option value="100">Up to 100K</option>
	   				<option value="150">Up to 150K</option>
	   				<option value="200">Up to 200K</option>
	   				<option value="300">Up to 300K</option>
	   			</select>
	   			<script type="text/javascript">
	   				$("#limit-{{.NameHash}}").val({{.LimitSpeed}})
	   			</script>
       		<div>
       </li>
    {{end}}
	</ul>`))

	tasks := download.GetTasks()
	download.SortTasksByCreateTime(tasks)
	t.Execute(w, tasks)
	w.Write([]byte(fmt.Sprintf("<h3>Go routine numbers: %d</h3>", runtime.NumGoroutine())))
}

type command struct {
	ack    chan bool
	result chan string

	name string
	arg  string
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	path := r.URL.Path[1:]
	if _, err := os.OpenFile(path, os.O_RDONLY, 0666); os.IsNotExist(err) {
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, path)
	}
}

func main() {
	// thunder.Login(config["thunder-user"], config["thunder-password"])
	// fmt.Println("thunder login success.")

	download.StartHandleCommands()

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/assets/", assetsHandler)

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/play/", playHandler)
	http.HandleFunc("/resume/", resumeHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/progress", progressHandler)
	http.HandleFunc("/new", newTaskHandler)
	http.HandleFunc("/limit/", limitHandler)

	fmt.Println("server start listern port 3824.")
	err := http.ListenAndServe(":3824", nil)
	if err != nil {
		log.Fatal(err)
	}
}
