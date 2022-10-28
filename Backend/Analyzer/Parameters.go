package analyzer

import (
	fs "Backend/FileSystem"
	"net/http"
	"strconv"
	"strings"
)

func Size(param string, w http.ResponseWriter) int {
	res, err := strconv.Atoi(param)

	if err != nil {
		fs.WriteResponse(w, "$Error: SIZE must be a number")
		return -1
	}

	return res
}

func Fit(param string, w http.ResponseWriter) rune {
	param = strings.ToLower(param)

	if param == "bf" {
		return 'b'
	} else if param == "wf" {
		return 'w'
	} else if param == "ff" {
		return 'f'
	}

	fs.WriteResponse(w, "$Error: FIT value is not valid")
	return 'e'
}

func UnitMkdisk(param string, w http.ResponseWriter) rune {
	param = strings.ToLower(param)

	if param == "k" {
		return 'k'
	} else if param == "m" {
		return 'm'
	}

	fs.WriteResponse(w, "$Error: UNIT value is not valid")
	return 'e'
}

func UnitFdisk(param string, w http.ResponseWriter) rune {
	param = strings.ToLower(param)

	if param == "k" {
		return 'k'
	} else if param == "m" {
		return 'm'
	} else if param == "b" {
		return 'b'
	}

	fs.WriteResponse(w, "$Error: UNIT value is not valid")
	return 'e'
}

func Type(param string, w http.ResponseWriter) rune {
	param = strings.ToLower(param)

	if param == "p" {
		return 'p'
	} else if param == "e" {
		return 'e'
	} else if param == "l" {
		return 'l'
	}

	fs.WriteResponse(w, "$Error: TYPE value is not valid")
	return 'i'
}

func Path(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: PATH cannot be empty")
	return ""
}

func Name(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: NAME cannot be empty")
	return ""
}

func TypeMkfs(param string, w http.ResponseWriter) bool {
	param = strings.ToLower(param)

	if param == "full" {
		return true
	}

	fs.WriteResponse(w, "$Error: TYPE value is not valid")
	return false
}

func Id(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: ID cannot be empty")
	return ""
}

func Usuario(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: USUARIO cannot be empty")
	return ""
}

func Password(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: PASSWORD cannot be empty")
	return ""
}

func Grp(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: GRP cannot be empty")
	return ""
}

func Cont(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: CONT cannot be empty")
	return ""
}
