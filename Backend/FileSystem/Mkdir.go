package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
)

func MkDir(path string, r bool, currentUser structs.Sesion, activeSession bool, perm int64, w http.ResponseWriter) {
	if !activeSession {
		WriteResponse(w, "$Error: there is no active session")
		return
	}

	start := int64(0)
	if currentUser.Mounted.IsLogic {
		start = ToInt(currentUser.Mounted.LogicPar.Part_start[:])
	} else {
		start = ToInt(currentUser.Mounted.Par.Part_start[:])
	}

	sb := structs.SuperBlock{}
	file, _ := os.OpenFile(currentUser.Mounted.Path, os.O_RDWR, 0777)
	defer file.Close()
	file.Seek(start, os.SEEK_SET)
	ReadSuperBlock(&sb, file)

	root := structs.Inode{}
	aux := structs.Inode{}
	file.Seek(ToInt(sb.S_inode_start[:]), os.SEEK_SET)
	ReadInode(&root, file)
	pointerOfFile := int64(0)
	aux = SearchFile(file, root, SplithPath("users.txt"), ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), &pointerOfFile)
	c := ReadFile(file, aux, ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), w)
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
		aux = SearchFile(file, root, p, ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), &pointerOfFile)
	}

	createdBlocks := int64(0)
	createdInodes := int64(0)
	if r {
		if aux.I_type == 'n' {
			if !CreateMultipleDirectories(p, root, file, root, ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), currentUser.Usr.Id, int64(GetGroupId(ToString(currentUser.Usr.Group[:]), c)), currentUser, sb, &createdBlocks, &createdInodes, &pointerOfFile, perm, w) {
				return
			}
			file.Seek(ToInt(sb.S_inode_start[:]), os.SEEK_SET)
			ReadInode(&root, file)
			aux = SearchFile(file, root, p, ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), &pointerOfFile)
		}
	}

	if aux.I_type == 'n' {
		WriteResponse(w, "$Error: The directory does not exist")
	} else {
		if GetPermission(aux, currentUser.Usr.Id, int64(GetGroupId(ToString(currentUser.Usr.Group[:]), c)), ToInt(aux.I_perm[:]), false, true, false) {
			if CreateDirectory(file, &aux, ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), currentUser.Usr, &currentUser, paux[len(paux)-1], 0, sb, &createdBlocks, &createdInodes, &pointerOfFile, int64(perm), w) == -1 {
				return
			}
			file.Seek(ToInt(sb.S_inode_start[:]), os.SEEK_SET)
			ReadInode(&root, file)

			binary.BigEndian.PutUint64(sb.S_free_blocks_count[:], uint64(ToInt(sb.S_free_blocks_count[:])-createdBlocks))
			binary.BigEndian.PutUint64(sb.S_free_inodes_count[:], uint64(ToInt(sb.S_free_inodes_count[:])-createdInodes))

			file.Seek(start, os.SEEK_SET)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &sb)
			writeBinary(file, buffer.Bytes())
			WriteResponse(w, "DIRECTORY CREATED SUCCESFULLY")
		} else {
			WriteResponse(w, "$Error: you dont have permission to write")
		}
	}
}
