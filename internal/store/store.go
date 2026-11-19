// Package store 는 스냅샷 파일 입출력을 담당한다.
//
//	data/last_snapshot.json  : 이전 수집분 (다음 Diff 의 기준)
//	data/history/YYYY-MM-DD.json : 날짜별 보관본(추적/감사용)
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// LoadSnapshot 은 스냅샷 JSON 을 읽는다. 파일이 없으면 (nil, nil) 을 돌려
// "첫 실행" 상황으로 취급한다.
func LoadSnapshot(path string) (*model.Snapshot, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("스냅샷 읽기 실패 %s: %w", path, err)
	}
	var s model.Snapshot
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("스냅샷 파싱 실패 %s: %w", path, err)
	}
	return &s, nil
}

// SaveSnapshot 은 스냅샷을 원자적으로 쓴다(임시 파일 → rename).
// GitHub Actions 가 중간에 죽어도 반쯤 쓰인 JSON 이 남지 않게 하기 위함이다.
func SaveSnapshot(path string, s model.Snapshot) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// SaveHistory 는 날짜별 보관본을 남긴다. 같은 날 재실행하면 덮어쓴다.
func SaveHistory(dir string, s model.Snapshot) error {
	if dir == "" || s.AsOf == "" {
		return nil
	}
	return SaveSnapshot(filepath.Join(dir, s.AsOf+".json"), s)
}

// LatestSnapshotBefore 는 dir 안의 날짜별 보관본 중 asOf(YYYY-MM-DD)보다 이른 날짜의
// 가장 최근 스냅샷을 돌려준다. 같은 날 재실행에서도 "직전 수집일"을 기준으로 삼기 위한 것이다.
// 해당하는 파일이 없으면 (nil, "", nil).
func LatestSnapshotBefore(dir, asOf string) (*model.Snapshot, string, error) {
	if dir == "" || asOf == "" {
		return nil, "", nil
	}
	ents, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("히스토리 폴더 읽기 실패 %s: %w", dir, err)
	}
	best := ""
	for _, e := range ents {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		day := strings.TrimSuffix(e.Name(), ".json")
		if len(day) != 10 {
			continue
		}
		if _, err := time.ParseInLocation("2006-01-02", day, time.UTC); err != nil {
			continue
		}
		if day >= asOf {
			continue
		}
		if day > best {
			best = day
		}
	}
	if best == "" {
		return nil, "", nil
	}
	path := filepath.Join(dir, best+".json")
	s, err := LoadSnapshot(path)
	if err != nil {
		return nil, path, err
	}
	return s, path, nil
}

// LoadChangeSet 은 하루치 변동 파일(data/changes-YYYY-MM-DD.json)을 읽는다.
// 파일이 없으면 (nil, nil) — 첫 실행으로 취급한다.
func LoadChangeSet(path string) (*model.ChangeSet, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("변동 파일 읽기 실패 %s: %w", path, err)
	}
	var cs model.ChangeSet
	if err := json.Unmarshal(b, &cs); err != nil {
		return nil, fmt.Errorf("변동 파일 파싱 실패 %s: %w", path, err)
	}
	return &cs, nil
}

// WriteJSON 은 임의의 값을 들여쓰기 JSON 으로 저장한다(웹 배포용 산출물).
func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}
