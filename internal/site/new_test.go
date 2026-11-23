package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yuundu/thehyundai-popup/internal/theedit"
)

func TestWriteNewPage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "new", "index.html")
	arch := theedit.Archive{
		GeneratedAt: "2026-11-23 16:00 KST",
		TotalItems:  2,
		Weeks: []theedit.Week{
			{
				ID: 2, Title: "[새로나왔] 11월 4주차 신상 리스트", Week: "11월 4주차",
				Date: "2026-11-23", Link: "https://the-edit.co.kr/89435", New: true,
				Intro: "이번 주 도입 문단.", Items: []theedit.Item{
					{Name: "테스트 아이템", Image: "https://the-edit.co.kr/a.webp", Summary: "설명 문단.",
						Links: []theedit.Link{{Label: "스마트스토어", URL: "https://smartstore.naver.com/x"}}},
					{Name: "링크 없는 아이템", Summary: "링크가 없는 설명."},
				},
			},
			{
				ID: 1, Title: "[새로나왔] 11월 3주차 신상 리스트", Week: "11월 3주차",
				Date: "2026-11-16", Link: "https://the-edit.co.kr/89328", Items: []theedit.Item{
					{Name: "지난 주 아이템", Summary: "지난 주 설명."},
				},
			},
		},
	}
	if err := WriteNewPage(out, arch); err != nil {
		t.Fatalf("페이지 생성 실패: %v", err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	for _, want := range []string{
		`href="../index.html"`,          // 더현대 팝업으로 돌아가는 상대경로
		`href="index.html" class="active"`,
		"ISSUE 11월 4주차",
		"NEW · 이번 주 업데이트",
		"테스트 아이템",
		`loading="lazy"`,
		`href="https://smartstore.naver.com/x"`,
		"링크 없는 아이템",
		`href="https://the-edit.co.kr/89435" target="_blank" rel="noopener">원문에서 보기`,
		`data-week="1"`,
		`hidden`, // 지난 주차는 처음에 접혀 있어야 한다
	} {
		if !strings.Contains(s, want) {
			t.Errorf("결과 HTML 에 %q 가 없습니다", want)
		}
	}

	// 첫 빌드/네트워크 실패로 비어 있어도 페이지는 나와야 한다.
	empty := filepath.Join(dir, "empty", "index.html")
	if err := WriteNewPage(empty, theedit.Archive{}); err != nil {
		t.Fatalf("빈 아카이브 페이지 생성 실패: %v", err)
	}
	eb, _ := os.ReadFile(empty)
	if !strings.Contains(string(eb), "아직 수집된 [새로나왔] 글이 없어요") {
		t.Error("빈 아카이브 안내 문구가 없습니다")
	}
}
