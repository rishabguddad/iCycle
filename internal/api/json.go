package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJSON sends data to the client as JSON, with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError sends an error message as JSON, e.g. {"error":"product not found"}.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// decodeJSON reads the request body into dst, which must be a pointer.
// Bodies over 1 MB and unknown fields (often typos) are rejected.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// pathID reads {id} from the URL. If it isn't a number, pathID sends a 400
// reply itself and returns false, and the handler should simply return.
func pathID(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, name+" id must be a number")
		return 0, false
	}
	return id, true
}
