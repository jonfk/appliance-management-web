/*
Copyright Inocybe Technlogies, 2014. All rights reserved.
*/

package main

import (
	//"os"
	//"encoding/json"
	"html/template"
	"log"
	"net/http"
	"flag"
	"strconv"
)

type WelcomePage struct {
	Pages []string
}

type ServerPage struct {
	Group string
	Server string
}

type NetworkPage struct {
	Dns string
	Address string
	Gateway string
}

func check(err error) {
    if err != nil {
	    log.Println(err)
    }
}

/* Flags */
var port int

var test bool

func init() {
	flag.IntVar(&port, "port", 8080, "set port to serve web app")
	flag.BoolVar(&test, "test", false, "set app to testing mode")
}

/* config files to modify */
var cloudconfigFilename string = "testconf/cloud-config.yml"
// TODO change to "/etc/systemd/network/static.network"
var networkConfigFilename string = "testconf/networktestconfig.network"
// TODO change to "/etc/coreos/update.conf"
var serverConfigFilename string = "testconf/testconfig.conf"

/* Welcome */

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	pages := []string{"update", "network", "cloudconfig"}
	p := WelcomePage{pages}
	renderTemplate(w, "welcome", p)
}


func main() {
	flag.Parse()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/save/update", saveUpdateHandler)
	http.HandleFunc("/network/", networkHandler)
	http.HandleFunc("/save/network", saveNetworkHandler)
	http.HandleFunc("/cloudconfig/", cloudConfigHandler)
	http.HandleFunc("/save/cloudconfig", saveCloudConfigHandler)

	/* Templates */
	var templatesPath string = "/usr/share/appliance-manager/templates/"

	if test {
		log.Println("Testing mode")
		templatesPath = "templates/"

		cloudconfigFilename = "testconf/cloud-config.yml"
		networkConfigFilename = "testconf/networktestconfig.network"
		serverConfigFilename = "testconf/testconfig.conf"
	} else {
		cloudconfigFilename = "/var/lib/coreos-install/user_data"
		networkConfigFilename = "/etc/systemd/network/static.network"
		serverConfigFilename = "/etc/coreos/update.conf"
	}
	log.Printf("Config files modified: \n%s\n%s\n%s", cloudconfigFilename, networkConfigFilename, serverConfigFilename)

	templates = template.Must(template.ParseFiles(templatesPath + "welcome.html", templatesPath+"updateserver.html", templatesPath+"networkconfig.html", templatesPath+"cloudconfig.html", templatesPath+"textcloudconfig.html"))

	addr := "0.0.0.0:" + strconv.Itoa(port)
	log.Println("Listening on : " + addr)
	http.ListenAndServe(addr , nil)
}



var templates *template.Template

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
