package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name     string
	DBPath   string
	Listen   string
	Divisors []int
}

func Profiles() []Profile {
	return []Profile{{Name: "local", DBPath: "stickerchallenge.db", Listen: ":8080", Divisors: []int{2, 3, 5, 7, 11}}, {Name: "review", DBPath: "review.db", Listen: ":8081", Divisors: []int{2, 3, 5}}, {Name: "audit", DBPath: "audit.db", Listen: ":8082", Divisors: []int{7, 11, 13}}}
}
func FindProfile(name string) (Profile, error) {
	for _, profile := range Profiles() {
		if profile.Name == name {
			return profile, nil
		}
	}
	return Profile{}, fmt.Errorf("profile %q not found", name)
}
func ResolvePath(base, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(base, path)
}
func NormalizeListen(address string) string {
	if strings.HasPrefix(address, ":") {
		return address
	}
	return ":" + address
}
func ProfileNames() []string {
	profiles := Profiles()
	names := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	return names
}
func DivisorText(profile Profile) string {
	values := make([]string, len(profile.Divisors))
	for index, divisor := range profile.Divisors {
		values[index] = fmt.Sprintf("%d", divisor)
	}
	return strings.Join(values, ",")
}
func IsSafePath(path string) bool { return path != "" && !strings.Contains(path, "..") }
func Merge(base Config, profile Profile) Config {
	if base.DBPath == "" {
		base.DBPath = profile.DBPath
	}
	if base.Listen == "" {
		base.Listen = profile.Listen
	}
	return base
}
