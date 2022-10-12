package filesystem

import (
	structs "Backend/Structures"
	"net/http"
)

func Logout(currentUser *structs.Sesion, activeSession *bool, w http.ResponseWriter) {
	if !*activeSession {
		WriteResponse(w, "$Error: there is no active session")
		return
	}

	*activeSession = false
	var s structs.Sesion
	WriteResponse(w, "GOOD BYE "+ToString(currentUser.Usr.Name[:])+"!!")
	currentUser = &s
}
