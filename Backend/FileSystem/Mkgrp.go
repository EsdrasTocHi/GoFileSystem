package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
	"strconv"
)

func Mkgrp(name string, currentUser *structs.Sesion, activeSession *bool, w http.ResponseWriter) {
	if len(name) > 10 {
		WriteResponse(w, "$Error: the name is too long")
		return
	}

	if !(*activeSession) {
		WriteResponse(w, "$Error: there is no active session")
		return
	}

	if ToString(currentUser.Usr.Name[:]) != "root" {
		WriteResponse(w, "$Error: you do not have permission to use this command")
		return
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

	root = SearchFile(file, root, SplithPath("users.txt"), ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)

	if root.I_type == 'n' {
		WriteResponse(w, "$Error: users.txt does not exist")
		return
	}

	content := ReadFile(file, root, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), w)

	if content == "" {
		return
	}

	if ExistGroup(name, content) {
		WriteResponse(w, "$Error: the group already exist")
		return
	}

	id := CountGroups(content) + 1
	newGroup := strconv.Itoa(id) + ",G," + name + "\n"

	freeBlocks := GetStartOfFreeBlocks(content, newGroup, file, rune(fit), sp, w)
	createdBlocks := int64(0)
	if freeBlocks == -1 {
		return
	}

	content += newGroup
	newSize := len(content)
	if WriteInFile(sp, &root, content, file, int64(pointerOfFile), &freeBlocks, &createdBlocks, w) {
		binary.BigEndian.PutUint64(sp.S_free_blocks_count[:], uint64(ToInt(sp.S_free_blocks_count[:])-createdBlocks))
		file.Seek(start, os.SEEK_SET)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, &sp)
		writeBinary(file, buffer.Bytes())

		binary.BigEndian.PutUint64(root.I_size[:], uint64(newSize))
		copy(root.I_mtime[:], []byte(getDate()))
		file.Seek(ToInt(sp.S_inode_start[:])+int64(binary.Size(root))*int64(pointerOfFile), os.SEEK_SET)
		buffer = bytes.Buffer{}
		binary.Write(&buffer, binary.BigEndian, &root)
		writeBinary(file, buffer.Bytes())
		WriteResponse(w, "GROUP CREATED SUCCESFULLY")
	}
}
