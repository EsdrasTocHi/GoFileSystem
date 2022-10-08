package analyzer

import (
	fs "Backend/FileSystem"
	structs "Backend/Structures"
	"net/http"
	"strings"
)

var number int = 0
var mountedPartitions []structs.MountedPartition = make([]structs.MountedPartition, 0)

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
	} else {
		fs.WriteResponse(w, "$Error: "+cmd+" command not recognized")
	}
}
