package mp4s_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/mp4s"
)

func TestGenWithDuration(t *testing.T) {
	d := time.Hour*2 + time.Minute*23 + time.Second*21 + time.Millisecond*90
	bytes := mp4s.GenWithDuration(d)
	if len(bytes) == 0 {
		t.Fatal("generated mp4 is empty")
	}
	out := filepath.Join(t.TempDir(), "test.mp4")
	if err := os.WriteFile(out, bytes, os.ModePerm); err != nil {
		t.Fatal(err)
	}
}
