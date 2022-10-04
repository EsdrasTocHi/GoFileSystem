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

func Path(param string, w http.ResponseWriter) string {
	if param != "" {
		return param
	}

	fs.WriteResponse(w, "$Error: PATH cannot be empty")
	return ""
}
