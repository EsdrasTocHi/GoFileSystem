package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
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

func readNextBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
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

func ReadMbr(mbr *structs.Mbr, file *os.File) {
	file.Seek(0, os.SEEK_SET)
	var size int64 = int64(binary.Size(*mbr))
	buffer := bytes.NewBuffer(readNextBytes(file, size))
	binary.Read(buffer, binary.BigEndian, mbr)
}

func ReadEbr(ebr *structs.Ebr, file *os.File) {
	buffer := bytes.NewBuffer(readNextBytes(file, int64(binary.Size(*ebr))))
	binary.Read(buffer, binary.BigEndian, ebr)
}

func ToInt(data []byte) int64 {
	res := binary.BigEndian.Uint64(data)

	return int64(res)
}

func ToString(data []byte) string {
	aux := false
	res := ""
	for i := len(data) - 1; i >= 0; i-- {
		if !aux {
			if data[i] != 0 {
				res = string(data[i]) + res
				aux = true
			}
		} else {
			res = string(data[i]) + res
		}
	}

	return res
}

func ToRune(data []byte) rune {
	num := int32(ToInt(data))

	return rune(num)
}
