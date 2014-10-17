package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"bytes"
)

/* Update Server Configuration */

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
