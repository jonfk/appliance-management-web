package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"bytes"
)

/* Network Configuration */
// TODO change to "/etc/systemd/network/static.network"
var networkConfigFilename string = "testconf/networktestconfig.network"

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
