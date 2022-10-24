package analyzer

import (
	fs "Backend/FileSystem"
	structs "Backend/Structures"
	"net/http"
	"strings"
)

var number int = 0
var mountedPartitions []structs.MountedPartition = make([]structs.MountedPartition, 0)
var activeSession bool = false
var currentUser structs.Sesion

func MkDisk(params []string, w http.ResponseWriter) {
	var path string = ""
	var unit rune = 'm'
	var fit rune = 'f'
	var size int = 0

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-size" {
			size = Size(value, w)
			if size == -1 {
				return
			}
		} else if name == "-path" {
			path = Path(value, w)
			if path == "" {
				return
			}
		} else if name == "-unit" {
			unit = UnitMkdisk(value, w)
			if unit == 'e' {
				return
			}
		} else if name == "-fit" {
			fit = Fit(value, w)
			if fit == 'e' {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if path == "" {
		fs.WriteResponse(w, "$Error: PATH is a mandatory parameter")
	}
	if size == 0 {
		fs.WriteResponse(w, "$Error: SIZE is a mandatory parameter")
	}

	fs.MkDisk(size, unit, fit, path, w)
}

func RmDisk(params []string, w http.ResponseWriter) {
	var path string = ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-path" {
			path = Path(value, w)
			if path == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if path == "" {
		fs.WriteResponse(w, "$Error: PATH is a mandatory parameter")
	}

	fs.Rmdisk(path, w)
}

func FDisk(params []string, w http.ResponseWriter) {
	var path string = ""
	var nameOfPar string = ""
	var unit rune = 'k'
	var fit rune = 'w'
	var tyype rune = 'p'
	var size int = 0

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-size" {
			size = Size(value, w)
			if size == -1 {
				return
			}
		} else if name == "-path" {
			path = Path(value, w)
			if path == "" {
				return
			}
		} else if name == "-unit" {
			unit = UnitFdisk(value, w)
			if unit == 'e' {
				return
			}
		} else if name == "-fit" {
			fit = Fit(value, w)
			if fit == 'e' {
				return
			}
		} else if name == "-name" {
			nameOfPar = Name(value, w)
			if nameOfPar == "" {
				return
			}
		} else if name == "-type" {
			tyype = Type(value, w)
			if tyype == 'i' {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if path == "" {
		fs.WriteResponse(w, "$Error: PATH is a mandatory parameter")
	}
	if nameOfPar == "" {
		fs.WriteResponse(w, "$Error: NAME is a mandatory parameter")
	}
	if size == 0 {
		fs.WriteResponse(w, "$Error: SIZE is a mandatory parameter")
	}

	fs.FDisk(size, unit, path, tyype, fit, nameOfPar, w)
}

func Mount(params []string, w http.ResponseWriter) {
	var path string = ""
	var nameOfPar string = ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-path" {
			path = Path(value, w)
			if path == "" {
				return
			}
		} else if name == "-name" {
			nameOfPar = Name(value, w)
			if nameOfPar == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if path == "" {
		fs.WriteResponse(w, "$Error: PATH is a mandatory parameter")
	}
	if nameOfPar == "" {
		fs.WriteResponse(w, "$Error: NAME is a mandatory parameter")
	}

	fs.Mount(path, nameOfPar, &mountedPartitions, &number, w)
}

func Mkfs(params []string, w http.ResponseWriter) {
	var id string = ""
	var tyype bool = false

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-id" {
			id = Id(value, w)
			if id == "" {
				return
			}
		} else if name == "-type" {
			tyype = TypeMkfs(value, w)
			if tyype == false {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if id == "" {
		fs.WriteResponse(w, "$Error: ID is a mandatory parameter")
	}

	fs.Mkfs(id, &mountedPartitions, w)
}

func Login(params []string, w http.ResponseWriter) {
	var id string = ""
	usuario := ""
	password := ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-id" {
			id = Id(value, w)
			if id == "" {
				return
			}
		} else if name == "-usuario" {
			usuario = Usuario(value, w)
			if usuario == "" {
				return
			}
		} else if name == "-password" {
			password = Password(value, w)
			if password == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if id == "" {
		fs.WriteResponse(w, "$Error: ID is a mandatory parameter")
	}
	if usuario == "" {
		fs.WriteResponse(w, "$Error: USUARIO is a mandatory parameter")
	}
	if password == "" {
		fs.WriteResponse(w, "$Error: PASSWORD is a mandatory parameter")
	}

	fs.Login(usuario, password, id, &mountedPartitions, &currentUser, &activeSession, w)
}

func Mkgrp(params []string, w http.ResponseWriter) {
	var id string = ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-name" {
			id = Name(value, w)
			if id == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if id == "" {
		fs.WriteResponse(w, "$Error: NAME is a mandatory parameter")
	}

	fs.Mkgrp(id, &currentUser, &activeSession, w)
}

func Mkusr(params []string, w http.ResponseWriter) {
	var grp string = ""
	usuario := ""
	password := ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-grp" {
			grp = Grp(value, w)
			if grp == "" {
				return
			}
		} else if name == "-usuario" {
			usuario = Usuario(value, w)
			if usuario == "" {
				return
			}
		} else if name == "-pwd" {
			password = Password(value, w)
			if password == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if grp == "" {
		fs.WriteResponse(w, "$Error: GRP is a mandatory parameter")
	}
	if usuario == "" {
		fs.WriteResponse(w, "$Error: USUARIO is a mandatory parameter")
	}
	if password == "" {
		fs.WriteResponse(w, "$Error: PASSWORD is a mandatory parameter")
	}

	fs.Mkusr(usuario, password, grp, &currentUser, &activeSession, w)
}

func Rmgrp(params []string, w http.ResponseWriter) {
	var id string = ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-name" {
			id = Name(value, w)
			if id == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if id == "" {
		fs.WriteResponse(w, "$Error: NAME is a mandatory parameter")
	}

	fs.Rmgrp(id, &currentUser, &activeSession, w)
}

func Rmusr(params []string, w http.ResponseWriter) {
	var id string = ""

	for i := 0; i < len(params); i++ {
		param := strings.Split(params[i], "=")
		name := strings.ToLower(param[0])
		value := param[1]
		if name == "-usuario" {
			id = Usuario(value, w)
			if id == "" {
				return
			}
		} else {
			fs.WriteResponse(w, "$Error: "+strings.Trim(name, "-")+" is not a valid parameter")
			return
		}
	}

	if id == "" {
		fs.WriteResponse(w, "$Error: USUARIO is a mandatory parameter")
	}

	fs.Rmusr(id, &currentUser, &activeSession, w)
}

func ReadCommand(cmd string, w http.ResponseWriter) {
	cmd = strings.Trim(cmd, " ")

	command := SplitCommand(cmd, ' ')
	cmd = strings.ToLower(command[0])
	params := command[1:]

	if cmd == "mkdisk" {
		MkDisk(params, w)
	} else if cmd == "rmdisk" {
		RmDisk(params, w)
	} else if cmd == "fdisk" {
		FDisk(params, w)
	} else if cmd == "mount" {
		Mount(params, w)
	} else if cmd == "mkfs" {
		Mkfs(params, w)
	} else if cmd == "login" {
		Login(params, w)
	} else if cmd == "logout" {
		fs.Logout(&currentUser, &activeSession, w)
	} else if cmd == "mkgrp" {
		Mkgrp(params, w)
	} else if cmd == "mkusr" {
		Mkusr(params, w)
	} else if cmd == "rmusr" {
		Rmusr(params, w)
	} else if cmd == "rmgrp" {
		Rmgrp(params, w)
	} else {
		fs.WriteResponse(w, "$Error: "+cmd+" command not recognized")
	}
}
