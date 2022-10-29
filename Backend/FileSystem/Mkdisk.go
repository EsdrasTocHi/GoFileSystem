package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func MkDisk(size int, unit rune, fit rune, path string, w http.ResponseWriter) {
	if !Exist(path) {

		paux := RemoveFileName(path)

		if paux != "" {
			exec.Command("mkdir", "-p", paux).Run()
			exec.Command("chmod", "-R", "777", paux).Run()
		}

		var tam int64
		if unit == 'k' {
			tam = int64(size * 1024)
		} else {
			tam = int64(size * 1024 * 1024)
		}

		file, err := os.Create(path)

		defer file.Close()

		if err != nil {
			log.Fatal(err)
			return
		}

		var temporal int8 = 0
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, &temporal)

		for i := 0; i < size; i++ {
			if !writeBinary(file, buffer.Bytes()) {
				return
			}
		}

		mbr := structs.Mbr{}

		copy(mbr.Mbr_fecha_creacion[:], []byte(getDate()))
		binary.BigEndian.PutUint64(mbr.Mbr_tamano[:], uint64(tam))
		binary.BigEndian.PutUint64(mbr.Mbr_dsk_fit[:], uint64(fit))
		binary.BigEndian.PutUint64(mbr.Mbr_dsk_signature[:], uint64(time.Now().UnixNano()))
		mbr.Mbr_partition_1 = structs.Partition{}
		mbr.Mbr_partition_2 = structs.Partition{}
		mbr.Mbr_partition_3 = structs.Partition{}
		mbr.Mbr_partition_4 = structs.Partition{}

		file.Seek(0, os.SEEK_SET)

		var buffer2 bytes.Buffer
		err = binary.Write(&buffer2, binary.BigEndian, &mbr)
		if err != nil {
			log.Fatal(err)
			return
		}

		writeBinary(file, buffer2.Bytes())
		WriteResponse(w, "DISK CREATED SUCCESFULLY!")
		return
	}

	WriteResponse(w, "$Error: The disk already exists")
}
