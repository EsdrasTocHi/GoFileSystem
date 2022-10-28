package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
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

func ReadByte(ebr *byte, file *os.File) {
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

func SearchInDirBlocks(pointer int64, file *os.File, name string, istart int64, bstart int64, p *int64) structs.Inode {
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
			file.Seek(istart+int64(dirBlock.B_content[i].B_inodo*int32(binary.Size(res))), os.SEEK_SET)
			ReadInode(&res, file)
			*p = int64(dirBlock.B_content[i].B_inodo)
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

func SearchFile(file *os.File, inode structs.Inode, path []string, istart int64, bstart int64, p *int64) structs.Inode {
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

func CountGroups(c string) int {
	content := GetLines(c)
	counter := 0

	for i := 0; i < len(content); i++ {
		if !IsUser(content[i]) {
			counter++
		}
	}

	return counter
}

func CountUsers(c string) int {
	content := GetLines(c)
	counter := 0

	for i := 0; i < len(content); i++ {
		if IsUser(content[i]) {
			counter++
		}
	}

	return counter
}

func ExistGroup(name string, c string) bool {
	content := GetLines(c)
	for i := 0; i < len(content); i++ {
		if !IsUser(content[i]) {
			nameGroup := ""
			aux := ""
			counter := 0
			for j := 0; j < len(content[i]); j++ {
				if content[i][j] == ',' || j == len(content[i])-1 {
					if counter == 2 {
						aux += string(content[i][j])
						nameGroup = aux
					}
					aux = ""
					counter++
					continue
				}
				aux += string(content[i][j])
			}

			if nameGroup == name && content[i][0] != '0' {
				return true
			}
		}
	}
	return false
}

func ExistUser(name string, c string) bool {
	content := GetLines(c)
	for i := 0; i < len(content); i++ {
		if IsUser(content[i]) {
			nameGroup := ""
			aux := ""
			counter := 0
			for j := 0; j < len(content[i]); j++ {
				if content[i][j] == ',' {
					if counter == 3 {
						nameGroup = aux
					}
					aux = ""
					counter++
					continue
				}
				aux += string(content[i][j])
			}

			if nameGroup == name && content[i][0] != '0' {
				return true
			}
		}
	}
	return false
}

func GetGroupId(name string, c string) int {
	content := GetLines(c)
	for i := 0; i < len(content); i++ {
		if !IsUser(content[i]) {
			nameGroup := ""
			aux := ""
			counter := 0
			id := 0
			for j := 0; j < len(content[i]); j++ {
				if content[i][j] == ',' || j == len(content[i])-1 {
					if counter == 2 {
						aux += string(content[i][j])
						nameGroup = aux
					} else if counter == 0 {
						id, _ = strconv.Atoi(aux)
					}
					aux = ""
					counter++
					continue
				}
				aux += string(content[i][j])
			}

			if nameGroup == name && content[i][0] != '0' {
				return id
			}
		}
	}
	return -1
}

func SeparateContent(content *string) string {
	res := ""
	aux := ""
	for i := 0; i < len(*content); i++ {
		if i < 64 {
			res += string((*content)[i])
			continue
		}

		aux += string((*content)[i])
	}

	*content = aux
	return res
}

func RoundToNext(number float64) int64 {
	r := math.Trunc(number)
	if number-r > 0 {
		return int64(r + 1)
	}

	return int64(r)
}

func GetStartOfFreeBlocks(oldContent string, newContent string, file *os.File, fit rune, sp structs.SuperBlock, w http.ResponseWriter) int64 {
	finalContent := oldContent + newContent
	blocksNeeded := RoundToNext(float64(len(finalContent)/64)) - RoundToNext(float64(len(oldContent)/64))
	if ToInt(sp.S_free_blocks_count[:]) < int64(blocksNeeded) {
		WriteResponse(w, "$Error: no space available")
		return -1
	}

	availableSpace := make([]int64, 0)
	firstBlock := make([]int64, 0)

	var b byte
	state := 0
	counter := 0
	for i := ToInt(sp.S_bm_block_start[:]); i < ToInt(sp.S_blocks_count[:])+ToInt(sp.S_bm_block_start[:]); i++ {
		file.Seek(i, os.SEEK_SET)
		ReadByte(&b, file)
		switch state {
		case 0:
			if b == '0' {
				counter++
				state = 1
				firstBlock = append(firstBlock, i)
			}
			break
		case 1:
			if b == '0' {
				counter++

				if i == ToInt(sp.S_blocks_count[:])+ToInt(sp.S_bm_block_start[:])-1 {
					availableSpace = append(availableSpace, int64(counter))
				}
			} else {
				state = 0
				availableSpace = append(availableSpace, int64(counter))
				counter = 0
			}
			break
		}
	}

	index := -1
	actualAvailable := int64(0)
	if fit == 'w' {
		for i := 0; i < len(availableSpace); i++ {
			if availableSpace[i] >= int64(blocksNeeded) && availableSpace[i] > int64(actualAvailable) {
				actualAvailable = availableSpace[i]
				index = i
			}
		}
	} else if fit == 'b' {
		for i := 0; i < len(availableSpace); i++ {
			if availableSpace[i] >= int64(blocksNeeded) && availableSpace[i] < int64(actualAvailable) {
				actualAvailable = availableSpace[i]
				index = i
			}
		}
	} else {
		for i := 0; i < len(availableSpace); i++ {
			if availableSpace[i] >= int64(blocksNeeded) {
				index = i
				break
			}
		}
	}

	if index == -1 {
		WriteResponse(w, "$Error: no space available")
		return -1
	}

	return firstBlock[index]
}

func GetStartOfFreeBlocks2(oldContent string, newContent string, file *os.File, fit rune, sp structs.SuperBlock, w http.ResponseWriter) int64 {
	blocksNeeded := RoundToNext(float64(len(newContent)/64)) - RoundToNext(float64(len(oldContent)/64))
	if ToInt(sp.S_free_blocks_count[:]) < int64(blocksNeeded) {
		WriteResponse(w, "$Error: no space available")
		return -1
	}

	availableSpace := make([]int64, 0)
	firstBlock := make([]int64, 0)

	var b byte
	state := 0
	counter := 0
	for i := ToInt(sp.S_bm_block_start[:]); i < ToInt(sp.S_blocks_count[:])+ToInt(sp.S_bm_block_start[:]); i++ {
		file.Seek(i, os.SEEK_SET)
		ReadByte(&b, file)
		switch state {
		case 0:
			if b == '0' {
				counter++
				state = 1
				firstBlock = append(firstBlock, i)
			}
			break
		case 1:
			if b == '0' {
				counter++

				if i == ToInt(sp.S_blocks_count[:])+ToInt(sp.S_bm_block_start[:])-1 {
					availableSpace = append(availableSpace, int64(counter))
				}
			} else {
				state = 0
				availableSpace = append(availableSpace, int64(counter))
				counter = 0
			}
			break
		}
	}

	index := -1
	actualAvailable := int64(0)
	if fit == 'w' {
		for i := 0; i < len(availableSpace); i++ {
			if availableSpace[i] >= int64(blocksNeeded) && availableSpace[i] > int64(actualAvailable) {
				actualAvailable = availableSpace[i]
				index = i
			}
		}
	} else if fit == 'b' {
		for i := 0; i < len(availableSpace); i++ {
			if availableSpace[i] >= int64(blocksNeeded) && availableSpace[i] < int64(actualAvailable) {
				actualAvailable = availableSpace[i]
				index = i
			}
		}
	} else {
		for i := 0; i < len(availableSpace); i++ {
			if availableSpace[i] >= int64(blocksNeeded) {
				index = i
				break
			}
		}
	}

	if index == -1 {
		WriteResponse(w, "$Error: no space available")
		return -1
	}

	return firstBlock[index]
}

func GetFreeBlock(sp structs.SuperBlock, file *os.File, block int64, createdBlocks *int64) int64 {
	file.Seek(block, os.SEEK_SET)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, '1')
	writeBinary(file, buffer.Bytes())

	(*createdBlocks)++

	return block - ToInt(sp.S_bm_block_start[:])
}

func GetFreeInode(sp structs.SuperBlock, file *os.File) int64 {
	file.Seek(ToInt(sp.S_bm_inode_start[:]), os.SEEK_SET)
	var c byte
	counter := int64(0)

	for i := ToInt(sp.S_bm_inode_start[:]); i < ToInt(sp.S_inodes_count[:])+ToInt(sp.S_bm_inode_start[:]); i++ {
		ReadByte(&c, file)
		if c == '0' {
			c = '1'
			file.Seek(i, os.SEEK_SET)
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, &c)
			writeBinary(file, buffer.Bytes())
			return counter
		}
		counter++
	}

	return -1
}

func WriteInContentBlock(content string, pointer int64, file *os.File, sp structs.SuperBlock, freeBlocks *int64, createdBlocks *int64, w http.ResponseWriter) int64 {
	p := pointer
	newBlock := false

	if pointer == -1 {
		p = GetFreeBlock(sp, file, *freeBlocks, createdBlocks)
		(*freeBlocks)++
		newBlock = true
	}

	var fb structs.FileBlock
	copy(fb.B_content[:], []byte(content))
	file.Seek(ToInt(sp.S_block_start[:])+(p*64), os.SEEK_SET)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &fb)
	writeBinary(file, buffer.Bytes())

	if newBlock {
		return p
	}

	return -2
}

