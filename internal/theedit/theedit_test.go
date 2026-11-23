package theedit

import (
	"os"
	"strings"
	"testing"
	"time"
)

// 11월 4주차 최신 글(https://the-edit.co.kr/89435) 실제 본문으로 파싱을 검증한다.
func TestParseLatestWeek(t *testing.T) {
	b, err := os.ReadFile("testdata/latest-89435.html")
	if err != nil {
		t.Fatalf("fixture 읽기 실패: %v", err)
	}
	intro, items := ParseContent(string(b))

	if len(items) != 12 {
		t.Fatalf("아이템 12개 기대, %d개 파싱됨", len(items))
	}
	if !strings.Contains(intro, "2027년 달력") {
		t.Errorf("도입 문단이 잡히지 않음: %q", intro)
	}

	first := items[0]
	if first.Name != "에이든 대한민국 박물관, 미술관 여권 스탬프북" {
		t.Errorf("첫 아이템명: %q", first.Name)
	}
	if !strings.HasPrefix(first.Image, "https://the-edit.co.kr/wp-content/uploads/") {
		t.Errorf("첫 아이템 이미지: %q", first.Image)
	}
	if !strings.Contains(first.Summary, "251곳") {
		t.Errorf("첫 아이템 설명 누락: %q", first.Summary)
	}
	if len(first.Links) != 1 || first.Links[0].Label != "스마트스토어" {
		t.Errorf("첫 아이템 링크: %+v", first.Links)
	}
	if strings.Contains(first.Links[0].URL, "&amp;") {
		t.Errorf("링크 HTML 엔티티가 남음: %s", first.Links[0].URL)
	}

	// 본문 중간의 여러 링크(알라딘·YES24)가 모두 뽑혀야 한다.
	var vinyl *Item
	for i := range items {
		if strings.Contains(items[i].Name, "검정치마") {
			vinyl = &items[i]
		}
	}
	if vinyl == nil {
		t.Fatal("검정치마 바이닐 아이템을 찾지 못함")
	}
	if len(vinyl.Links) != 2 {
		t.Errorf("검정치마 링크 2개 기대, %+v", vinyl.Links)
	}

	// 뉴스레터 안내문은 아이템으로 새지 않아야 한다.
	for _, it := range items {
		if strings.Contains(it.Name, "까탈로그") || strings.Contains(it.Name, "세상은 넓고") {
			t.Errorf("안내문이 아이템으로 파싱됨: %q", it.Name)
		}
		if it.Name == "" || it.Image == "" {
			t.Errorf("이름/이미지 누락: %+v", it)
		}
	}
}

// 순서·이미지 유무·문단 수가 다른 변형 본문도 견뎌야 한다.
func TestParseFlexibleLayout(t *testing.T) {
	body := `
<p>도입 문단입니다.</p>
<hr class="wp-block-separator"/>
<p><strong>이미지 없는 아이템</strong></p>
<p>설명 첫 문단.</p>
<p>설명 둘째 문단. <a href="https://shop.example.com/item?x=1&amp;y=2">구매처</a></p>
<hr/>
<p><strong>이미지가 먼저 오지 않는 아이템</strong></p>
<figure class="wp-block-image"><img src="https://cdn.example.com/a.webp" alt="브랜드"/></figure>
<p>설명. <a href="https://shop.example.com/a">A몰</a> <a href="https://shop.example.com/a">A몰 중복</a></p>
<hr/>
<p><strong><em>구독 안내 <a href="https://the-edit.co.kr/newsletter">구독</a></em></strong></p>
`
	intro, items := ParseContent(body)
	if intro != "도입 문단입니다." {
		t.Errorf("도입 문단: %q", intro)
	}
	if len(items) != 2 {
		t.Fatalf("아이템 2개 기대, %d개: %+v", len(items), items)
	}
	if items[0].Name != "이미지 없는 아이템" || items[0].Image != "" {
		t.Errorf("첫 아이템: %+v", items[0])
	}
	if items[0].Summary != "설명 첫 문단.\n\n설명 둘째 문단. 구매처" {
		t.Errorf("여러 문단 설명: %q", items[0].Summary)
	}
	if len(items[0].Links) != 1 || items[0].Links[0].URL != "https://shop.example.com/item?x=1&y=2" {
		t.Errorf("링크 추출/디코딩: %+v", items[0].Links)
	}
	if items[1].Image != "https://cdn.example.com/a.webp" {
		t.Errorf("figure 가 뒤에 와도 이미지 추출: %+v", items[1])
	}
	if len(items[1].Links) != 1 {
		t.Errorf("중복 링크 제거: %+v", items[1].Links)
	}
}

