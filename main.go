/*
Copyright Inocybe Technlogies, 2014. All rights reserved.
*/

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"os"
	"io"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

type ServerPage struct {
	Group string
	Server string
}

type NetworkPage struct {
	Dns string
	Address string
	Gateway string
}

type CloudConfigPage struct {
	Body string
}

func check(err error) {
    if err != nil {
	    log.Println(err)
    }
}

/* Update Server Configuration */

// TODO change to "/etc/coreos/update.conf"
//"/usr/share/coreos/update.conf"
//"/usr/share/coreos/release"
var serverConfigFilename string = "testconfig.conf"

func loadUpdateServerConfig() (*ServerPage, error){
	body, err := ioutil.ReadFile(serverConfigFilename)
	if err != nil {
		log.Println("ERROR reading server config file")
		log.Println(err)
		return nil, err
	}

	lines := strings.Split(string(body), "\n")

	var server, group string = "", ""

	for _, line := range lines {
		arguments := strings.Split(line, "=")
		switch arguments[0] {
		case "GROUP":
			group = arguments[1]
		case "SERVER":
			server = arguments[1]
		}
	}

	return &ServerPage{Group: group, Server: server}, nil
}

func writeUpdateServerConfig(page *ServerPage) {
	var body bytes.Buffer
	_, err := body.WriteString("GROUP=" + page.Group)
	check(err)
	_, err = body.WriteString("\n")
	check(err)
	_, err = body.WriteString("SERVER=" + page.Server)
	check(err)

	ioutil.WriteFile(serverConfigFilename, body.Bytes(), 0644)
}

func saveUpdateHandler(w http.ResponseWriter, r *http.Request) {
	group := r.FormValue("group")
	server := r.FormValue("server")
	log.Println("Set group: " + group)
	log.Println("Set server: " + server)

	writeUpdateServerConfig(&ServerPage{group, server})

	http.Redirect(w, r, "/update/", http.StatusFound)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	p,_ := loadUpdateServerConfig()

	renderTemplate(w, "updateserver", p)
}

/* Network Configuration */
// TODO change to "/etc/systemd/network/static.network"
var networkConfigFilename string = "networktestconfig.network"

func writeNetworkConfig(page *NetworkPage) {
	var body bytes.Buffer
	_, err := body.WriteString("[Match]\n")
	check(err)
	_, err = body.WriteString("Name=enp2s0\n\n")
	check(err)
	_, err = body.WriteString("[Network]\n")
	check(err)
	_, err = body.WriteString("DNS=" + page.Dns)
	check(err)
	_, err = body.WriteString("\n")
	check(err)
	_, err = body.WriteString("Address=" + page.Address)
	check(err)
	_, err = body.WriteString("\n")
	check(err)
	_, err = body.WriteString("Gateway=" + page.Gateway)
	check(err)

	ioutil.WriteFile(networkConfigFilename, body.Bytes(), 0644)
}

func saveNetworkHandler(w http.ResponseWriter, r *http.Request) {
	dns := r.FormValue("dns")
	address := r.FormValue("address")
	gateway := r.FormValue("gateway")
	log.Println("Set dns: " + dns)
	log.Println("Set address: " + address)
	log.Println("Set gateway: " + gateway)

	writeNetworkConfig(&NetworkPage{dns, address, gateway})

	http.Redirect(w, r, "/network/", http.StatusFound)
}

func networkHandler(w http.ResponseWriter, r *http.Request) {
	p := NetworkPage{"", "", ""}
	renderTemplate(w, "networkconfig", p)
}

/* cloud-config upload */

var cloudconfigFilename string = "cloud-config.yml"

func writeCloudConfig(file multipart.File) {
	body, err := ioutil.ReadAll(file)
	check(err)

	ioutil.WriteFile(cloudconfigFilename, body, 0644)
}

func loadCloudConfig() (*CloudConfigPage, error){
	body, err := ioutil.ReadFile(serverConfigFilename)
	if err != nil {
		log.Println("ERROR reading server config file")
		log.Println(err)
		return nil, err
	}

	return &CloudConfigPage{string(body)}, nil
}

func cloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	p := CloudConfigPage{""}
	renderTemplate(w, "cloudconfig", p)
}

func saveCloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	check(err)

	file, fileHeader, err := r.FormFile("file")
	check(err)
	defer file.Close()

	log.Println("\n\nfilename : " + fileHeader.Filename + "\n\n")

	if filepath.Ext(fileHeader.Filename) == "yml" {
		newFile, err := os.OpenFile(cloudconfigFilename, os.O_WRONLY|os.O_CREATE, 0666)
		check(err)
		defer newFile.Close()
		io.Copy(newFile, file)
	} else {
		log.Println("Incorrect filen extension. File was not a yml file")
	}

	http.Redirect(w, r, "/cloudconfig/", http.StatusFound)
}

/* text cloud-config */

func writeTextCloudConfig(body string) {
	ioutil.WriteFile(cloudconfigFilename, []byte(body), 0644)
}

func textcloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	p := CloudConfigPage{"test"}
	renderTemplate(w, "textcloudconfig", p)
}

func savetextCloudConfigHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	log.Println("Set cloud-config: " + body)

	writeTextCloudConfig(body)

	http.Redirect(w, r, "/cloudconfig/", http.StatusFound)
}

/* Templates */

var templates = template.Must(template.ParseFiles("updateserver.html", "networkconfig.html", "cloudconfig.html", "textcloudconfig.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/save/update", saveUpdateHandler)
	http.HandleFunc("/network/", networkHandler)
	http.HandleFunc("/save/network", saveNetworkHandler)
	http.HandleFunc("/cloudconfig/", textcloudConfigHandler)
	http.HandleFunc("/save/cloudconfig/", savetextCloudConfigHandler)

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}
