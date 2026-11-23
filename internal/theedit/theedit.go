// Package theedit 는 디에디트(THE EDIT)의 주간 신상 시리즈 [새로나왔] 을 수집·파싱한다.
//
//	https://the-edit.co.kr/wp-json/wp/v2/posts?search=새로나왔
//
// 한 주차 글은 아래처럼 <hr> 구분선으로 아이템이 나뉜다.
//
//	<figure class="wp-block-image"><img src="..."/></figure>
//	<p><strong>아이템명</strong></p>
//	<p>설명 ... <a href="...">구매처</a> ...</p>
//
// 글마다 figure/p 순서나 이미지 유무, 문단 수가 다를 수 있어 순서에 기대지 않고
// "구분선으로 나뉜 덩어리" 단위로 유연하게 파싱한다. 네트워크가 실패해도 기존
// 아카이브(data/theedit-archive.json)를 그대로 써서 빌드가 깨지지 않게 한다.
package theedit

import (
	"context"
	"encoding/json"
	"fmt"
	stdhtml "html"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	xhtml "golang.org/x/net/html"
)

// APIURL 은 [새로나왔] 글을 최신순으로 주는 WordPress REST API 주소.
const APIURL = "https://the-edit.co.kr/wp-json/wp/v2/posts?search=%EC%83%88%EB%A1%9C%EB%82%98%EC%99%94&per_page=10"

// Source 는 아카이브에 기록하는 출처 표기.
const Source = "디에디트 THE EDIT · [새로나왔]"

const kstOffset = 9 * 60 * 60

// Link 는 본문에서 뽑아낸 구매·참고 링크. Label 은 앵커 텍스트("스마트스토어" 등).
type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Item 은 신상 한 개.
type Item struct {
	Name    string `json:"name"`
	Image   string `json:"image,omitempty"`
	Summary string `json:"summary,omitempty"`
	Links   []Link `json:"links,omitempty"`
}

