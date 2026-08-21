package domain

import "errors"

type Status string

const (
	StatusRegistered Status = "registered"
	StatusReviewing  Status = "reviewing"
	StatusConfirmed  Status = "confirmed"
	StatusPublished  Status = "published"
	StatusArchived   Status = "archived"
)

var (
	ErrNotFound   = errors.New("entity not found")
	ErrInvalid    = errors.New("invalid entity")
	ErrConflict   = errors.New("version conflict")
	ErrTransition = errors.New("invalid lifecycle transition")
)

type StickerRecord struct {
	ID        string `json:"id"`
	BatchID   string `json:"batch_id"`
	Number    int    `json:"number"`
	Divisors  []int  `json:"divisors"`
	Result    string `json:"result"`
	Confirmed bool   `json:"confirmed"`
	UpdatedBy string `json:"updated_by"`
}

type Batch struct {
	ID        string          `json:"id"`
	Label     string          `json:"label"`
	Status    Status          `json:"status"`
	Owner     string          `json:"owner"`
	Records   []StickerRecord `json:"records"`
	Version   int             `json:"version"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

type AuditEvent struct {
	ID      string `json:"id"`
	BatchID string `json:"batch_id"`
	Action  string `json:"action"`
	Actor   string `json:"actor"`
	Detail  string `json:"detail"`
	At      string `json:"at"`
}

type CollaborationNote struct {
	ID      string `json:"id"`
	BatchID string `json:"batch_id"`
	Author  string `json:"author"`
	Body    string `json:"body"`
	At      string `json:"at"`
}

type ExportSnapshot struct {
	ID        string `json:"id"`
	BatchID   string `json:"batch_id"`
	Format    string `json:"format"`
	Payload   string `json:"payload"`
	CreatedBy string `json:"created_by"`
	At        string `json:"at"`
}

type Candidate struct {
	ID     string
	Number int
}

type SearchQuery struct {
	Label  string
	Status Status
}

type Summary struct {
	BatchID   string `json:"batch_id"`
	Label     string `json:"label"`
	Total     int    `json:"total"`
	Confirmed int    `json:"confirmed"`
	Passing   int    `json:"passing"`
	Failing   int    `json:"failing"`
	Lifecycle Status `json:"lifecycle"`
}
