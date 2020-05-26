package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/242617/core/config/source"
	"github.com/242617/core/config/source/file"
	"github.com/pkg/errors"
)

type Item struct {
	User struct {
		Name struct {
			First  string `env:"USER_FIRST_NAME"`
			Second string
		}
		Age     uint    `env:"USER_AGE"`
		Balance float64 `env:"USER_BALANCE" default:"10.25"`
		Active  bool    `env:"USER_ACTIVE" default:"true"`
	} `env:"USER"`
	Status  string        `yaml:"status_string" default:"ok"`
	Timeout time.Duration `env:"TIMEOUT" default:"10s"`
}

func TestDefault(t *testing.T) {
	var cfg Item

	config := New()
	if err := config.Scan(&cfg); err != nil {
		t.Fatal(errors.Wrap(err, "cannot scan config"))
	}

	if !cfg.User.Active {
		log.Fatalf("unexpected user activity: want %t, got %t", true, cfg.User.Active)
	}

	if cfg.Status != "ok" {
		log.Fatalf("unexpected status: want %q, got %q", "ok", cfg.Status)
	}

	if cfg.User.Balance != 10.25 {
		log.Fatalf("unexpected user balance: want %f, got %f", 10.25, cfg.User.Balance)
	}

	if cfg.Timeout != 10*time.Second {
		log.Fatalf("unexpected timeout: want %s, got %s", time.Duration(10*time.Second).String(), cfg.Timeout)
	}
}

func TestEnvBasic(t *testing.T) {
	for k, v := range map[string]string{
		"USER_FIRST_NAME": "Vasily",
		"USER_ACTIVE":     "true",
		"USER_AGE":        "30",
		"USER_BALANCE":    "2.5",
		"TIMEOUT":         "20s",
	} {
		if err := os.Setenv(k, v); err != nil {
			t.Fatal(errors.Wrap(err, "cannot send env"))
		}
	}

	var cfg Item

	config := New().With(source.Env())
	if err := config.Scan(&cfg); err != nil {
		t.Fatal(errors.Wrap(err, "cannot scan config"))
	}

	if cfg.User.Name.First != "Vasily" {
		log.Fatalf("unexpected user first name: want %q, got %q", "Vasily", cfg.User.Name.First)
	}

	if !cfg.User.Active {
		log.Fatalf("unexpected user activity: want %t, got %t", true, cfg.User.Active)
	}

	if cfg.User.Age != 30 {
		log.Fatalf("unexpected user age: want %d, got %d", 30, cfg.User.Age)
	}

	if cfg.User.Balance != 2.5 {
		log.Fatalf("unexpected balance: want %f, got %f", 2.5, cfg.User.Balance)
	}

	if cfg.Timeout != 20*time.Second {
		log.Fatalf("unexpected timeout: want %s, got %s", time.Duration(20*time.Second).String(), cfg.Timeout)
	}
}

func TestYAMLBasic(t *testing.T) {
	content := []byte(strings.Join([]string{
		"user:",
		"   name:",
		"       first: Ivan",
		"   active: true",
		"status_string: idle",
	}, "\n"))

	dir, err := ioutil.TempDir(os.TempDir(), "config")
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot create temp directory"))
	}
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, "config.yaml")
	if err := ioutil.WriteFile(filename, content, 0666); err != nil {
		t.Fatal(errors.Wrap(err, "cannot write file"))
	}

	var cfg Item

	config := New().With(file.YAML(filename))
	if err := config.Scan(&cfg); err != nil {
		t.Fatal(errors.Wrap(err, "cannot scan config"))
	}

	if cfg.User.Name.First != "Ivan" {
		log.Fatalf("unexpected user first name: want %q, got %q", "Vasily", cfg.User.Name.First)
	}

	if !cfg.User.Active {
		log.Fatalf("unexpected user activity: want %t, got %t", true, cfg.User.Active)
	}

	if cfg.User.Balance != 10.25 {
		log.Fatalf("unexpected user balance: want %f, got %f", 10.25, cfg.User.Balance)
	}

	if cfg.Status != "idle" {
		log.Fatalf("unexpected status: want %q, got %q", "idle", cfg.Status)
	}
}