func WriteInFile(sp structs.SuperBlock, f *structs.Inode, content string, file *os.File, pointerInode int64, freeBlocks *int64, createdBlocks *int64, w http.ResponseWriter) bool {
	for i := 0; i < 16; i++ {
		if content == "" {
			f.I_block[i] = -1
			copy(f.I_mtime[:], []byte(getDate()))

			file.Seek(ToInt(sp.S_inode_start[:])+(pointerInode*int64(binary.Size(*f))), os.SEEK_SET)
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, &f)
			writeBinary(file, buffer.Bytes())
			continue
		}

		response := WriteInContentBlock(SeparateContent(&content), f.I_block[i], file, sp, freeBlocks, createdBlocks, w)

		if response == -1 {
			return false
		} else if response == -2 {
			continue
		}

		f.I_block[i] = response
		copy(f.I_mtime[:], []byte(getDate()))

		file.Seek(ToInt(sp.S_inode_start[:])+(pointerInode*int64(binary.Size(*f))), os.SEEK_SET)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, &f)
		writeBinary(file, buffer.Bytes())
	}
	return true
}

func Permission(read *bool, write *bool, exec *bool, num int64) {
	switch num {
	case 0:
		*read = false
		*write = false
		*exec = false
		break
	case 1:
		*read = false
		*write = false
		*exec = true
		break
	case 2:
		*read = false
		*write = true
		*exec = false
		break
	case 3:
		*read = false
		*write = true
		*exec = true
		break
	case 4:
		*read = true
		*write = false
		*exec = false
		break
	case 5:
		*read = true
		*write = false
		*exec = true
		break
	case 6:
		*read = true
		*write = true
		*exec = false
		break
	case 7:
		*read = true
		*write = true
		*exec = true
		break
	}
}

