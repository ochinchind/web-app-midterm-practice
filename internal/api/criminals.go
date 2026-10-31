package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type criminal struct {
	ID            int64    `json:"id"`
	Name          *string  `json:"name"`
	Nickname      string   `json:"nickname"`
	Reward        int64    `json:"reward"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	LastKnownNote *string  `json:"last_known_note"`
	Difficulty    int      `json:"difficulty"`
	Armed         bool     `json:"armed"`
	HasCrew       bool     `json:"has_crew"`
	EscapeMaster  bool     `json:"escape_master"`
}

const criminalColumns = `id, name, nickname, reward, latitude, longitude, last_known_note, difficulty, armed, has_crew, escape_master`

func scanCriminal(row pgx.Row) (criminal, error) {
	var c criminal
	err := row.Scan(&c.ID, &c.Name, &c.Nickname, &c.Reward, &c.Latitude, &c.Longitude,
		&c.LastKnownNote, &c.Difficulty, &c.Armed, &c.HasCrew, &c.EscapeMaster)
	return c, err
}

func (s *server) listCriminals(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var (
		where []string
		args  []any
	)
	add := func(cond string, v any) {
		args = append(args, v)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}

	for _, p := range []struct{ param, cond string }{
		{"min_reward", "reward >= $%d"},
		{"max_reward", "reward <= $%d"},
	} {
		if raw := q.Get(p.param); raw != "" {
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil || v < 0 {
				writeError(w, http.StatusBadRequest, p.param+" must be a non-negative integer")
				return
			}
			add(p.cond, v)
		}
	}
	if q.Get("min_reward") != "" && q.Get("max_reward") != "" {
		lo, _ := strconv.ParseInt(q.Get("min_reward"), 10, 64)
		hi, _ := strconv.ParseInt(q.Get("max_reward"), 10, 64)
		if lo > hi {
			writeError(w, http.StatusBadRequest, "min_reward must not exceed max_reward")
			return
		}
	}

	for _, col := range []string{"armed", "has_crew", "escape_master"} {
		if raw := q.Get(col); raw != "" {
			v, err := strconv.ParseBool(raw)
			if err != nil {
				writeError(w, http.StatusBadRequest, col+" must be true or false")
				return
			}
			add(col+" = $%d", v)
		}
	}

	order := "id ASC"
	switch q.Get("sort") {
	case "":
	case "reward":
		order = "reward ASC, id ASC"
	case "-reward":
		order = "reward DESC, id ASC"
	default:
		writeError(w, http.StatusBadRequest, "sort must be one of: reward, -reward")
		return
	}

	sql := "SELECT " + criminalColumns + " FROM wanted_criminals"
	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}
	sql += " ORDER BY " + order

	rows, err := s.db.Query(r.Context(), sql, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer rows.Close()

	out := []criminal{}
	for rows.Next() {
		c, err := scanCriminal(rows)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}
		out = append(out, c)
	}
	if rows.Err() != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *server) getCriminal(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	c, err := scanCriminal(s.db.QueryRow(r.Context(),
		"SELECT "+criminalColumns+" FROM wanted_criminals WHERE id = $1", id))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "criminal not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, c)
}
