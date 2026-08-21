package config

import "testing"

func TestLoad(t *testing.T) {
	cfg, err := Load([]string{"-db", "data.db", "-listen", ":9000", "-demo"})
	if err != nil || cfg.DBPath != "data.db" || !cfg.Demo {
		t.Fatalf("%+v %v", cfg, err)
	}
	if err := Validate(Default()); err != nil {
		t.Fatal(err)
	}
}
