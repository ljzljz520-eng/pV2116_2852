package domain

import (
	"encoding/json"
	"fmt"
)

func EncodeBatch(batch Batch) ([]byte, error) { return json.Marshal(batch) }
func DecodeBatch(data []byte) (Batch, error) {
	var batch Batch
	err := json.Unmarshal(data, &batch)
	return batch, err
}
func EncodeRecord(record StickerRecord) ([]byte, error) { return json.Marshal(record) }
func DecodeRecord(data []byte) (StickerRecord, error) {
	var record StickerRecord
	err := json.Unmarshal(data, &record)
	return record, err
}
func EncodeEvent(event AuditEvent) ([]byte, error) { return json.Marshal(event) }
func DecodeEvent(data []byte) (AuditEvent, error) {
	var event AuditEvent
	err := json.Unmarshal(data, &event)
	return event, err
}
func ParseStatus(value string) (Status, error) {
	for _, status := range Statuses() {
		if string(status) == value {
			return status, nil
		}
	}
	return "", fmt.Errorf("unknown status %q", value)
}
func MustStatus(value string) Status {
	status, err := ParseStatus(value)
	if err != nil {
		return StatusRegistered
	}
	return status
}
