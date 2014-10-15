/*
Copyright Inocybe Technlogies, 2014. All rights reserved.
*/

package main

import (
	"html/template"
	"log"
	"net/http"
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


/* Welcome */

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	pages := []string{"update", "network", "cloudconfig"}
	p := WelcomePage{pages}
	renderTemplate(w, "welcome", p)
}

/* Templates */

var templates = template.Must(template.ParseFiles("templates/welcome.html", "templates/updateserver.html", "templates/networkconfig.html", "templates/cloudconfig.html", "templates/textcloudconfig.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/save/update", saveUpdateHandler)
	http.HandleFunc("/network/", networkHandler)
	http.HandleFunc("/save/network", saveNetworkHandler)
	http.HandleFunc("/cloudconfig", cloudConfigHandler)
	http.HandleFunc("/save/cloudconfig", saveCloudConfigHandler)


	http.ListenAndServe(":8080", nil)
}