func TestWeekLabel(t *testing.T) {
	cases := map[string]string{
		"[새로나왔] 11월 4주차 신상 리스트":  "11월 4주차",
		"[새로나왔] 9월 5주차 신상  리스트": "9월 5주차",
		"[새로나왔] 신상 리스트":          "[새로나왔] 신상 리스트",
	}
	for in, want := range cases {
		if got := WeekLabel(Norm(in)); got != want {
			t.Errorf("WeekLabel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMergeArchive(t *testing.T) {
	now := time.Date(2026, 11, 23, 16, 0, 0, 0, time.FixedZone("KST", kstOffset))
	fetched := []Week{
		{ID: 2, Date: "2026-11-23", Title: "[새로나왔] 11월 4주차", Items: []Item{{Name: "B"}}},
		{ID: 1, Date: "2026-11-16", Title: "[새로나왔] 11월 3주차", Items: []Item{{Name: "A"}}},
	}

	// 첫 빌드: 가장 최신 주차만 NEW.
	first, discovered := Merge(nil, fetched, now)
	if len(first.Weeks) != 2 || first.Weeks[0].ID != 2 {
		t.Fatalf("정렬/병합 결과: %+v", first.Weeks)
	}
	if !first.Weeks[0].New || first.Weeks[1].New {
		t.Errorf("첫 빌드 NEW 표시: %+v", first.Weeks)
	}
	if len(discovered) != 1 || discovered[0] != 2 {
		t.Errorf("첫 빌드 발견 목록: %v", discovered)
	}

	// 다음 빌드에서 새 주차가 올라오면 새 주차만 NEW.
	fetched2 := append([]Week{{ID: 3, Date: "2026-11-30", Title: "[새로나왔] 11월 5주차", Items: []Item{{Name: "C"}}}}, fetched...)
	next, discovered := Merge(&first, fetched2, now.AddDate(0, 0, 7))
	if len(next.Weeks) != 3 {
		t.Fatalf("누적 주차: %d", len(next.Weeks))
	}
	if !next.Weeks[0].New || next.Weeks[0].ID != 3 {
		t.Errorf("새 주차 NEW: %+v", next.Weeks[0])
	}
	for _, w := range next.Weeks[1:] {
		if w.New {
			t.Errorf("이전 주차 NEW 가 내려가지 않음: %+v", w)
		}
	}
	if len(discovered) != 1 || discovered[0] != 3 {
		t.Errorf("발견 목록: %v", discovered)
	}
	if next.TotalItems != 3 {
		t.Errorf("총 아이템: %d", next.TotalItems)
	}

	// API 가 잠깐 비어도 기존 아카이브가 유지되어야 한다.
	kept, _ := Merge(&next, nil, now.AddDate(0, 0, 8))
	if len(kept.Weeks) != 3 || kept.TotalItems != 3 {
		t.Errorf("빈 수집 시 보관본 유지 실패: %+v", kept)
	}
}

func TestLoadMissing(t *testing.T) {
	a, err := Load("testdata/does-not-exist.json")
	if err != nil || a != nil {
		t.Fatalf("없는 파일: %v %v", a, err)
	}
}
