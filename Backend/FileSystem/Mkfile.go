package filesystem

import (
	structs "Backend/Structures"
	"bufio"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
	"strconv"
)

func Mkfile(path string, r bool, size int64, contPath string, currentUser structs.Sesion, activeSession bool, perm int, w http.ResponseWriter) {
	if !activeSession {
		WriteResponse(w, "$Error: there is no active session")
		return
	}

	content := ""
	if contPath != "" {
		if !Exist(contPath) {
			WriteResponse(w, "$Error: the file for cont does not exist")
			return
		}

		readFile, _ := os.Open(contPath)
		defer readFile.Close()

		fileScanner := bufio.NewScanner(readFile)

		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			content += fileScanner.Text() + "\n"
		}
	} else if size != 0 {
		for i := int64(0); i < size; i++ {
			content += strconv.Itoa(int(i % 10))
		}
	}

	var mountedPartition *structs.MountedPartition
	mountedPartition = &(currentUser.Mounted)
	file, _ := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0777)
	defer file.Close()

	var sp structs.SuperBlock
	start := int64(0)
	var fit byte
	if mountedPartition.IsLogic {
		start = ToInt(mountedPartition.LogicPar.Part_start[:])
		fit = mountedPartition.LogicPar.Part_fit
	} else {
		start = ToInt(mountedPartition.Par.Part_start[:])
		fit = mountedPartition.Par.Part_fit
	}

	file.Seek(start, os.SEEK_SET)
	ReadSuperBlock(&sp, file)

	var root structs.Inode
	file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
	ReadInode(&root, file)

	pointerOfFile := int64(0)

	aux := SearchFile(file, root, SplithPath("users.txt"), ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)
	c := ReadFile(file, aux, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), w)
	p := make([]string, 0)
	paux := SplithPath(path)

	for i := 0; i < len(paux); i++ {
		if i < len(paux)-1 {
			p = append(p, paux[i])
		}
	}

	if len(p) == 0 {
		aux = root
	} else {
		aux = SearchFile(file, root, p, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)
	}

	createdBlocks := int64(0)
	createdInodes := int64(0)

	if r {
		if aux.I_type == 'n' {
			if !CreateMultipleDirectories(p, root, file, root, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), currentUser.Usr.Id, int64(GetGroupId(ToString(currentUser.Usr.Group[:]), c)), currentUser, sp, &createdBlocks, &createdInodes, &pointerOfFile, int64(perm), w) {
				return
			}
			file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
			ReadInode(&root, file)
			aux = SearchFile(file, root, p, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)
		}
	}

	freeBlocks := GetStartOfFreeBlocks("", content, file, rune(fit), sp, w)
	if aux.I_type == 'n' {
		WriteResponse(w, "$Error: The directory does not exist")
		return
	} else {
		if GetPermission(aux, currentUser.Usr.Id, int64(GetGroupId(ToString(currentUser.Usr.Group[:]), c)), ToInt(aux.I_perm[:]), false, true, false) {
			file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
			ReadInode(&root, file)
			if CreateDirectory(file, &aux, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), currentUser.Usr, &currentUser, paux[len(paux)-1], 1, sp, &createdBlocks, &createdInodes, &pointerOfFile, int64(perm), w) == -1 {
				return
			}
			file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
			ReadInode(&root, file)
			aux = SearchFile(file, root, SplithPath(path), ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)

			if WriteInFile(sp, &aux, content, file, pointerOfFile, &freeBlocks, &createdBlocks, w) {
				binary.BigEndian.PutUint64(aux.I_size[:], uint64(len(content)))
				file.Seek(ToInt(sp.S_inode_start[:])+(pointerOfFile*int64(binary.Size(root))), os.SEEK_SET)
				buffer := bytes.Buffer{}
				binary.Write(&buffer, binary.BigEndian, &aux)
				writeBinary(file, buffer.Bytes())
			}
			file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
			ReadInode(&root, file)

			binary.BigEndian.PutUint64(sp.S_free_blocks_count[:], uint64(ToInt(sp.S_free_blocks_count[:])-createdBlocks))
			binary.BigEndian.PutUint64(sp.S_free_inodes_count[:], uint64(ToInt(sp.S_free_inodes_count[:])-createdInodes))

			file.Seek(start, os.SEEK_SET)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &sp)
			writeBinary(file, buffer.Bytes())
			WriteResponse(w, "FILE CREATED SUCCESFULLY")
		} else {
			WriteResponse(w, "$Error: you dont have permission to write")
		}
	}
}