func GetPermission(inode structs.Inode, userId int64, groupId int64, ugo int64, read bool, write bool, execute bool) bool {
	u := math.Floor(float64(ugo / 100))
	g := math.Floor(float64((ugo % 100) / 10))
	o := ugo % 10
	r := false
	w := false
	x := false

	if ToInt(inode.I_uid[:]) == userId {
		Permission(&r, &w, &x, int64(u))
	} else if ToInt(inode.I_gid[:]) == groupId {
		Permission(&r, &w, &x, int64(g))
	} else {
		Permission(&r, &w, &x, o)
	}

	if read {
		return r
	}

	if write {
		return w
	}

	return x
}

func GetFreeBlock2(sp structs.SuperBlock, file *os.File) int64 {
	file.Seek(ToInt(sp.S_bm_block_start[:]), os.SEEK_SET)
	var c byte
	counter := int64(0)
	for i := ToInt(sp.S_bm_block_start[:]); i < ToInt(sp.S_inode_start[:]); i++ {
		ReadByte(&c, file)
		if c == '0' {
			c = '1'
			file.Seek(i, os.SEEK_SET)
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, &c)
			writeBinary(file, buffer.Bytes())
			return counter
		}
		counter++
	}
	return -1
}

func NewDirBlock(n *structs.Dirblock) {
	for i := 0; i < 4; i++ {
		n.B_content[i].B_inodo = -1
		copy(n.B_content[i].B_name[:], "")
	}
}

func WriteInode(tyype int64, user structs.User, file *os.File, istart int64, bstart int64, currentUser *structs.Sesion, father int64, createdBlocks *int64, createdInodes *int64, perm int64, w http.ResponseWriter) int64 {
	mountedPartition := &(currentUser.Mounted)
	var sp structs.SuperBlock

	start := int64(0)
	if mountedPartition.IsLogic {
		start = ToInt(mountedPartition.LogicPar.Part_start[:])
	} else {
		start = ToInt(mountedPartition.Par.Part_start[:])
	}

	file.Seek(start, os.SEEK_SET)
	ReadSuperBlock(&sp, file)

	var root structs.Inode
	file.Seek(istart, os.SEEK_SET)
	ReadInode(&root, file)
	pointerOfFile := int64(0)

	root = SearchFile(file, root, SplithPath("users.txt"), istart, bstart, &pointerOfFile)
	if root.I_type == 'n' {
		WriteResponse(w, "$Error: users.txt does not exist")
		return -1
	}

	content := ReadFile(file, root, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), w)

	if content == "" {
		return -1
	}

	groupId := GetGroupId(ToString(user.Group[:]), content)

	newInode := structs.Inode{}
	newInodePointer := GetFreeInode(sp, file)
	if newInodePointer == -1 {
		return -1
	}

	for i := 0; i < 16; i++ {
		newInode.I_block[i] = -1
	}

	copy(newInode.I_ctime[:], []byte(getDate()))
	binary.BigEndian.PutUint64(newInode.I_uid[:], uint64(user.Id))
	binary.BigEndian.PutUint64(newInode.I_gid[:], uint64(groupId))
	binary.BigEndian.PutUint64(newInode.I_perm[:], uint64(perm))
	newInode.I_type = '1'
	binary.BigEndian.PutUint64(newInode.I_size[:], uint64(0))
	if tyype == 0 {
		newInode.I_type = '0'
		newBlock := GetFreeBlock2(sp, file)
		if newBlock == -1 {
			return -1
		}

		(*createdBlocks)++
		db := structs.Dirblock{}
		NewDirBlock(&db)
		db.B_content[0].B_inodo = int32(newInodePointer)
		copy(db.B_content[0].B_name[:], ".")
		db.B_content[1].B_inodo = int32(father)
		copy(db.B_content[1].B_name[:], "..")

		file.Seek(bstart+(64*newBlock), os.SEEK_SET)
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.BigEndian, &db)
		writeBinary(file, buffer.Bytes())
		newInode.I_block[0] = newBlock
	}

	file.Seek(istart+int64(binary.Size(newInode)*int(newInodePointer)), os.SEEK_SET)
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, &newInode)
	(*createdInodes)++
	writeBinary(file, buffer.Bytes())

	return newInodePointer
}

