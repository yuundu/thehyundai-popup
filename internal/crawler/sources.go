package crawler

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yuundu/thehyundai-popup/internal/model"
)

// Source 는 수집 대상 한 곳. Parse 는 해당 사이트 DOM 을 우리 모델로 바꾼다.
type Source struct {
	Name    string
	URL     string
	Kind    model.SourceKind
	Enabled bool
	Parse   func(doc *goquery.Document, asOf time.Time, meta model.Source) []model.Popup
}

var (
	reWeekday  = regexp.MustCompile(`\([일월화수목금토]\)|\[[일월화수목금토]\]`)
	reDDay     = regexp.MustCompile(`D-\d+|D-DAY|오픈\s*예정|진행\s*중|새로\s*오픈|종료`)
	rePrice    = regexp.MustCompile(`(\d{1,3}(,\d{3})*)\s*원`)
	reNonWordy = regexp.MustCompile(`이용약관|개인정보|고객센터|로그인|회원가입|전체보기|더\s*보기|바로가기|공지사항`)
)

// DefaultSources 는 하루 한 번 도는 기본 수집 목록이다.
//
// 인스타그램(@thehyundai_seoul)과 네이버 예약은 로그인/차단 정책 때문에
// 자동 수집에서 제외하고, 대신 수동 확인 채널로 문서화한다(README 참고).
func DefaultSources() []Source {
	return []Source{
		{
			Name:    "더현대 서울 공식",
			URL:     "https://thehyundaiseoul.ehyundai.com/",
			Kind:    model.KindOfficial,
			Enabled: true,
			Parse:   parseEventCards,
		},
		{
			Name:    "현대백화점 이벤트",
			URL:     "https://ehyundai.com/newPortal/EV/EV000003_14_V.do?eventCd=B0349900&eventNo=37057",
			Kind:    model.KindOfficial,
			Enabled: true,
			Parse:   parseEventCards,
		},
		{
			Name:    "팝업나우",
			URL:     "https://popupnow.plentyer.com/area/thehyundai-seoul/",
			Kind:    model.KindPlatform,
			Enabled: true,
			Parse:   parsePopupNow,
		},
		{
			Name:    "팝플리",
			URL:     "https://poply.kr/",
			Kind:    model.KindPlatform,
			Enabled: false, // SPA: 렌더링 후 DOM 이 필요해 기본 비활성
			Parse:   parseEventCards,
		},
		{
			Name:    "디스위켄드",
			URL:     "https://www.disweekend.com/",
			Kind:    model.KindPlatform,
			Enabled: false,
			Parse:   parseEventCards,
		},
	}
}

// parseEventCards 는 "제목 + 기간"이 함께 있는 카드형 목록을 범용으로 긁는다.
// 공식몰/이벤트 페이지처럼 서버 렌더링 되는 목록에 쓴다.
func parseEventCards(doc *goquery.Document, asOf time.Time, meta model.Source) []model.Popup {
	var out []model.Popup
	seen := map[string]bool{}

	doc.Find("li, article, div.card, div.item, div.event, tr").Each(func(_ int, s *goquery.Selection) {
		txt := NormalizeSpace(s.Text())
		if len(txt) < 6 || len(txt) > 160 {
			return
		}
		if !strings.Contains(txt, "더현대") && FloorFromText(txt) == "" {
			return
		}
		if reNonWordy.MatchString(txt) {
			return
		}
		start, end, ok := ExtractRange(stripWeekday(txt), asOf)
		if !ok {
			return
		}
		name := cardTitle(txt)
		if name == "" || seen[Key(name)] {
			return
		}
		seen[Key(name)] = true

		p := model.Popup{
			ID:        Slug(name),
			Name:      name,
			Floor:     FloorFromText(txt),
			Category:  CategoryFromText(name + " " + txt),
			StartDate: start,
			EndDate:   end,
			Venue:     venueFromText(txt),
			Sources:   []model.Source{meta},
		}
		if m := rePrice.FindStringSubmatch(txt); m != nil {
			p.Price = m[0]
		}
		out = append(out, p)
	})
	return out
}

