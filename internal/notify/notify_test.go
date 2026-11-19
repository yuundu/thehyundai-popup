package notify

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// 같은 내용이면 파일을 다시 쓰지 않아야 한다(수정 시각만 바뀌는 불필요한 커밋 방지).
func TestFileSkipsIdenticalRewrite(t *testing.T) {
	dir := t.TempDir()
	set := model.ChangeSet{
		AsOf:       "2026-11-19",
		PreviousOf: "2026-11-18",
		Counts:     map[string]int{"removed": 1},
		Changes: []model.Change{
			{Type: model.ChangeRemoved, ID: "x", Name: "X", Message: "🏁 종료: X"},
		},
	}
	if err := (File{Dir: dir}).Send(set); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "changes-2026-11-19.json")
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}

	if err := (File{Dir: dir}).Send(set); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.ModTime().After(old.Add(time.Minute)) {
		t.Fatalf("같은 내용인데 파일을 다시 씀: %v", fi.ModTime())
	}
}

// 내용이 달라지면 정상적으로 덮어쓴다.
func TestFileOverwritesWhenChanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "changes-2026-11-19.json")
	_ = os.WriteFile(path, []byte("{\"as_of\":\"2026-11-19\",\"changes\":[]}\n"), 0o644)

	set := model.ChangeSet{
		AsOf:    "2026-11-19",
		Counts:  map[string]int{"added": 1},
		Changes: []model.Change{{Type: model.ChangeAdded, ID: "y", Name: "Y", Message: "🆕 신규 오픈: Y"}},
	}
	if err := (File{Dir: dir}).Send(set); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) < 20 || !contains(string(b), "🆕 신규 오픈: Y") {
		t.Fatalf("갱신 실패: %s", b)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
