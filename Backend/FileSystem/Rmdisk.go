package filesystem

import (
	"net/http"
	"os"
)

func Rmdisk(path string, w http.ResponseWriter){
	if Exist(path){
		os.Remove(path)
		WriteResponse(w, path+" REMOVED SUCCESSFULLY")
	}else{
		WriteResponse(w, "$Error: the disk does not exists")
	}
}