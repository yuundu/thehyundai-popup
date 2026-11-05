package crawler

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 이 파일은 "사람이 읽는 표기"를 기계가 비교할 수 있는 값으로 바꾸는 순수 함수들만 담는다.
// 크롤링이 실패해도 날짜 비교/Diff 는 계속 동작해야 하므로 단위 테스트 대상이다.

var (
	reYMD     = regexp.MustCompile(`(\d{4})[.\-/](\d{1,2})[.\-/](\d{1,2})`)
	reMD      = regexp.MustCompile(`(\d{1,2})\s*[./]\s*(\d{1,2})`)
	reRangeMD = regexp.MustCompile(`(\d{1,2})\.(\d{1,2})\s*[~\-–]\s*(\d{1,2})\.(\d{1,2})`)
	reFloor   = regexp.MustCompile(`지하\s*(\d+)\s*층|(\d+)\s*층`)
	reSpaces  = regexp.MustCompile(`\s+`)
)

// NormalizeSpace 는 연속 공백/줄바꿈을 한 칸으로 줄인다.
func NormalizeSpace(s string) string {
	return strings.TrimSpace(reSpaces.ReplaceAllString(strings.ReplaceAll(s, "\u00a0", " "), " "))
}

// ParseDate 는 "2026-11-06", "2026.11.6", "2026/11/06" 을 time.Time 으로 바꾼다.
func ParseDate(s string) (time.Time, bool) {
	m := reYMD.FindStringSubmatch(s)
	if m == nil {
		return time.Time{}, false
	}
	y, _ := strconv.Atoi(m[1])
	mo, _ := strconv.Atoi(m[2])
	d, _ := strconv.Atoi(m[3])
	if mo < 1 || mo > 12 || d < 1 || d > 31 {
		return time.Time{}, false
	}
	return time.Date(y, time.Month(mo), d, 0, 0, 0, 0, time.UTC), true
}

// FormatDate 는 time.Time 을 스냅샷 표준 표기로 바꾼다.
func FormatDate(t time.Time) string { return t.Format("2006-01-02") }

// InferYear 는 연도 없는 "11.06" 표기에 연도를 붙인다.
//
// 규칙: 기준일보다 8개월 이상 과거로 보이면 내년 행사로 본다.
// (12월에 올라온 "01.05" 는 내년 1월 오픈으로 해석해야 한다.)
func InferYear(month, day int, asOf time.Time) time.Time {
	y := asOf.Year()
	cand := time.Date(y, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if cand.Before(asOf.AddDate(0, -8, 0)) {
		cand = time.Date(y+1, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	return cand
}

// ExtractRange 는 텍스트에서 기간을 뽑아 절대 날짜로 바꾼다.
// "2026-11-06 ~ 2026-11-12", "11.04 ~ 11.11", "~ 11.16" 을 모두 처리한다.
func ExtractRange(text string, asOf time.Time) (start, end string, ok bool) {
	// 1) 연도가 붙은 범위
	reYMDRange := regexp.MustCompile(`(\d{4})[.\-/](\d{1,2})[.\-/](\d{1,2})\s*[~\-–]\s*(\d{4})[.\-/](\d{1,2})[.\-/](\d{1,2})`)
	if m := reYMDRange.FindStringSubmatch(text); m != nil {
		s, ok1 := ParseDate(m[1] + "-" + m[2] + "-" + m[3])
		e, ok2 := ParseDate(m[4] + "-" + m[5] + "-" + m[6])
		if ok1 && ok2 {
			return FormatDate(s), FormatDate(e), true
		}
	}
	// 2) 연도 없는 범위
	if m := reRangeMD.FindStringSubmatch(text); m != nil {
		sm, _ := strconv.Atoi(m[1])
		sd, _ := strconv.Atoi(m[2])
		em, _ := strconv.Atoi(m[3])
		ed, _ := strconv.Atoi(m[4])
		if sm >= 1 && sm <= 12 && em >= 1 && em <= 12 {
			s := InferYear(sm, sd, asOf)
			e := time.Date(s.Year(), time.Month(em), ed, 0, 0, 0, 0, time.UTC)
			if e.Before(s) {
				e = e.AddDate(1, 0, 0)
			}
			return FormatDate(s), FormatDate(e), true
		}
	}
	// "~ 11.16" : 종료일만 있는 경우
	if m := regexp.MustCompile(`[~\-–]\s*(\d{1,2})\.(\d{1,2})`).FindStringSubmatch(text); m != nil {
		mo, _ := strconv.Atoi(m[1])
		d, _ := strconv.Atoi(m[2])
		if mo >= 1 && mo <= 12 {
			return "", FormatDate(InferYear(mo, d, asOf)), true
		}
	}
	return "", "", false
}

// FloorFromText 는 "지하 1층", "5층" 같은 표기를 층 코드로 바꾼다.
func FloorFromText(text string) string {
	m := reFloor.FindStringSubmatch(text)
	if m == nil {
		return ""
	}
	if m[1] != "" {
		switch m[1] {
		case "1":
			return "B1"
		case "2":
			return "B2"
		case "3":
			return "B3"
		case "4":
			return "B4"
		}
		return ""
	}
	switch m[2] {
	case "1":
		return "1F"
	case "2":
		return "2F"
	case "3":
		return "3F"
	case "4":
		return "4F"
	case "5":
		return "5F"
	case "6":
		return "6F"
	}
	return ""
}

// FloorLabelKO 는 층 코드를 한국어 표기로 되돌린다(알림 메시지용).
func FloorLabelKO(code string) string {
	switch code {
	case "B1":
		return "지하 1층"
	case "B2":
		return "지하 2층"
	case "1F":
		return "1층"
	case "5F":
		return "5층"
	}
	return "층 미확인"
}

// Slug 는 이름을 안정적인 ID 로 바꾼다. 한글은 그대로 두고 공백/기호만 정리한다.
func Slug(s string) string {
	s = strings.ToLower(NormalizeSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 0xAC00 && r <= 0xD7A3: // 한글 음절
			b.WriteRune(r)
		case r == ' ', r == '-', r == '_', r == '&':
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	if out == "" {
		out = "popup"
	}
	return out
}

// Key 는 병합용 정규화 키. 브랜드명 표기 차이(공백/기호)를 흡수한다.
func Key(name string) string {
	k := Slug(name)
	k = strings.ReplaceAll(k, "-", "")
	return k
}
