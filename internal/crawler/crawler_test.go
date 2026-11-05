package crawler

import (
	"testing"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

func TestExtractRange(t *testing.T) {
	asOf := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		in         string
		start, end string
	}{
		{"11.04 ~ 11.11", "2026-11-04", "2026-11-11"},
		{"2026-11-06 ~ 2026-11-12", "2026-11-06", "2026-11-12"},
		{"10.20~11.16", "2026-10-20", "2026-11-16"},
		{"~ 11.16", "", "2026-11-16"},
	}
	for _, c := range cases {
		s, e, ok := ExtractRange(c.in, asOf)
		if !ok || s != c.start || e != c.end {
			t.Errorf("ExtractRange(%q) = %q,%q,%v; want %q,%q", c.in, s, e, ok, c.start, c.end)
		}
	}
}

func TestFloorFromText(t *testing.T) {
	if got := FloorFromText("더현대 서울 지하 1층 팝업 존"); got != "B1" {
		t.Errorf("B1 변환 실패: %q", got)
	}
	if got := FloorFromText("5층 사운즈 포레스트"); got != "5F" {
		t.Errorf("5F 변환 실패: %q", got)
	}
}

func TestInferYear(t *testing.T) {
	asOf := time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC)
	got := InferYear(1, 5, asOf)
	if got.Year() != 2027 {
		t.Errorf("12월의 01.05 는 내년이어야 함: %v", got)
	}
}

func TestSlugAndKey(t *testing.T) {
	if Slug("강릉감자 닭강정!") != "강릉감자-닭강정" {
		t.Errorf("Slug 오류: %q", Slug("강릉감자 닭강정!"))
	}
	if Key("오키뷰티") != Key("오키뷰티 ") {
		t.Errorf("Key 정규화 실패")
	}
}

func TestMergeDedupAndPriority(t *testing.T) {
	base := model.Popup{ID: "oki", Name: "오키뷰티", Floor: "B1", StartDate: "2026-11-04", EndDate: "2026-11-11"}
	items := []Tagged{
		{Popup: base, Kind: model.KindPlatform},
		{Popup: base, Kind: model.KindSNS},
	}
	out := Merge(items)
	if len(out) != 1 {
		t.Fatalf("병합 실패: %d건", len(out))
	}
	if out[0].Confidence < 60 {
		t.Fatalf("신뢰도 산정 오류: %d", out[0].Confidence)
	}
}
