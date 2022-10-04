package filesystem

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var output string

type Control struct {
	FirstSpace int64
}

func getDate() string {
	t := time.Now()

	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return fecha
}

func writeBinary(file *os.File, content []byte) bool {
	_, err := file.Write(content)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func Exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func WriteResponse(w http.ResponseWriter, content string) {
	fmt.Fprintln(w, "{\n\"response\" : \""+content+"\"\n}")
}
