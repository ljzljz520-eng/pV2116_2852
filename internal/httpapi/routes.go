package httpapi

import (
	"encoding/json"
	"net/http"
	"stickerchallenge/internal/domain"
	"strings"
)

type registerRequest struct {
	ID         string             `json:"id"`
	Label      string             `json:"label"`
	Owner      string             `json:"owner"`
	Candidates []domain.Candidate `json:"candidates"`
}
type reviewRequest struct {
	Actor string `json:"actor"`
}
type updateRequest struct {
	Actor    string `json:"actor"`
	RecordID string `json:"record_id"`
	Number   int    `json:"number"`
	Version  int    `json:"version"`
}
type noteRequest struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (s *Server) registerRoutes() {
	s.Mux.HandleFunc("/health", s.Health)
	s.Mux.HandleFunc("/batches", s.handleBatches)
	s.Mux.HandleFunc("/batches/", s.handleBatch)
}

func (s *Server) handleBatches(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		query := domain.SearchQuery{Label: r.URL.Query().Get("label")}
		batches, err := s.Service.Search(query)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, batches)
		return
	}
	if r.Method != http.MethodPost {
		writeErrorStatus(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}
	batch, err := s.Service.RegisterBatch(request.ID, request.Label, request.Owner, request.Candidates)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, batch)
}

func (s *Server) handleBatch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		writeErrorStatus(w, http.StatusNotFound, "batch path required")
		return
	}
	id := parts[1]
	if len(parts) == 2 && r.Method == http.MethodGet {
		batch, err := s.Service.Store.GetBatch(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, batch)
		return
	}
	if len(parts) != 3 {
		writeErrorStatus(w, http.StatusNotFound, "unknown batch action")
		return
	}
	switch parts[2] {
	case "review":
		s.handleReview(w, r, id)
	case "confirm":
		s.handleConfirm(w, r, id)
	case "publish":
		s.handlePublish(w, r, id)
	case "archive":
		s.handleArchive(w, r, id)
	case "export":
		s.handleExport(w, r, id)
	case "notes":
		s.handleNote(w, r, id)
	case "records":
		s.handleUpdate(w, r, id)
	default:
		writeErrorStatus(w, http.StatusNotFound, "unknown batch action")
	}
}

func (s *Server) decodeActor(r *http.Request) string {
	var req reviewRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	return req.Actor
}
func (s *Server) handleReview(w http.ResponseWriter, r *http.Request, id string) {
	batch, err := s.Service.StartReview(id, s.decodeActor(r))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}
func (s *Server) handleConfirm(w http.ResponseWriter, r *http.Request, id string) {
	batch, err := s.Service.ConfirmBatch(id, s.decodeActor(r))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}
func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request, id string) {
	batch, err := s.Service.PublishBatch(id, s.decodeActor(r))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}
func (s *Server) handleArchive(w http.ResponseWriter, r *http.Request, id string) {
	batch, err := s.Service.ArchiveBatch(id, s.decodeActor(r))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}
func (s *Server) handleExport(w http.ResponseWriter, r *http.Request, id string) {
	snapshot, err := s.Service.ExportConfirmed(id, s.decodeActor(r))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) handleNote(w http.ResponseWriter, r *http.Request, id string) {
	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}
	note, err := s.Service.AddNote(id, req.Author, req.Body)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, note)
}
func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, err.Error())
		return
	}
	batch, err := s.Service.UpdateRecord(id, req.RecordID, req.Actor, req.Number, req.Version)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, err error) {
	writeErrorStatus(w, http.StatusBadRequest, err.Error())
}
func writeErrorStatus(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
