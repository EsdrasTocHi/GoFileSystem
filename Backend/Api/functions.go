package api

import (
	reader "Backend/Analyzer"
	fs "Backend/FileSystem"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HolaMundo(w http.ResponseWriter, r *http.Request) {
	log.Println("hola mundo")
	fmt.Fprintln(w, "hola mundo")
}

func mkdisk(w http.ResponseWriter, r *http.Request) {
	fs.MkDisk(10, 'k', 'f', "disco.dk", w)
}

func ReadCommand(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var data map[string]string

	json.Unmarshal(body, &data)

	if val, ok := data["command"]; ok {
		fmt.Println(val)
		reader.ReadCommand(val, w)
	} else {
		fs.WriteResponse(w, "$Error: something went wrong")
	}
}
