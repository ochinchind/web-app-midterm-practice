package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type hunt struct {
	ID         int64     `json:"id"`
	CriminalID int64     `json:"criminal_id"`
	Success    bool      `json:"success"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

const huntColumns = `id, criminal_id, success, notes, created_at`

func scanHunt(row pgx.Row) (hunt, error) {
	var h hunt
	err := row.Scan(&h.ID, &h.CriminalID, &h.Success, &h.Notes, &h.CreatedAt)
	return h, err
}

// GET /api/hunts/?criminal_id=
func (s *server) listHunts(w http.ResponseWriter, r *http.Request) {
	sql := "SELECT " + huntColumns + " FROM hunt_logs"
	var args []any
	if raw := r.URL.Query().Get("criminal_id"); raw != "" {
		cid, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || cid < 1 {
			writeError(w, http.StatusBadRequest, "criminal_id must be a positive integer")
			return
		}
		sql += " WHERE criminal_id = $1"
		args = append(args, cid)
	}
	sql += " ORDER BY id ASC"

	rows, err := s.db.Query(r.Context(), sql, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer rows.Close()

	out := []hunt{}
	for rows.Next() {
		h, err := scanHunt(rows)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}
		out = append(out, h)
	}
	if rows.Err() != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

// GET /api/hunts/{id}/
func (s *server) getHunt(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	h, err := scanHunt(s.db.QueryRow(r.Context(),
		"SELECT "+huntColumns+" FROM hunt_logs WHERE id = $1", id))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "hunt log not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, h)
}

// POST /api/hunts/ 
func (s *server) createHunt(w http.ResponseWriter, r *http.Request) {
	var in struct {
		CriminalID *int64  `json:"criminal_id"`
		Success    *bool   `json:"success"`
		Notes      *string `json:"notes"`
	}
	if err := decodeBody(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	switch {
	case in.CriminalID == nil || *in.CriminalID < 1:
		writeError(w, http.StatusBadRequest, "criminal_id is required and must be a positive integer")
		return
	case in.Success == nil:
		writeError(w, http.StatusBadRequest, "success is required and must be a boolean")
		return
	}
	notes := ""
	if in.Notes != nil {
		notes = *in.Notes
	}

	h, err := scanHunt(s.db.QueryRow(r.Context(),
		`INSERT INTO hunt_logs (criminal_id, success, notes) VALUES ($1, $2, $3) RETURNING `+huntColumns,
		*in.CriminalID, *in.Success, notes))
	if isForeignKeyViolation(err) {
		writeError(w, http.StatusBadRequest, "criminal_id does not reference an existing criminal")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	w.Header().Set("Location", "/api/hunts/"+strconv.FormatInt(h.ID, 10)+"/")
	writeJSON(w, http.StatusCreated, h)
}

// PATCH /api/hunts/{id}/ 
func (s *server) updateHunt(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in struct {
		Success *bool   `json:"success"`
		Notes   *string `json:"notes"`
	}
	if err := decodeBody(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	if in.Success == nil && in.Notes == nil {
		writeError(w, http.StatusBadRequest, "provide at least one of: success (boolean), notes (string)")
		return
	}

	h, err := scanHunt(s.db.QueryRow(r.Context(),
		`UPDATE hunt_logs
		    SET success = COALESCE($2, success),
		        notes   = COALESCE($3, notes)
		  WHERE id = $1
		  RETURNING `+huntColumns, id, in.Success, in.Notes))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "hunt log not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, h)
}

// DELETE /api/hunts/{id}/
func (s *server) deleteHunt(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	tag, err := s.db.Exec(r.Context(), "DELETE FROM hunt_logs WHERE id = $1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	if tag.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "hunt log not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeBody(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)
	return json.NewDecoder(r.Body).Decode(dst)
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