// Week 는 [새로나왔] 한 주차 글.
type Week struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`           // "[새로나왔] 11월 4주차 신상 리스트"
	Week         string `json:"week,omitempty"`  // "11월 4주차"
	Date         string `json:"date"`            // "2026-11-23"
	Link         string `json:"link"`            // 원문 주소
	Intro        string `json:"intro,omitempty"` // 도입 문단
	New          bool   `json:"new,omitempty"`   // 이번 빌드에서 처음 발견
	DiscoveredAt string `json:"discovered_at,omitempty"`
	Items        []Item `json:"items"`
}

// Archive 는 누적 보관본. data/theedit-archive.json 으로 저장된다.
type Archive struct {
	GeneratedAt string `json:"generated_at"`
	Source      string `json:"source"`
	TotalItems  int    `json:"total_items"`
	Weeks       []Week `json:"weeks"` // 최신 주차가 앞
}

type apiPost struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Link  string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
}

// Fetch 는 최근 [새로나왔] 글들을 수집해 주차별로 돌려준다. 최신 주차가 앞.
func Fetch(ctx context.Context) ([]Week, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "thehyundai-popup/1.0 (+https://yuundu.github.io/thehyundai-popup/)")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("디에디트 API 요청 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("디에디트 API 응답 %d", resp.StatusCode)
	}

	var posts []apiPost
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, fmt.Errorf("디에디트 API 파싱 실패: %w", err)
	}

	weeks := make([]Week, 0, len(posts))
	for _, p := range posts {
		title := Norm(stdhtml.UnescapeString(p.Title.Rendered))
		if !strings.Contains(title, "[새로나왔]") {
			continue
		}
		intro, items := ParseContent(p.Content.Rendered)
		if len(items) == 0 {
			// 본문 구조가 예상과 다른 글은 비어 있는 카드로 남기지 않는다.
			continue
		}
		date := p.Date
		if len(date) > 10 {
			date = date[:10]
		}
		w := Week{
			ID:    p.ID,
			Title: title,
			Week:  WeekLabel(title),
			Date:  date,
			Link:  p.Link,
			Intro: intro,
			Items: items,
		}
		weeks = append(weeks, w)
	}
	if len(weeks) == 0 {
		return nil, fmt.Errorf("디에디트 [새로나왔] 글을 찾지 못했습니다")
	}
	SortWeeks(weeks)
	return weeks, nil
}

// ParseContent 는 본문 HTML 을 (도입 문단, 아이템 목록) 으로 나눈다.
func ParseContent(content string) (intro string, items []Item) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", nil
	}
	body := doc.Find("body").First()

	var segments [][]*xhtml.Node
	var cur []*xhtml.Node
	body.Children().Each(func(_ int, s *goquery.Selection) {
		n := s.Get(0)
		if n == nil {
			return
		}
		if n.Type == xhtml.ElementNode && n.Data == "hr" {
			segments = append(segments, cur)
			cur = nil
			return
		}
		cur = append(cur, n)
	})
	segments = append(segments, cur)

	for i, seg := range segments {
		it, ok := parseSegment(seg)
		if !ok {
			// 첫 덩어리(구분선 앞)는 도입 문단이다.
			if i == 0 && intro == "" {
				intro = segmentText(seg)
			}
			continue
		}
		items = append(items, it)
	}
	return intro, items
}

// parseSegment 는 <hr> 로 나뉜 한 덩어리에서 아이템 하나를 꺼낸다.
func parseSegment(nodes []*xhtml.Node) (Item, bool) {
	var it Item
	var descParts []string

	for _, n := range nodes {
		if n.Type != xhtml.ElementNode {
			continue
		}
		sel := goquery.NewDocumentFromNode(n).Selection

		switch n.Data {
		case "figure", "div":
			// 이미지: 덩어리 안 첫 이미지를 썸네일로.
			if it.Image == "" {
				if src, ok := firstImage(sel); ok {
					it.Image = src
				}
			}
		case "p":
			txt := Norm(sel.Text())
			if txt == "" {
				continue
			}
			if isFooter(txt) {
				continue
			}
			if it.Name == "" && isTitleParagraph(sel) {
				it.Name = strongText(sel)
				continue
			}
			descParts = append(descParts, txt)
			collectLinks(sel, &it.Links)
		case "h2", "h3", "h4":
			if it.Name == "" {
				if t := Norm(sel.Text()); t != "" {
					it.Name = t
					continue
				}
			}
			descParts = append(descParts, Norm(sel.Text()))
		default:
			// 그 밖의 블록: 이미지가 있으면 썸네일 후보로만 사용.
			if it.Image == "" {
				if src, ok := firstImage(sel); ok {
					it.Image = src
				}
			}
		}
	}

	it.Summary = strings.Join(descParts, "\n\n")
	it.Links = dedupeLinks(it.Links)

	if it.Name == "" || (it.Image == "" && it.Summary == "") {
		return Item{}, false
	}
	return it, true
}

// isTitleParagraph 는 <p><strong>아이템명</strong></p> 형태인지 본다.
func isTitleParagraph(sel *goquery.Selection) bool {
	if sel.Find("img").Length() > 0 {
		return false
	}
	strong := sel.Find("strong")
	if strong.Length() != 1 {
		return false
	}
	// 강조 안에 <em> 이 있으면 안내문(구독 권유 등)이다.
	if strong.Find("em").Length() > 0 {
		return false
	}
	st := Norm(strong.Text())
	if st == "" {
		return false
	}
	return st == Norm(sel.Text())
}

func strongText(sel *goquery.Selection) string {
	return Norm(sel.Find("strong").First().Text())
}

// isFooter 는 모든 글 끝에 붙는 뉴스레터 안내 문단.
func isFooter(txt string) bool {
	switch {
	case strings.Contains(txt, "까탈로그"),
		strings.Contains(txt, "구독하자"),
		strings.Contains(txt, "세상은 넓고 신제품은 많다"):
		return true
	}
	return false
}

func firstImage(sel *goquery.Selection) (string, bool) {
	img := sel.Find("img").First()
	if img.Length() == 0 {
		return "", false
	}
	src, ok := img.Attr("src")
	src = strings.TrimSpace(src)
	if !ok || src == "" {
		if ds, ok := img.Attr("data-src"); ok {
			src = strings.TrimSpace(ds)
		}
	}
	if !strings.HasPrefix(src, "http") {
		return "", false
	}
	return src, true
}

func collectLinks(sel *goquery.Selection, out *[]Link) {
	sel.Find("a").Each(func(_ int, a *goquery.Selection) {
		href, ok := a.Attr("href")
		href = strings.TrimSpace(href)
		if !ok || !strings.HasPrefix(href, "http") {
			return
		}
		if strings.Contains(href, "the-edit.co.kr/newsletter") {
			return
		}
		label := Norm(a.Text())
		if label == "" {
			label = hostLabel(href)
		}
		*out = append(*out, Link{Label: label, URL: href})
	})
}

func hostLabel(raw string) string {
	s := raw
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	if i := strings.IndexAny(s, "/?#"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimPrefix(s, "www.")
}

func dedupeLinks(in []Link) []Link {
	if len(in) == 0 {
		return nil
	}
	seen := map[string]bool{}
	out := make([]Link, 0, len(in))
	for _, l := range in {
		key := l.URL
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, l)
	}
	return out
}

func segmentText(nodes []*xhtml.Node) string {
	var parts []string
	for _, n := range nodes {
		if n.Type != xhtml.ElementNode {
			continue
		}
		if t := Norm(goquery.NewDocumentFromNode(n).Text()); t != "" && !isFooter(t) {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, " ")
}

var weekRe = regexp.MustCompile(`(\d{1,2})\s*월\s*(\d{1,2})\s*주차`)

// WeekLabel 은 제목에서 "11월 4주차" 를 뽑는다. 없으면 제목 그대로.
func WeekLabel(title string) string {
	m := weekRe.FindStringSubmatch(title)
	if len(m) == 3 {
		return m[1] + "월 " + m[2] + "주차"
	}
	return title
}

// Norm 은 화면에 쓸 텍스트를 정리한다(NBSP·제로폭 문자 제거, 공백 축약).
func Norm(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.ReplaceAll(s, "\u200b", "")
	s = strings.ReplaceAll(s, "\u2009", " ")
	return strings.Join(strings.Fields(s), " ")
}

// SortWeeks 는 최신 주차가 앞에 오도록 정렬한다.
func SortWeeks(ws []Week) {
	sort.SliceStable(ws, func(i, j int) bool {
		if ws[i].Date != ws[j].Date {
			return ws[i].Date > ws[j].Date
		}
		return ws[i].ID > ws[j].ID
	})
}

// NewBadgeWindow 는 NEW 배지를 유지하는 기간.
const NewBadgeWindow = 72 * time.Hour

// FreshDate 는 "2026-11-23" 발행일이 now 기준 window 안에 있는지 본다(KST 0시 기준).
func FreshDate(date string, now time.Time, window time.Duration) bool {
	t, err := time.ParseInLocation("2006-01-02", date, time.FixedZone("KST", kstOffset))
	if err != nil {
		return false
	}
	d := now.Sub(t)
	return d <= window && d >= -24*time.Hour
}

// Load 는 아카이브를 읽는다. 파일이 없으면 (nil, nil).
func Load(path string) (*Archive, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("디에디트 아카이브 읽기 실패 %s: %w", path, err)
	}
	var a Archive
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, fmt.Errorf("디에디트 아카이브 파싱 실패 %s: %w", path, err)
	}
	return &a, nil
}

// Merge 는 기존 아카이브와 새로 수집한 주차를 합친다.
// 이번에 처음 본 주차는 New=true 로 표시하고, 그동안의 NEW 표시는 내린다.
func Merge(old *Archive, fetched []Week, now time.Time) (Archive, []int) {
	stamp := now.In(time.FixedZone("KST", kstOffset)).Format(time.RFC3339)

	byID := map[int]Week{}
	if old != nil {
		for _, w := range old.Weeks {
			w.New = false
			byID[w.ID] = w
		}
	}

	var discovered []int
	for _, w := range fetched {
		if prev, ok := byID[w.ID]; ok {
			w.DiscoveredAt = prev.DiscoveredAt
			w.New = false
		} else {
			w.DiscoveredAt = stamp
			w.New = true
			discovered = append(discovered, w.ID)
		}
		if w.DiscoveredAt == "" {
			w.DiscoveredAt = stamp
		}
		byID[w.ID] = w
	}

	weeks := make([]Week, 0, len(byID))
	for _, w := range byID {
		weeks = append(weeks, w)
	}
	SortWeeks(weeks)

	// NEW 배지는 가장 최신 주차가 발행 3일 이내일 때만 붙인다.
	// 매일 다시 빌드해도 다음 주차가 올라오기 전까지 "이번 주 새 소식"임을 알 수 있다.
	for i := range weeks {
		weeks[i].New = false
	}
	if len(weeks) > 0 && FreshDate(weeks[0].Date, now, NewBadgeWindow) {
		weeks[0].New = true
	}

	total := 0
	for _, w := range weeks {
		total += len(w.Items)
	}

	return Archive{
		GeneratedAt: stamp,
		Source:      Source,
		TotalItems:  total,
		Weeks:       weeks,
	}, discovered
}
