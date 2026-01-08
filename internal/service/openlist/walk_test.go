package openlist_test

import (
	"log"
	"os"
	"testing"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/config"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/service/openlist"
)

func TestWalkFsList(t *testing.T) {
	cfg := os.Getenv("GE2O_CONFIG")
	if cfg == "" {
		cfg = "../../../config.yml"
	}
	if _, err := os.Stat(cfg); err != nil {
		t.Skipf("skip integration test: missing config file: %s", cfg)
	}

	err := config.ReadFromFile(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}

	walker := openlist.WalkFsList("/", 4)
	page, err := walker.Next()
	for err == nil {
		log.Println("page: ", page)
		page, err = walker.Next()
	}
	if err == openlist.ErrWalkEOF {
		return
	}
	t.Fatal(err)
}
