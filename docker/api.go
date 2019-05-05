package main

import (
	"encoding/xml"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type InvalidIDCard struct {
	XMLName     xml.Name `xml:"doklady_neplatne"`
	Text        string   `xml:",chardata"`
	PoslZmena   string   `xml:"posl_zmena,attr"`
	PristiZmeny string   `xml:"pristi_zmeny,attr"`
	Dotaz       struct {
		Text  string `xml:",chardata"`
		Typ   string `xml:"typ,attr"`
		Cislo string `xml:"cislo,attr"`
		Serie string `xml:"serie,attr"`
	} `xml:"dotaz"`
	Odpoved struct {
		Text          string `xml:",chardata"`
		Aktualizovano string `xml:"aktualizovano,attr"`
		Evidovano     string `xml:"evidovano,attr"`
		EvidovanoOd   string `xml:"evidovano_od,attr"`
	} `xml:"odpoved"`
}

func isValidIDCard(document_id string, document_type string) bool {
	url := "https://aplikace.mvcr.cz/neplatne-doklady/Doklady.aspx?dotaz=" + document_id + "&doklad=" + document_type
	var contents []byte
	var id InvalidIDCard

	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(contents, &id)
	if err != nil {
		panic(err)
	}

	if id.Odpoved.Evidovano == "ne" {
		return true
	}

	return false
}

func IDCardChecker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	if isValidIDCard(string(vars["id"]), "0") {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			panic(err)
		}
	} else {
		_, err := w.Write([]byte("NOK"))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/validate/{id}", IDCardChecker)
	log.Fatal(http.ListenAndServe(":80", router))
}
