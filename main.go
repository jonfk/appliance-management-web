/*
Copyright Inocybe Technlogies, 2014. All rights reserved.
*/

package main

import (
	"os"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Configuration struct {
	Port string
	TemplatesPath string
}

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

var configuration Configuration


/* Welcome */

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	pages := []string{"update", "network", "cloudconfig"}
	p := WelcomePage{pages}
	renderTemplate(w, "welcome", p)
}


func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/save/update", saveUpdateHandler)
	http.HandleFunc("/network/", networkHandler)
	http.HandleFunc("/save/network", saveNetworkHandler)
	http.HandleFunc("/cloudconfig", cloudConfigHandler)
	http.HandleFunc("/save/cloudconfig", saveCloudConfigHandler)

	file, _ := os.Open("conf/conf.json")
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	check(err)

	log.Printf("Starting appliance-management-web with configuration: %#v\n",configuration)

	http.ListenAndServe("0.0.0.0:"+configuration.Port, nil)
}

/* Templates */

var templatesPath string = "/usr/share/appliance-manager/templates/"

var templates = template.Must(template.ParseFiles(templatesPath + "welcome.html", templatesPath+"updateserver.html", templatesPath+"networkconfig.html", templatesPath+"cloudconfig.html", templatesPath+"textcloudconfig.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
