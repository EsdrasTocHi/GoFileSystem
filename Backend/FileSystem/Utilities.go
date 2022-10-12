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
	"strings"
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

func ReadDirBlock(ebr *structs.Dirblock, file *os.File) {
	buffer := bytes.NewBuffer(readNextBytes(file, int64(binary.Size(*ebr))))
	binary.Read(buffer, binary.BigEndian, ebr)
}

func ReadFileBlock(ebr *structs.FileBlock, file *os.File) {
	buffer := bytes.NewBuffer(readNextBytes(file, int64(binary.Size(*ebr))))
	binary.Read(buffer, binary.BigEndian, ebr)
}

func ReadInode(ebr *structs.Inode, file *os.File) {
	buffer := bytes.NewBuffer(readNextBytes(file, int64(binary.Size(*ebr))))
	binary.Read(buffer, binary.BigEndian, ebr)
}

func ReadSuperBlock(ebr *structs.SuperBlock, file *os.File) {
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

func SplithPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func SearchInDirBlocks(pointer int64, file *os.File, name string, istart int64, bstart int64, p *int) structs.Inode {
	var dirBlock structs.Dirblock
	var res structs.Inode
	res.I_type = 'n'

	if pointer == -1 {
		return res
	}

	pointer = bstart + (pointer * 64)
	file.Seek(pointer, os.SEEK_SET)
	ReadDirBlock(&dirBlock, file)
	for i := 0; i < 4; i++ {
		if ToString(dirBlock.B_content[i].B_name[:]) == name {
			file.Seek(istart+(dirBlock.B_content[i].B_inodo*int64(binary.Size(res))), os.SEEK_SET)
			ReadInode(&res, file)
			*p = int(dirBlock.B_content[i].B_inodo)
			return res
		}
	}
	return res
}

func ReadInFileBLocks(pointer int64, file *os.File, istart int64, bstart int64) string {
	var fileBlock structs.FileBlock
	res := ""

	if pointer == -1 {
		return res
	}

	pointer = bstart + (pointer * 64)
	file.Seek(pointer, os.SEEK_SET)
	ReadFileBlock(&fileBlock, file)
	res += ToString(fileBlock.B_content[:])
	return res
}

func SearchFile(file *os.File, inode structs.Inode, path []string, istart int64, bstart int64, p *int) structs.Inode {
	var res structs.Inode
	res.I_type = 'n'

	for i := 0; i < len(path); i++ {
		if inode.I_type == byte('0') {
			for j := 0; j < 16; j++ {
				res = SearchInDirBlocks(inode.I_block[j], file, path[i], istart, bstart, p)

				if res.I_type != 'n' {
					inode = res
					break
				}
			}

			if res.I_type == 'n' {
				break
			}
		} else {
			res.I_type = 'n'
			return res
		}
	}
	return res
}

func ReadFile(file *os.File, inode structs.Inode, istart int64, bstart int64, w http.ResponseWriter) string {
	res := ""

	if inode.I_type == '1' {
		for j := 0; j < 16; j++ {
			res += ReadInFileBLocks(inode.I_block[j], file, istart, bstart)
		}
	} else {
		res = "$Error: yo can not read a directory"
		WriteResponse(w, "$Error: yo can not read a directory")
		return res
	}
	return res
}

func GetLines(content string) []string {
	return strings.Split(content, "\n")
}

func IsUser(line string) bool {
	counter := 0
	for i := 0; i < len(line); i++ {
		if line[i] == ',' {
			counter++
		}
	}

	return counter == 4
}
