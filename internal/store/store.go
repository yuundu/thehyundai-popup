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
