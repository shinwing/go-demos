package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

func main() {
	// routers
	http.Handle("/logs/", http.HandlerFunc(newLogFilesHandler().handerLogFiles))
	http.ListenAndServe(":55555", nil)
}

// LogFilesHandler providers a handler to download log files
type LogFilesHandler struct {
	LogPath string
}

func newLogFilesHandler() *LogFilesHandler {
	// get log file path from os environment
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "/tmp/"
	}
	return &LogFilesHandler{logPath}
}

type dirFileInfo struct {
	FileName string
	FileSize string
}

func (logHandler *LogFilesHandler) handerLogFiles(w http.ResponseWriter, r *http.Request) {
	// Content-Type handling
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err == nil && len(query["dl"]) > 0 {
		w.Header().Set("Content-Type", "application/octet-stream")
		filePathName := logHandler.LogPath + path.Clean(r.URL.Path)[5:]
		needDownloadFile, err := os.Open(filePathName)
		if err != nil {
			http.Error(w, "404 Not Found : Error while opening the file.", 404)
			return
		}
		defer needDownloadFile.Close()

		d, err := needDownloadFile.Stat()
		if err != nil {
			http.Error(w, "404 Not Found : Error while opening the file.", 404)
			return
		}
		// serveContent will check modification time
		http.ServeContent(w, r, d.Name(), d.ModTime(), needDownloadFile)
		return
	}

	dir, err := os.Open(logHandler.LogPath)
	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		fmt.Println(err)
		return
	}
	defer dir.Close()

	allFileInfos := make([]dirFileInfo, 0)
	files, _ := dir.Readdir(-1)
	for _, f := range files {
		// avoid hidden files
		if f.Name()[0] == '.' || f.IsDir() {
			continue
		}
		allFileInfos = append(allFileInfos, dirFileInfo{
			FileName: f.Name(),
			FileSize: strconv.FormatFloat(float64(f.Size())/1024/1024, 'f', 2, 64) + " Mb",
		})
	}

	tpl, err := template.New("tpl").Parse(dirPageTpl)
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		return
	}

	err = tpl.Execute(w, allFileInfos)
	if err != nil {
		fmt.Println(err)
	}
}

const dirPageTpl = `<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="cn">
<head>
<style type="text/css">
a, a:active {text-decoration: none; color: blue;}
a:visited {color: #48468F;}
a:hover, a:focus {text-decoration: underline; color: red;}
body {background-color: #F5F5F5;}
h2 {margin-bottom: 12px;}
table {margin-left: 20;}
th, td { font: 120% monospace; text-align: left;}
th { font-weight: bold; padding-right: 14px; padding-bottom: 3px;}
td {padding-right: 200px;}
td.s, th.s {text-align: right;}
div.list { background-color: white; border-top: 1px solid #646464; border-bottom: 1px solid #646464; padding-top: 10px; padding-bottom: 14px;}
div.foot { font: 120% monospace; color: #787878; padding-top: 4px;}
</style>
</head>
<body>
<div class="list">
<table summary="Directory Listing" cellpadding="0" cellspacing="0">
<thead><tr><th class="dl">Name</th><th class="t">Size</th></tr></thead>
<tbody>
{{range .}}
<tr><td class="dl"><a href="/logs/{{.FileName}}?dl">{{.FileName}}</a></td><td class="t">{{.FileSize}}</td></tr>
{{end}}
</tbody>
</table>
</div>
</body>
</html>`
