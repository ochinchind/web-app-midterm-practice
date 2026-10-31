package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type server struct {
	db *pgxpool.Pool
}

func NewRouter(db *pgxpool.Pool) http.Handler {
	s := &server{db: db}
	mux := http.NewServeMux()

	handle := func(method, path string, h http.HandlerFunc) {
		mux.HandleFunc(method+" "+path+"{$}", h)
		mux.HandleFunc(method+" "+path[:len(path)-1], h)
	}

	handle("GET", "/api/criminals/", s.listCriminals)
	handle("GET", "/api/criminals/{id}/", s.getCriminal)

	handle("GET", "/api/hunts/", s.listHunts)
	handle("POST", "/api/hunts/", s.createHunt)
	handle("GET", "/api/hunts/{id}/", s.getHunt)
	handle("PATCH", "/api/hunts/{id}/", s.updateHunt)
	handle("DELETE", "/api/hunts/{id}/", s.deleteHunt)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "route not found")
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// pathID parses the {id} path value; on failure it writes a 400 and returns false.
func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}
