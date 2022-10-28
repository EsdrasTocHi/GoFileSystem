package filesystem

import (
	structs "Backend/Structures"
	"net/http"
	"os"
	"strconv"
)

func Authenticate(usr string, passw string, currentUser *structs.Sesion, content string, w http.ResponseWriter) bool {
	usersContent := GetLines(content)
	for i := 0; i < len(usersContent); i++ {
		if IsUser(usersContent[i]) {
			id := 0
			counter := 0
			group := ""
			user := ""
			pass := ""
			aux := ""

			for j := 0; j < len(usersContent[i]); j++ {
				if usersContent[i][j] == ',' || j == len(usersContent[i])-1 {
					switch counter {
					case 0:
						id, _ = strconv.Atoi(aux)
						break
					case 2:
						group = aux
						break
					case 3:
						user = aux
						break
					case 4:
						aux += string(usersContent[i][j])
						pass = aux
						break
					}
					aux = ""
					counter++
					continue
				}
				aux += string(usersContent[i][j])
			}

			if usr == user {
				if pass == passw {
					if id == 0 {
						WriteResponse(w, "$Error: user does not exist")
						return false
					}

					currentUser.Usr = structs.User{}

					copy(currentUser.Usr.Group[:], []byte(group))
					currentUser.Usr.Id = int64(id)
					copy(currentUser.Usr.Name[:], []byte(user))
					copy(currentUser.Usr.Password[:], []byte(pass))
					return true
				}
				WriteResponse(w, "$Error: incorrect password")
				return false
			}
		}
	}

	WriteResponse(w, "$Error: incorrect user")
	return false
}

func Login(usr string, passw string, id string, partitions *[]structs.MountedPartition, currentUser *structs.Sesion, activeSession *bool, w http.ResponseWriter) {
	if *activeSession {
		WriteResponse(w, "$Error: active session")
		return
	}

	i := 0
	var mountedPartition *structs.MountedPartition

	for i = 0; i < len(*partitions); i++ {
		if id == (*partitions)[i].Id {
			mountedPartition = &((*partitions)[i])
			break
		}
	}

	if i == len(*partitions) {
		WriteResponse(w, "$Error: "+id+" is not mounted")
		return
	}

	file, _ := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0777)
	defer file.Close()

	var sp structs.SuperBlock
	var start int64
	if mountedPartition.IsLogic {
		start = ToInt(mountedPartition.LogicPar.Part_start[:])
	} else {
		start = ToInt(mountedPartition.Par.Part_start[:])
	}

	file.Seek(start, os.SEEK_SET)
	ReadSuperBlock(&sp, file)

	var root structs.Inode
	file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
	ReadInode(&root, file)

	p := int64(0)
	root = SearchFile(file, root, SplithPath("users.txt"), ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &p)

	if root.I_type == 'n' {
		WriteResponse(w, "$Error: users.txt does not exist")
		return
	}

	content := ReadFile(file, root, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), w)
	if content == "" {
		return
	}

	*activeSession = Authenticate(usr, passw, currentUser, content, w)

	if *activeSession {
		currentUser.Mounted = *mountedPartition
		WriteResponse(w, "WELCOME "+usr+"!!")
	}
}
