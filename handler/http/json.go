package handler

import (
	"net/http"
)

// RespondWithError ..
/*func RespondWithMessage(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"message": msg})

}
*/
// RespondWithJSON ..
func RespondWithJSON(w http.ResponseWriter, code int, resp []byte) {

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("encoding", "utf-8")
	w.WriteHeader(code)
	w.Write(resp)
}
