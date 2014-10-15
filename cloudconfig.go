package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"io"
	"fmt"
	"time"
	"crypto/md5"
	"strconv"
	"strings"
)

type CloudConfigPage struct {
	Body string
}

/* cloud-config upload */

var cloudconfigFilename string = "testconf/cloud-config.yml"

func loadCloudConfig() (*CloudConfigPage, error){
	body, err := ioutil.ReadFile(cloudconfigFilename)
	if err != nil {
		log.Println("ERROR reading server config file")
		log.Println(err)
		return nil, err
	}
	safe := template.HTMLEscapeString(string(body))
	newbody := strings.Replace(safe, "\n", "<br>", -1)

	return &CloudConfigPage{newbody}, nil
}

func allinonecloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		log.Println("Saving uploaded file")
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		fmt.Fprintf(w, "File Uploaded")
		f, err := os.OpenFile("test"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func cloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("cloud-config method:", r.Method)
	var p *CloudConfigPage
	p, err := loadCloudConfig()
	if err != nil {
		p = &CloudConfigPage{"None"}
	}
	renderTemplate(w, "cloudconfig", p)
}

func saveCloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Saving uploaded file")
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	check(err)
	defer file.Close()
	//fmt.Fprintf(w, "%v", handler.Header)
	fmt.Fprintf(w, "File Uploaded")
	f, err := os.Create(cloudconfigFilename)
	check(err)
	defer f.Close()
	io.Copy(f, file)
}