// parsePopupNow 는 "오키뷰티 더현대 서울 지하 1층 팝업 존 11.04 ~ 11.11 새로 오픈 D-5"
// 처럼 한 줄로 정리된 목록을 파싱한다. 상호명은 항목의 첫 링크 텍스트에서 가져와
// 이모지·브랜드 중복이 섞이는 걸 막는다.
func parsePopupNow(doc *goquery.Document, asOf time.Time, meta model.Source) []model.Popup {
	var out []model.Popup
	seen := map[string]bool{}

	doc.Find("li, article").Each(func(_ int, s *goquery.Selection) {
		txt := NormalizeSpace(s.Text())
		if !strings.Contains(txt, "더현대") || len(txt) > 160 {
			return
		}
		start, end, ok := ExtractRange(stripWeekday(txt), asOf)
		if !ok {
			return
		}
		// 종료된 행사는 오늘 목록에서 제외(연도/기간 필터 원칙)
		if end != "" {
			if e, err := time.ParseInLocation("2006-01-02", end, asOf.Location()); err == nil && e.Before(asOf) {
				return
			}
		}

		name := cleanName(NormalizeSpace(s.Find("a").First().Text()))
		if name == "" {
			name = cleanName(cardTitle(txt))
		}
		if name == "" && strings.Contains(txt, "더현대") {
			if i := strings.Index(txt, "더현대"); i > 0 {
				name = cleanName(strings.TrimSpace(txt[:i]))
			}
		}
		if name == "" || seen[Key(name)] {
			return
		}
		seen[Key(name)] = true

		p := model.Popup{
			ID:        Slug(name),
			Name:      name,
			Floor:     FloorFromText(txt),
			Category:  CategoryFromText(name + " " + txt),
			StartDate: start,
			EndDate:   end,
			Venue:     venueFromText(txt),
			Sources:   []model.Source{meta},
		}
		if strings.Contains(txt, "오픈 예정") {
			p.Tags = append(p.Tags, "upcoming")
		}
		out = append(out, p)
	})
	return out
}

// cleanName 은 상호명에서 지점 표기·상태 토큰·중복 토큰을 정리한다.
func cleanName(s string) string {
	for _, cut := range []string{"더현대 서울", "더현대서울", "더현대", "현대백화점"} {
		if i := strings.Index(s, cut); i >= 0 {
			s = s[:i]
		}
	}
	s = strings.TrimSpace(strings.Trim(s, "-·| :"))
	for _, suffix := range []string{"팝업스토어", "팝업 행사", "팝업 이벤트", " 팝업", "팝업", "여의도 오픈", "오픈 예정", "오픈"} {
		s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), suffix))
	}
	// 연속 중복 토큰 정리: "오키 오키 오키뷰티" → "오키 오키뷰티"
	toks := strings.Fields(s)
	out := toks[:0]
	for i, t := range toks {
		if i > 0 && toks[i-1] == t {
			continue
		}
		out = append(out, t)
	}
	s = strings.Join(out, " ")
	if n := len([]rune(s)); n < 2 || n > 30 {
		return ""
	}
	return s
}

// stripWeekday 는 날짜 표기 앞뒤의 요일/상태 토큰을 제거해 파싱률을 올린다.
func stripWeekday(s string) string {
	return strings.TrimSpace(reDDay.ReplaceAllString(reWeekday.ReplaceAllString(s, ""), ""))
}

// cardTitle 은 카드 텍스트에서 상호명만 뽑는다. 기간 표기 앞부분을 제목으로 본다.
func cardTitle(txt string) string {
	if m := reRangeMD.FindStringIndex(txt); m != nil {
		head := strings.TrimSpace(txt[:m[0]])
		// "브랜드명 더현대 서울 지하 1층" → 브랜드명만
		for _, cut := range []string{"더현대 서울", "더현대서울", "현대백화점", "지하 "} {
			if i := strings.Index(head, cut); i > 0 {
				head = strings.TrimSpace(head[:i])
			}
		}
		head = strings.TrimSpace(strings.Trim(head, "-·|"))
		if len([]rune(head)) >= 2 && len([]rune(head)) <= 30 {
			return head
		}
	}
	return ""
}

// venueFromText 는 층 표기 뒤쪽 문구를 장소로 쓴다.
func venueFromText(txt string) string {
	m := reFloor.FindStringIndex(stripWeekday(txt))
	if m == nil {
		return ""
	}
	tail := stripWeekday(txt)[m[0]:]
	if r := reRangeMD.FindStringIndex(tail); r != nil {
		tail = tail[:r[0]]
	}
	tail = strings.TrimSpace(reDDay.ReplaceAllString(tail, ""))
	if len([]rune(tail)) > 40 {
		tail = string([]rune(tail)[:40])
	}
	return tail
}

// CategoryFromText 는 키워드 사전으로 카테고리를 추정한다.
// 확신이 없으면 fashion(라이프)으로 두고, 틀렸을 때 티가 적은 쪽을 고른다.
func CategoryFromText(s string) string {
	low := strings.ToLower(s)
	switch {
	case matchAny(low, "타코야키", "닭강정", "에그타르트", "초콜릿", "마켓", "한우", "사과", "식품", "디저트", "베이커리", "커피", "food"):
		return model.CategoryFood
	case matchAny(low, "뷰티", "립밤", "스킨", "화장", "코스메", "beauty", "향수"):
		return model.CategoryBeauty
	case matchAny(low, "크리스마스", "빌리지", "시즌", "holiday", "christmas"):
		return model.CategorySeason
	}
	return model.CategoryFashion
}

func matchAny(s string, subs ...string) bool {
	for _, x := range subs {
		if strings.Contains(s, x) {
			return true
		}
	}
	return false
}
