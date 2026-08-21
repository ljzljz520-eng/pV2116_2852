package importer

import (
	"encoding/json"
	"stickerchallenge/internal/domain"
)

func ParseJSON(data []byte) (Result, error) {
	var candidates []domain.Candidate
	if err := json.Unmarshal(data, &candidates); err != nil {
		return Result{}, err
	}
	return Result{Candidates: Normalize(candidates), Warnings: ValidateCandidates(candidates)}, nil
}

func ToJSON(result Result) ([]byte, error) { return json.Marshal(result) }
