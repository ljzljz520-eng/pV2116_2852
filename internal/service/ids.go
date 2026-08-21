package service

import "fmt"

func BatchID(prefix string, number int) string { return fmt.Sprintf("%s-%02d", prefix, number) }
func RecordID(batch string, number int) string { return fmt.Sprintf("%s-record-%02d", batch, number) }
func ActorName(value string) string {
	if value == "" {
		return "system"
	}
	return value
}
func IsPublished(status string) bool { return status == "published" || status == "archived" }
