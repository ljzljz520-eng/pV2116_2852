package config

import (
	"flag"
	"fmt"
)

type Config struct {
	DBPath string
	Listen string
	Demo   bool
}

func Default() Config { return Config{DBPath: "stickerchallenge.db", Listen: ":8080"} }

func Load(args []string) (Config, error) {
	cfg := Default()
	flags := flag.NewFlagSet("stickerctl", flag.ContinueOnError)
	flags.StringVar(&cfg.DBPath, "db", cfg.DBPath, "database path")
	flags.StringVar(&cfg.Listen, "listen", cfg.Listen, "listen address")
	flags.BoolVar(&cfg.Demo, "demo", false, "run deterministic demo")
	if err := flags.Parse(args); err != nil {
		return cfg, err
	}
	if cfg.DBPath == "" || cfg.Listen == "" {
		return cfg, fmt.Errorf("db and listen are required")
	}
	return cfg, nil
}

func Usage() string          { return "stickerctl [-db path] [-listen address] [-demo]" }
func IsDemo(cfg Config) bool { return cfg.Demo }
func Validate(cfg Config) error {
	if cfg.DBPath == "" || cfg.Listen == "" {
		return fmt.Errorf("invalid config")
	}
	return nil
}
