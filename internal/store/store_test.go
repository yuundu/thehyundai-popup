package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

func TestLatestSnapshotBefore(t *testing.T) {
	dir := t.TempDir()
	for _, day := range []string{"2026-11-13", "2026-11-17", "2026-11-18", "2026-11-19"} {
		if err := SaveSnapshot(filepath.Join(dir, day+".json"), model.Snapshot{AsOf: day}); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 오늘 11-19 를 돌 때 기준은 어제(11-18)여야 한다.
	// 같은 날 재실행으로 history/2026-11-19.json 이 있어도 오늘 것은 기준이 아니다.
	s, path, err := LatestSnapshotBefore(dir, "2026-11-19")
	if err != nil {
		t.Fatal(err)
	}
	if s == nil || s.AsOf != "2026-11-18" {
		t.Fatalf("기준일 선택 오류: %+v", s)
	}
	if filepath.Base(path) != "2026-11-18.json" {
		t.Fatalf("경로 오류: %s", path)
	}

	// 주말을 건너뛰는 월요일: 11-16 기준은 11-13.
	if s, _, _ := LatestSnapshotBefore(dir, "2026-11-16"); s == nil || s.AsOf != "2026-11-13" {
		t.Fatalf("주말 건너뛰기 오류: %+v", s)
	}

	// 더 이른 날짜가 없으면 nil — 첫 실행.
	if s, _, err := LatestSnapshotBefore(dir, "2026-11-06"); err != nil || s != nil {
		t.Fatalf("첫 실행 처리 오류: %+v %v", s, err)
	}

	// 폴더가 없어도 죽지 않는다.
	if s, _, err := LatestSnapshotBefore(filepath.Join(dir, "nope"), "2026-11-19"); err != nil || s != nil {
		t.Fatalf("폴더 없음 처리 오류: %+v %v", s, err)
	}
}

func TestLoadChangeSetMissing(t *testing.T) {
	cs, err := LoadChangeSet(filepath.Join(t.TempDir(), "changes-2026-11-19.json"))
	if err != nil || cs != nil {
		t.Fatalf("없는 파일은 (nil, nil): %+v %v", cs, err)
	}
}
