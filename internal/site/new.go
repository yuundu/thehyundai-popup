package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuundu/thehyundai-popup/internal/theedit"
)

// newFuncMap 은 [새로나왔] 페이지 템플릿용 함수 모음.
var newFuncMap = template.FuncMap{
	"dateKO":     DateKO,
	"searchText": searchText,
	"itemCount":  func(w theedit.Week) int { return len(w.Items) },
}

// searchText 는 카드 검색용 색인 문자열(소문자)을 만든다.
func searchText(it theedit.Item) string {
	parts := []string{it.Name, it.Summary}
	for _, l := range it.Links {
		parts = append(parts, l.Label)
	}
	return strings.ToLower(strings.Join(parts, " "))
}

// newPageData 는 out/new/index.html 전체의 입력.
type newPageData struct {
	GeneratedAt string
	Weeks       []theedit.Week
	TotalItems  int
	Empty       bool
	Source      string
}

// WriteNewPage 는 디에디트 [새로나왔] 아카이브 페이지(out/new/index.html)를 만든다.
func WriteNewPage(out string, arch theedit.Archive) error {
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	tpl, err := template.New("new").Funcs(newFuncMap).Parse(newTmpl)
	if err != nil {
		return fmt.Errorf("[새로나왔] 템플릿 파싱 실패: %w", err)
	}
	data := newPageData{
		GeneratedAt: arch.GeneratedAt,
		Weeks:       arch.Weeks,
		TotalItems:  arch.TotalItems,
		Empty:       len(arch.Weeks) == 0,
		Source:      theedit.APIURL,
	}
	if data.TotalItems == 0 {
		for _, w := range arch.Weeks {
			data.TotalItems += len(w.Items)
		}
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.Execute(f, data)
}

// newTmpl 은 index.html 과 같은 순백 배경·29CM 에디토리얼 톤의 단일 파일 페이지다.
const newTmpl = `<!DOCTYPE html>
<html lang="ko">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>THE EDIT · [새로나왔] 신상 리스트 | NEW ARRIVALS</title>
<meta name="description" content="디에디트가 매주 고른 신상 리스트 [새로나왔]을 한 페이지에. 주차별 아카이브와 키워드 검색, 구매 링크까지 한눈에.">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Pretendard:wght@400;500;600;700;800&family=Playfair+Display:ital,wght@0,600;1,600&display=swap" rel="stylesheet">
<style>
:root {
  --bg: #ffffff;
  --bg-sub: #f8f9fa;
  --card: #ffffff;
  --line: #e5e5e5;
  --line-bold: #000000;
  --text-main: #111111;
  --text-body: #374151;
  --text-sub: #6b7280;
  --text-muted: #9ca3af;
  --point-green: #006b54;
  --point-green-light: #eef8f5;
  --tag-red: #d92d20;
  --tag-red-bg: #fff0f0;
  --tag-red-border: #ffc9c9;
}

* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  background: var(--bg);
  color: var(--text-main);
  font-family: 'Pretendard', -apple-system, BlinkMacSystemFont, 'Apple SD Gothic Neo', 'Noto Sans KR', system-ui, sans-serif;
  line-height: 1.55;
  -webkit-font-smoothing: antialiased;
}
a { color: inherit; text-decoration: none; }

.wrap { max-width: 1140px; margin: 0 auto; padding: 0 20px 80px; }

/* 사이트 내비게이션 — 두 페이지를 오가는 탭 */
.site-nav {
  display: flex;
  gap: 2px;
  border-bottom: 1px solid var(--line);
  margin-bottom: 30px;
  overflow-x: auto;
  scrollbar-width: none;
}
.site-nav::-webkit-scrollbar { display: none; }
.site-nav a {
  white-space: nowrap;
  padding: 12px 16px;
  font-size: 13.5px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--text-sub);
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: color 0.15s ease, border-color 0.15s ease;
}
.site-nav a:hover { color: var(--text-main); }
.site-nav a.active { color: var(--text-main); border-bottom-color: var(--line-bold); font-weight: 750; }

header { padding: 8px 0 24px; border-bottom: 1px solid var(--line); }

.brand-eyebrow {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  font-size: 11.5px;
  font-weight: 700;
  letter-spacing: 0.22em;
  color: var(--text-main);
  text-transform: uppercase;
  margin-bottom: 14px;
}
.brand-eyebrow .right-badge {
  font-size: 11px;
  letter-spacing: 0.08em;
  color: var(--point-green);
  font-weight: 600;
  background: var(--point-green-light);
  padding: 3px 8px;
  border-radius: 2px;
  white-space: nowrap;
}

h1.main-title {
  font-size: clamp(30px, 5.2vw, 50px);
  font-weight: 850;
  letter-spacing: -0.025em;
  line-height: 1.1;
  margin-bottom: 10px;
}
h1.main-title .serif {
  font-family: 'Playfair Display', 'Didot', serif;
  font-style: italic;
  font-weight: 600;
}

.sub-desc {
  font-size: 15px;
  color: var(--text-sub);
  line-height: 1.5;
  margin-bottom: 16px;
  word-break: keep-all;
}

.header-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px 12px;
  padding-top: 6px;
  font-size: 13px;
  color: var(--text-sub);
}
.asof-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: 1px solid var(--text-main);
  background: #fff;
  color: var(--text-main);
  font-weight: 700;
  padding: 4px 12px;
  border-radius: 2px;
  font-size: 12.5px;
}
.new-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11.5px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--tag-red);
  background: var(--tag-red-bg);
  border: 1px solid var(--tag-red-border);
  border-radius: 2px;
  padding: 3px 9px;
}
.origin-link {
  color: var(--text-sub);
  border-bottom: 1px dotted var(--text-muted);
  font-size: 12.5px;
}
.origin-link:hover { color: var(--text-main); border-bottom-color: var(--text-main); }
.updated-line { margin-top: 10px; font-size: 11.5px; color: var(--text-muted); letter-spacing: 0.01em; }

/* Sticky controls */
.controls-bar {
  position: sticky;
  top: 0;
  z-index: 40;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--line);
  padding: 16px 0 14px;
}
.search-wrap { position: relative; width: 100%; }
.search-icon { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); color: var(--text-muted); pointer-events: none; }
.search-input {
  width: 100%;
  background: var(--bg-sub);
  border: 1px solid var(--line);
  border-radius: 4px;
  color: var(--text-main);
  padding: 12px 16px 12px 40px;
  font-size: 14.5px;
  outline: none;
  transition: all 0.2s ease;
}
.search-input:focus { background: #fff; border-color: var(--text-main); box-shadow: 0 0 0 1px var(--text-main); }

.weektabs {
  display: flex;
  gap: 6px;
  margin-top: 12px;
  overflow-x: auto;
  scrollbar-width: none;
  padding-bottom: 2px;
}
.weektabs::-webkit-scrollbar { display: none; }
.ftab {
  white-space: nowrap;
  border: 1px solid var(--line);
  background: #fff;
  color: var(--text-sub);
  border-radius: 4px;
  padding: 8px 15px;
  font-size: 13.5px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.ftab:hover { color: var(--text-main); border-color: var(--text-main); }
.ftab.active { background: var(--text-main); border-color: var(--text-main); color: #fff; }
.ftab .tab-new {
  font-size: 9.5px;
  font-weight: 800;
  letter-spacing: 0.06em;
  color: var(--tag-red);
  background: var(--tag-red-bg);
  border-radius: 2px;
  padding: 1px 4px;
  margin-left: 2px;
  vertical-align: 1px;
}
.ftab.active .tab-new { background: #fff; }

.count-bar {
  margin: 24px 0 0;
  font-size: 14px;
  color: var(--text-sub);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}
.count-bar .total-highlight { color: var(--point-green); font-weight: 700; font-size: 15px; }

/* 주차 섹션 헤더 */
.week-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin: 40px 0 18px;
  padding-bottom: 10px;
  border-bottom: 3px solid var(--line-bold);
}
.week-head .head-left { display: flex; align-items: baseline; gap: 10px; flex-wrap: wrap; }
.week-head .week-code { font-size: 25px; font-weight: 850; letter-spacing: -0.02em; }
.week-head .week-title { font-size: 14.5px; font-weight: 500; color: var(--text-sub); }
.week-head .week-title a { border-bottom: 1px dotted var(--text-muted); }
.week-head .week-title a:hover { color: var(--text-main); border-bottom-color: var(--text-main); }
.week-head .week-count { font-size: 12.5px; color: var(--text-muted); font-weight: 500; white-space: nowrap; }
.week-new {
  font-size: 10.5px;
  font-weight: 800;
  letter-spacing: 0.05em;
  color: var(--tag-red);
  background: var(--tag-red-bg);
  border: 1px solid var(--tag-red-border);
  border-radius: 2px;
  padding: 2px 7px;
}
.week-intro {
  font-size: 13.5px;
  color: var(--text-body);
  line-height: 1.7;
  background: var(--bg-sub);
  border-left: 3px solid var(--line-bold);
  padding: 12px 16px;
  margin-bottom: 22px;
  word-break: keep-all;
}

/* 아이템 그리드 */
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 26px 22px;
}
.card {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--line);
  border-radius: 4px;
  overflow: hidden;
  background: var(--card);
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}
.card:hover { transform: translateY(-3px); border-color: var(--text-main); box-shadow: 0 10px 24px rgba(0,0,0,0.06); }
.thumb { width: 100%; aspect-ratio: 1 / 1; background: #f4f4f4; overflow: hidden; }
.thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.card-body { display: flex; flex-direction: column; gap: 10px; padding: 16px 16px 18px; flex: 1; }
.pname {
  font-size: 17px;
  font-weight: 750;
  letter-spacing: -0.02em;
  line-height: 1.4;
  word-break: keep-all;
  overflow-wrap: break-word;
}
.summary-desc {
  font-size: 13px;
  color: var(--text-body);
  line-height: 1.65;
  white-space: pre-line;
  word-break: keep-all;
  overflow-wrap: break-word;
}
.links { display: flex; flex-wrap: wrap; gap: 6px; margin-top: auto; padding-top: 6px; }
.link-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid var(--text-main);
  border-radius: 3px;
  background: var(--text-main);
  color: #fff;
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: -0.01em;
  padding: 5px 10px;
  transition: opacity 0.15s ease, background 0.15s ease;
  word-break: break-all;
}
.link-btn:hover { opacity: 0.82; }
.link-btn .arr { font-size: 10.5px; opacity: 0.8; }
.link-btn.ghost {
  background: #fff;
  color: var(--text-sub);
  border-color: var(--line);
}
.link-btn.ghost:hover { color: var(--text-main); border-color: var(--text-main); opacity: 1; }

.empty-state { padding: 70px 0; text-align: center; color: var(--text-muted); font-size: 15px; }
.empty-state a { color: var(--text-main); border-bottom: 1px solid var(--text-main); }

/* Floating top button */
.top-btn {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 44px;
  height: 44px;
  background: #000;
  color: #fff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px rgba(0,0,0,0.18);
  transition: transform 0.2s ease, opacity 0.2s ease;
  z-index: 90;
  opacity: 0;
  visibility: hidden;
}
.top-btn.visible { opacity: 1; visibility: visible; }
.top-btn:hover { transform: translateY(-2px); }

footer {
  margin-top: 64px;
  padding-top: 30px;
  border-top: 2px solid var(--line-bold);
  color: var(--text-sub);
  font-size: 12.5px;
  line-height: 1.7;
}
.footer-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 12px;
}
.footer-brand { font-size: 14px; font-weight: 800; letter-spacing: 0.15em; color: var(--text-main); }
.disclaimer-text { font-size: 11.5px; color: var(--text-muted); margin-top: 6px; }

@media (max-width: 768px) {
  .wrap { padding: 0 16px 60px; }
  header { padding: 4px 0 20px; }
  h1.main-title { font-size: 32px; }
  .grid { grid-template-columns: 1fr; gap: 20px; }
  .week-head { margin-top: 30px; }
  .week-head .week-code { font-size: 21px; }
  .site-nav a { padding: 11px 12px; font-size: 13px; }
}
</style>
</head>
<body>
<div class="wrap">
<header>
  <nav class="site-nav">
    <a href="../index.html">더현대 팝업</a>
    <a href="index.html" class="active" aria-current="page">디에디트 [새로나왔]</a>
  </nav>
  <div class="brand-eyebrow">
    <span>THE EDIT · WEEKLY NEW ARRIVALS</span>
    <span class="right-badge">[새로나왔] ARCHIVE</span>
  </div>
  <h1 class="main-title">NEW <span class="serif">ARRIVALS</span></h1>
  <div class="sub-desc">디에디트가 매주 고른 신상 리스트 [새로나왔] · 이번 주 새 제품부터 지난 주차 아카이브까지 한 페이지에</div>
  {{if not .Empty}}
  {{$latest := index .Weeks 0}}
  <div class="header-meta">
    <span class="asof-badge">ISSUE {{$latest.Week}}</span>
    {{if $latest.New}}<span class="new-pill">● NEW · 이번 주 업데이트</span>{{end}}
    <span>{{dateKO $latest.Date}} 발행 · {{itemCount $latest}}개 신상</span>
    <a class="origin-link" href="{{$latest.Link}}" target="_blank" rel="noopener">원문 보기 ↗</a>
  </div>
  {{end}}
  <div class="updated-line">최근 수집 {{.GeneratedAt}} · 매일 아침 자동 갱신 (매주 월요일 새 주차)</div>
</header>

<div class="controls-bar">
  <div class="search-wrap">
    <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="11" cy="11" r="8"></circle>
      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
    </svg>
    <input class="search-input" id="q" type="search" placeholder="신상 이름 · 설명 · 구매처 검색 (예: 달력, 틴트, 콘칲, 스마트스토어)">
  </div>
  <div class="weektabs" id="weektabs">
    {{range $i, $w := .Weeks}}
    <button class="ftab{{if eq $i 0}} active{{end}}" data-week="{{$w.ID}}">{{$w.Week}}{{if $w.New}} <span class="tab-new">NEW</span>{{end}}</button>
    {{end}}
  </div>
</div>

<div class="count-bar" id="count-bar">
  <div id="count"></div>
  <div style="font-size:12px; color:var(--text-muted)">THE EDIT [새로나왔] 에디토리얼 뷰</div>
</div>

<div id="list">
{{range $w := .Weeks}}
  <section class="week-section" data-week="{{$w.ID}}" data-items="{{itemCount $w}}" data-label="{{$w.Week}}" data-count-label="{{itemCount $w}}개 · {{dateKO $w.Date}}"{{if $w.New}} data-new="1"{{end}}{{if ne $w.ID (index $.Weeks 0).ID}} hidden{{end}}>
    <div class="week-head">
      <div class="head-left">
        <span class="week-code">{{$w.Week}}</span>
        <span class="week-title"><a href="{{$w.Link}}" target="_blank" rel="noopener">{{$w.Title}}</a></span>
        {{if $w.New}}<span class="week-new">NEW</span>{{end}}
      </div>
      <span class="week-count">{{itemCount $w}}개 · {{dateKO $w.Date}}</span>
    </div>
    {{if $w.Intro}}<p class="week-intro">{{$w.Intro}}</p>{{end}}
    <div class="grid">
      {{range $w.Items}}
      <article class="card" data-q="{{searchText .}}">
        {{if .Image}}<div class="thumb"><img src="{{.Image}}" alt="{{.Name}}" loading="lazy" decoding="async"></div>{{end}}
        <div class="card-body">
          <div class="pname">{{.Name}}</div>
          {{if .Summary}}<div class="summary-desc">{{.Summary}}</div>{{end}}
          <div class="links">
            {{if .Links}}
              {{range .Links}}<a class="link-btn" href="{{.URL}}" target="_blank" rel="noopener nofollow">{{.Label}} <span class="arr">↗</span></a>{{end}}
            {{else}}
              <a class="link-btn ghost" href="{{$w.Link}}" target="_blank" rel="noopener">원문에서 보기 <span class="arr">↗</span></a>
            {{end}}
          </div>
        </div>
      </article>
      {{end}}
    </div>
  </section>
{{end}}
</div>

{{if .Empty}}
<div class="empty-state">
  아직 수집된 [새로나왔] 글이 없어요.<br>
  <a href="https://the-edit.co.kr/" target="_blank" rel="noopener">디에디트에서 직접 보기 ↗</a>
</div>
{{end}}
<div class="empty-state" id="empty" hidden>조건에 맞는 신상이 없어요. 검색어나 주차 탭을 바꿔 보세요.</div>

<footer>
  <div class="footer-top">
    <div class="footer-brand">THE EDIT × 29MAGAZINE STYLE</div>
    <div>자동 수집: 매일 아침 KST (GitHub Actions) · 새 주차는 매주 월요일</div>
  </div>
  <div>출처: 디에디트 [새로나왔] 주간 신상 리스트 (the-edit.co.kr) — 본문의 제품명·설명·구매 링크를 그대로 옮겨 담았습니다.</div>
  <div class="disclaimer-text">ⓒ 디에디트. 이미지·설명의 저작권은 디에디트와 해당 브랜드에 있으며, 이 페이지는 개인 열람용 큐레이션 뷰입니다.</div>
</footer>
</div>

<button class="top-btn" id="topBtn" aria-label="맨 위로 가기">
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
    <polyline points="18 15 12 9 6 15"></polyline>
  </svg>
</button>

<script>
const sections = Array.from(document.querySelectorAll(".week-section"));
const tabs = Array.from(document.querySelectorAll(".ftab"));
const q = document.getElementById("q");
const totalItems = {{.TotalItems}};
const weekCount = {{len .Weeks}};
let curWeek = sections.length ? sections[0].dataset.week : null;

function apply() {
  const term = q.value.trim().toLowerCase();
  let shownTotal = 0;
  let pickLabel = "이번 주";
  let pickCount = 0;

  sections.forEach((sec, idx) => {
    const cards = Array.from(sec.querySelectorAll(".card"));
    let shown = 0;
    cards.forEach(c => {
      const hit = !term || (c.dataset.q || "").includes(term);
      c.hidden = !hit;
      if (hit) shown++;
    });

    const visible = term ? shown > 0 : sec.dataset.week === curWeek;
    sec.hidden = !visible;
    const cnt = sec.querySelector(".week-count");
    if (cnt) cnt.textContent = term ? shown + "개 검색됨" : (sec.dataset.countLabel || sec.dataset.items + "개");
    if (visible) {
      shownTotal += shown;
      if (!term) {
        pickLabel = idx === 0 ? "이번 주" : (sec.dataset.label || "선택한 주차");
        pickCount = shown;
      }
    }
  });

  const countEl = document.getElementById("count");
  if (term) {
    countEl.innerHTML = '검색 결과 <span class="total-highlight">' + shownTotal + '</span>개 (전체 ' + totalItems + '개 중)';
  } else {
    countEl.innerHTML = pickLabel + ' <span class="total-highlight">' + pickCount + '</span>개 신상 · 아카이브 ' + weekCount + '주 / 총 ' + totalItems + '개';
  }
  document.getElementById("empty").hidden = shownTotal > 0;
}

q.addEventListener("input", apply);

tabs.forEach(t => t.addEventListener("click", () => {
  tabs.forEach(x => x.classList.remove("active"));
  t.classList.add("active");
  curWeek = t.dataset.week;
  if (q.value.trim()) q.value = "";
  apply();
}));

const topBtn = document.getElementById("topBtn");
window.addEventListener("scroll", () => {
  topBtn.classList.toggle("visible", window.scrollY > 300);
});
topBtn.addEventListener("click", () => window.scrollTo({ top: 0, behavior: "smooth" }));

apply();
</script>
</body>
</html>
`
