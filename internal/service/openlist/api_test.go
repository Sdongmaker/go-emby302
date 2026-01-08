package openlist_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/config"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/service/openlist"
)

func TestFetch(t *testing.T) {
	cfg := os.Getenv("GE2O_CONFIG")
	if cfg == "" {
		cfg = "../../../config.yml"
	}
	if _, err := os.Stat(cfg); err != nil {
		t.Skipf("skip integration test: missing config file: %s", cfg)
	}

	err := config.ReadFromFile(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	var res openlist.FsList
	err = openlist.Fetch("/api/fs/list", http.MethodPost, nil, map[string]any{
		"refresh":  true,
		"password": "",
		"path":     "/",
	}, &res, true)
	if err != nil {
		t.Error(err)
		return
	}

	log.Printf("请求成功, data: %v", res)
}