func CreateDirectory(file *os.File, father *structs.Inode, istart int64, bstart int64, user structs.User, currentUser *structs.Sesion, name string, tyype int64, sp structs.SuperBlock, createdBlocks *int64, createdInodes *int64, p *int64, perm int64, w http.ResponseWriter) int64 {
	file.Seek(bstart+(64*father.I_block[0]), os.SEEK_SET)
	dbf := structs.Dirblock{}
	ReadDirBlock(&dbf, file)
	fatherPointer := dbf.B_content[0].B_inodo
	for i := 0; i < 16; i++ {
		if father.I_block[i] == -1 {
			pointerOfNewDir := GetFreeBlock2(sp, file)
			if pointerOfNewDir == -1 {
				return -1
			}
			newDir := structs.Dirblock{}
			NewDirBlock(&newDir)

			newInode := WriteInode(tyype, user, file, istart, bstart, currentUser, int64(fatherPointer), createdBlocks, createdInodes, perm, w)
			if newInode == -1 {
				return -1
			}

			newDir.B_content[0].B_inodo = int32(newInode)
			copy(newDir.B_content[0].B_name[:], []byte(name))

			father.I_block[i] = pointerOfNewDir

			file.Seek(bstart+(64*pointerOfNewDir), os.SEEK_SET)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &newDir)
			writeBinary(file, buffer.Bytes())

			file.Seek(istart+(int64(binary.Size(father))*int64(fatherPointer)), os.SEEK_SET)
			buffer = bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, father)
			writeBinary(file, buffer.Bytes())
			return -2
		} else {
			dirBlock := structs.Dirblock{}
			file.Seek(bstart+(father.I_block[i]*64), os.SEEK_SET)
			ReadDirBlock(&dirBlock, file)

			for j := 0; j < 4; j++ {
				if ToString(dirBlock.B_content[j].B_name[:]) == name {
					WriteResponse(w, "$Error: "+name+" already exist")
					return -1
				}
				if dirBlock.B_content[j].B_inodo == -1 {
					newInode := WriteInode(tyype, user, file, istart, bstart, currentUser, int64(fatherPointer), createdBlocks, createdInodes, perm, w)
					if newInode == -1 {
						return -1
					}
					dirBlock.B_content[j].B_inodo = int32(newInode)
					copy(dirBlock.B_content[j].B_name[:], []byte(name))

					file.Seek(bstart+(64*father.I_block[i]), os.SEEK_SET)
					buffer := bytes.Buffer{}
					binary.Write(&buffer, binary.BigEndian, &dirBlock)
					writeBinary(file, buffer.Bytes())
					return -2
				}
			}
		}
	}

	return -1
}

func CreateMultipleDirectories(path []string, father structs.Inode, file *os.File, root structs.Inode, istart int64, bstart int64, userId int64, groupId int64, currentUser structs.Sesion, sp structs.SuperBlock, createdBlocks *int64, createdInodes *int64, p *int64, perm int64, w http.ResponseWriter) bool {
	auxPath := make([]string, 0)
	response := root
	for i := 0; i < len(path); i++ {
		auxPath = append(auxPath, path[i])
		response = SearchFile(file, root, auxPath, istart, bstart, p)
		if response.I_type == 'n' {
			if GetPermission(father, userId, groupId, ToInt(father.I_perm[:]), false, true, false) {
				CreateDirectory(file, &father, istart, bstart, currentUser.Usr, &currentUser, path[i], 0, sp, createdBlocks, createdInodes, p, perm, w)
				file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
				ReadInode(&root, file)
				father = SearchFile(file, root, auxPath, istart, bstart, p)
			} else {
				WriteResponse(w, "$Error: you dont have permission to write")
				return false
			}
		} else {
			father = response
		}
	}

	return true
}
