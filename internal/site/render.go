package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// 카드에 찍히는 상태 문구. 29CM 및 더현대 스타일에 맞춘 간결하고 감각적인 배지.
var statusBadge = map[string]string{
	model.StatusNew:      "🔥 오늘 오픈",
	model.StatusEnding:   "⏳ 종료 임박",
	model.StatusClosed:   "⏰ 오늘 마지막",
	model.StatusDeal:     "🛒 한정 특가",
	model.StatusConflict: "⚠️ 정보 상이",
}

var categoryChip = map[string]string{
	model.CategoryFood:    "식품/디저트",
	model.CategoryBeauty:  "뷰티",
	model.CategoryFashion: "패션/라이프",
	model.CategorySeason:  "시즌/크리스마스",
}

// FloorLabel 은 층 코드의 표기. data-floor 속성과 화면 라벨을 함께 돌려준다.
func FloorLabel(f string) string {
	switch f {
	case model.FloorB1:
		return "B1"
	case model.FloorB2:
		return "B2"
	case model.Floor1F:
		return "1F"
	case model.Floor4F:
		return "4F"
	case model.Floor5F:
		return "5F"
	}
	return "?"
}

// DateKO 는 2026-11-06 → 11.06(금) 형태로 바꾼다.
func DateKO(d string) string {
	t, err := time.ParseInLocation("2006-01-02", d, time.UTC)
	if err != nil {
		return d
	}
	wd := []string{"일", "월", "화", "수", "목", "금", "토"}[int(t.Weekday())]
	return fmt.Sprintf("%02d.%02d(%s)", int(t.Month()), t.Day(), wd)
}

// PeriodKO 는 기간 표기를 만든다. 종료 미정이면 "소진 시까지".
func PeriodKO(p model.Popup) string {
	s, e := DateKO(p.StartDate), DateKO(p.EndDate)
	switch {
	case p.StartDate != "" && p.EndDate != "":
		return s + " ~ " + e
	case p.StartDate != "":
		return s + " ~ 종료일 미정"
	case p.EndDate != "":
		return "~ " + e
	}
	return "기간 미확인"
}

// Badge 는 상태 코드의 화면 배지.
func Badge(code string) string { return statusBadge[code] }

// Chip 은 카테고리 코드의 화면 칩.
func Chip(code string) string {
	if v, ok := categoryChip[code]; ok {
		return v
	}
	return "패션/라이프"
}

// funcMap 을 템플릿에 물린다.
var funcMap = template.FuncMap{
	"floorLabel": FloorLabel,
	"dateKO":     DateKO,
	"periodKO":   PeriodKO,
	"badge":      Badge,
	"chip":       Chip,
	"safeHTML":   func(s string) template.HTML { return template.HTML(s) },
}

// pageData 는 index.html 전체의 입력.
type pageData struct {
	GeneratedAt string
	AsOfLabel   string
	Summary     string
	Floors      []model.Popup
}

// WriteIndex 는 정적 index.html 을 만든다.
func WriteIndex(out string, snap model.Snapshot, changes model.ChangeSet, newToday int) error {
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	tpl, err := template.New("index").Funcs(funcMap).Parse(indexTmpl)
	if err != nil {
		return fmt.Errorf("템플릿 파싱 실패: %w", err)
	}
	data := pageData{
		GeneratedAt: snap.GeneratedAt,
		AsOfLabel:   asOfLabel(snap.AsOf),
		Summary:     SummaryLine(changes, newToday),
		Floors:      snap.Popups,
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.Execute(f, data)
}

// asOfLabel 은 헤더의 기준 일시 표기.
func asOfLabel(asOf string) string {
	t, err := time.ParseInLocation("2006-01-02", asOf, time.UTC)
	if err != nil {
		return asOf
	}
	wd := []string{"일", "월", "화", "수", "목", "금", "토"}[int(t.Weekday())]
	return fmt.Sprintf("%s (%s)", asOf, wd)
}

// SummaryLine 은 상단 변동 요약. 오늘 오픈 건수와 어제 대비 변동을 문장으로 보여준다.
func SummaryLine(set model.ChangeSet, newToday int) string {
	parts := []string{}
	if newToday > 0 {
		parts = append(parts, fmt.Sprintf("오늘 오픈 %d건", newToday))
	}
	if len(set.Changes) == 0 {
		parts = append(parts, "어제와 동일")
		return strings.Join(parts, " · ")
	}
	var bits []string
	if n := set.Counts["added"]; n > 0 {
		bits = append(bits, fmt.Sprintf("신규 %d건", n))
	}
	if n := set.Counts["removed"]; n > 0 {
		bits = append(bits, fmt.Sprintf("종료 %d건", n))
	}
	if n := set.Counts["updated"]; n > 0 {
		bits = append(bits, fmt.Sprintf("변경 %d건", n))
	}
	parts = append(parts, "어제 대비 "+strings.Join(bits, ", "))
	return strings.Join(parts, " · ")
}
