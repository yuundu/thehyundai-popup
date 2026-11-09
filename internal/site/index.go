package site

// indexTmpl 은 단일 파일로 배포 가능한 index.html 템플릿이다.
// 더현대 서울 공식 웹사이트와 29CM 매거진의 감각적인 에디토리얼 레이아웃을 반영했다.
const indexTmpl = `<!DOCTYPE html>
<html lang="ko">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>THE HYUNDAI SEOUL · POP-UP RADAR | 매거진 에디션</title>
<meta name="description" content="더현대 서울 팝업스토어를 층별로 한눈에. 매일 아침 07:30 KST 자동 업데이트되는 에디토리얼 팝업 매거진.">
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
  --tag-amber: #b45309;
  --tag-amber-bg: #fffbeb;
  --tag-amber-border: #fde68a;
  --tag-deal: #047857;
  --tag-deal-bg: #ecfdf5;
  --tag-deal-border: #a7f3d0;
  --tag-warn: #c2410c;
  --tag-warn-bg: #fff7ed;
  --tag-warn-border: #fed7aa;
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

.wrap {
  max-width: 1140px;
  margin: 0 auto;
  padding: 0 20px 80px;
}

/* 29CM & THE HYUNDAI Editorial Header */
header {
  padding: 44px 0 24px;
  border-bottom: 1px solid var(--line);
}

.brand-eyebrow {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
}

h1.main-title {
  font-size: clamp(30px, 5.2vw, 50px);
  font-weight: 850;
  letter-spacing: -0.025em;
  line-height: 1.1;
  color: var(--text-main);
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
  gap: 12px;
  padding-top: 8px;
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
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 2px;
  font-size: 12px;
}

.summary-text b {
  color: var(--point-green);
  font-weight: 700;
}

/* Sticky Controls */
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

.search-wrap {
  position: relative;
  width: 100%;
}

.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
}

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

.search-input:focus {
  background: #ffffff;
  border-color: var(--text-main);
  box-shadow: 0 0 0 1px var(--text-main);
}

/* 29CM Tab Buttons */
.floortabs {
  display: flex;
  gap: 6px;
  margin-top: 12px;
  overflow-x: auto;
  scrollbar-width: none;
  padding-bottom: 2px;
}
.floortabs::-webkit-scrollbar { display: none; }

.ftab {
  white-space: nowrap;
  border: 1px solid var(--line);
  background: #ffffff;
  color: var(--text-sub);
  border-radius: 4px;
  padding: 8px 16px;
  font-size: 13.5px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.ftab:hover {
  color: var(--text-main);
  border-color: var(--text-main);
}

.ftab.active {
  background: var(--text-main);
  border-color: var(--text-main);
  color: #ffffff;
}

/* Category Chips */
.chips {
  display: flex;
  gap: 6px;
  margin-top: 10px;
  flex-wrap: wrap;
}

.chip {
  border: 1px solid var(--line);
  background: #ffffff;
  color: var(--text-sub);
  border-radius: 4px;
  padding: 6px 13px;
  font-size: 12.5px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

.chip:hover {
  border-color: var(--text-sub);
  color: var(--text-main);
}

.chip.active {
  border-color: var(--point-green);
  color: var(--point-green);
  background: var(--point-green-light);
  font-weight: 600;
}

/* Count bar */
.count-bar {
  margin: 24px 0 16px;
  font-size: 14px;
  color: var(--text-sub);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.count-bar .total-highlight {
  color: var(--point-green);
  font-weight: 700;
  font-size: 15px;
}

/* 29CM Editorial Section Header */
.floor-section-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin: 44px 0 20px;
  padding-bottom: 10px;
  border-bottom: 3px solid var(--line-bold);
}

.floor-section-head .head-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.floor-section-head .floor-code {
  font-size: 26px;
  font-weight: 850;
  letter-spacing: -0.02em;
  color: var(--text-main);
}

.floor-section-head .floor-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-sub);
}

.floor-section-head .floor-count {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}

/* Cards Grid */
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(330px, 1fr));
  gap: 22px;
}

/* 29CM Magazine Editorial Card */
.card {
  background: var(--card);
  border: 1px solid var(--line);
  border-radius: 4px;
  padding: 22px 20px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
  position: relative;
  overflow: hidden;
}

.card:hover {
  transform: translateY(-3px);
  border-color: var(--text-main);
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.06);
}

.card.is-new {
  border-top: 3px solid var(--tag-red);
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.badges {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  align-items: center;
}

.bdg {
  font-size: 11px;
  font-weight: 600;
  border-radius: 2px;
  padding: 3px 7px;
  letter-spacing: -0.01em;
  border: 1px solid transparent;
}

.bdg.floor {
  background: #f4f4f4;
  color: var(--text-main);
  border-color: #e0e0e0;
  font-weight: 700;
}

.bdg.new {
  background: var(--tag-red-bg);
  color: var(--tag-red);
  border-color: var(--tag-red-border);
}

.bdg.ending {
  background: var(--tag-amber-bg);
  color: var(--tag-amber);
  border-color: var(--tag-amber-border);
}

.bdg.deal {
  background: var(--tag-deal-bg);
  color: var(--tag-deal);
  border-color: var(--tag-deal-border);
}

.bdg.warn {
  background: var(--tag-warn-bg);
  color: var(--tag-warn);
  border-color: var(--tag-warn-border);
}

.bdg.closed {
  background: #f3f4f6;
  color: var(--text-muted);
}

.category-label {
  font-size: 11.5px;
  color: var(--text-muted);
  font-weight: 500;
}

.pname {
  font-size: 19px;
  font-weight: 750;
  letter-spacing: -0.02em;
  color: var(--text-main);
  line-height: 1.35;
  word-break: keep-all;
  overflow-wrap: break-word;
}

.period {
  font-size: 13.5px;
  color: var(--point-green);
  font-weight: 600;
  letter-spacing: -0.01em;
}

.summary-desc {
  font-size: 13.5px;
  color: var(--text-body);
  line-height: 1.6;
  word-break: keep-all;
}

.meta-box {
  display: flex;
  flex-direction: column;
  gap: 5px;
  font-size: 13px;
  color: var(--text-sub);
  background: var(--bg-sub);
  padding: 10px 12px;
  border-radius: 3px;
}

.meta-row {
  display: flex;
  gap: 8px;
  align-items: baseline;
}

.meta-row .label {
  font-weight: 600;
  color: var(--text-main);
  min-width: 32px;
}

.meta-row .price-val {
  color: var(--text-main);
  font-weight: 700;
}

.waiting-box {
  font-size: 12.5px;
  color: #1e293b;
  background: #f8fafc;
  border-left: 3px solid var(--point-green);
  border-radius: 0 3px 3px 0;
  padding: 8px 12px;
  line-height: 1.45;
  word-break: keep-all;
}

.alt-box {
  font-size: 12px;
  color: #7c2d12;
  background: var(--tag-warn-bg);
  border-left: 3px solid var(--tag-warn);
  border-radius: 0 3px 3px 0;
  padding: 8px 12px;
  line-height: 1.45;
  word-break: keep-all;
}

.card-foot {
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11.5px;
  color: var(--text-muted);
}

.card-foot a {
  color: var(--text-sub);
  border-bottom: 1px dotted var(--text-muted);
}

.card-foot a:hover {
  color: var(--text-main);
  border-bottom-color: var(--text-main);
}

.empty-state {
  padding: 60px 0;
  text-align: center;
  color: var(--text-muted);
  font-size: 15px;
}

/* Floating Top Button */
.top-btn {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 44px;
  height: 44px;
  background: #000000;
  color: #ffffff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.18);
  transition: transform 0.2s ease, opacity 0.2s ease;
  z-index: 90;
  opacity: 0;
  visibility: hidden;
}

.top-btn.visible {
  opacity: 1;
  visibility: visible;
}

.top-btn:hover {
  transform: translateY(-2px);
}

/* Editorial Footer */
footer {
  margin-top: 60px;
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

.footer-brand {
  font-size: 14px;
  font-weight: 800;
  letter-spacing: 0.15em;
  color: var(--text-main);
}

.disclaimer-text {
  font-size: 11.5px;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .wrap { padding: 0 16px 60px; }
  header { padding: 28px 0 20px; }
  h1.main-title { font-size: 32px; }
  .grid { grid-template-columns: 1fr; }
  .floor-section-head { margin-top: 32px; }
}
</style>
</head>
<body>
<div class="wrap">
<header>
  <div class="brand-eyebrow">
    <span>THE HYUNDAI SEOUL</span>
    <span class="right-badge">POP-UP RADAR · 2026</span>
  </div>
  <h1 class="main-title">POP-UP <span class="serif">RADAR</span></h1>
  <div class="sub-desc">더현대 서울 팝업스토어 & 브랜드 큐레이션 매거진 · 매일 아침 달라지는 행사 한눈에 보기</div>
  <div class="header-meta">
    <div class="asof-badge">ISSUE {{.AsOfLabel}}</div>
    <div class="summary-text">🔎 <b>{{.Summary}}</b> · 매일 평일 07:30 자동 수집(KST)</div>
  </div>
</header>

<div class="controls-bar">
  <div class="search-wrap">
    <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="11" cy="11" r="8"></circle>
      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
    </svg>
    <input class="search-input" id="q" type="search" placeholder="팝업명 · 브랜드 · 메뉴 · 층별 검색 (예: 닭강정, B1, 프라다, 타코야키)">
  </div>
  
  <div class="floortabs" id="floortabs">
    <button class="ftab active" data-floor="all">ALL (전체보기)</button>
    <button class="ftab" data-floor="B1">B1 테이스티 & 이슈</button>
    <button class="ftab" data-floor="B2">B2 영 & 트렌디</button>
    <button class="ftab" data-floor="1F">1F 럭셔리 & 뷰티</button>
    <button class="ftab" data-floor="4F">4F 라이프스타일</button>
    <button class="ftab" data-floor="5F">5F 사운즈포레스트</button>
  </div>
  
  <div class="chips" id="chips">
    <button class="chip active" data-cat="all">#전체</button>
    <button class="chip" data-cat="food">#🍽️ 식품/디저트</button>
    <button class="chip" data-cat="beauty">#💄 뷰티</button>
    <button class="chip" data-cat="fashion">#🛍️ 패션/라이프</button>
    <button class="chip" data-cat="season">#🎄 크리스마스/시즌</button>
  </div>
</div>

<div class="count-bar" id="count-bar">
  <div id="count"></div>
  <div style="font-size:12px; color:var(--text-muted)">매거진 에디토리얼 뷰</div>
</div>

<div id="list"></div>
<div class="empty-state" id="empty" hidden>조건에 맞는 팝업스토어가 없어요. 검색어나 층별 필터를 바꿔 보세요.</div>

<footer>
  <div class="footer-top">
    <div class="footer-brand">THE HYUNDAI SEOUL × 29MAGAZINE STYLE</div>
    <div>자동 수집 스케줄: 매일 평일 07:30 KST (GitHub Actions)</div>
  </div>
  <div>데이터 출처: 현대백화점 공식 공지(이슈 다이어리), 팝플리, 팝놀, 보도자료, 네이버 예약 전수 비교 검증</div>
  <div class="disclaimer-text">
    ⚠️ 출처별 안내가 상이한 경우(예: 올기 B1/B2, 피크퍼포먼스 종료일) 두 정보를 모두 투명하게 병기합니다. 연도(2026) 검증을 통과한 정보만 제공되며, 현장 사정에 따라 변동될 수 있습니다.
  </div>
</footer>
</div>

<button class="top-btn" id="topBtn" aria-label="맨 위로 가기">
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
    <polyline points="18 15 12 9 6 15"></polyline>
  </svg>
</button>

<script>
const DATA = [
{{- range .Floors}}
  {
    floor: "{{floorLabel .Floor}}",
    name: {{.Name}},
    brand: {{.Brand}},
    period: "{{periodKO .}}",
    venue: {{.Venue}},
    venueAlt: {{.VenueAlt}},
    floorAlt: {{.FloorAlt}},
    startAlt: {{.StartDateAlt}},
    endAlt: {{.EndDateAlt}},
    price: {{.Price}},
    summary: {{.Summary}},
    waiting: {{.Waiting}},
    category: {{.Category}},
    chip: "{{chip .Category}}",
    statuses: [{{range .Status}}"{{.}}",{{end}}],
    confidence: {{.Confidence}},
    sources: [
      {{- range .Sources}}
      {name: {{.Name}}, url: {{.URL}}, claim: {{.Claim}}},
      {{- end}}
    ]
  },
{{- end}}
];

const badgeText = {
  new: "🔥 오늘 오픈!",
  ending: "⏳ 종료 임박",
  closed: "⏰ 오늘 마지막",
  deal: "🛒 한정 특가",
  conflict: "⚠️ 정보 상이"
};
const badgeClass = { new: "new", ending: "ending", deal: "deal", conflict: "warn", closed: "closed" };
const floorNames = {
  B1: "테이스티 & 이슈 · 식품관 F&B",
  B2: "영 & 트렌디 · 스트리트 & 컬처",
  "1F": "럭셔리 & 뷰티 · 명품관",
  "2F": "모던 & 트렌디",
  "4F": "라이프스타일 & 레저",
  "5F": "사운즈 포레스트 & 스튜디오",
  "?": "기타 · 층 정보 확인 중"
};

let curFloor = "all", curCat = "all", curQ = "";

const esc = s => String(s ?? "").replace(/[&<>"']/g, c => ({
  "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;"
}[c]));

function cardHTML(p) {
  const isNew = p.statuses.includes("new");
  const floorBadgeText = p.floor === "?" ? "#기타" : "#" + p.floor;
  const badges = ['<span class="bdg floor">' + esc(floorBadgeText) + '</span>']
    .concat(p.statuses.map(s => '<span class="bdg ' + (badgeClass[s] || "") + '">' + (badgeText[s] || esc(s)) + '</span>'));
  
  let h = '<div class="card ' + (isNew ? 'is-new' : '') + '" data-floor="' + esc(p.floor) + '" data-cat="' + esc(p.category) + '">';
  
  h += '<div class="card-top">';
  h += '  <div class="badges">' + badges.join("") + '</div>';
  if (p.chip) h += '  <span class="category-label">' + esc(p.chip) + '</span>';
  h += '</div>';

  h += '<div class="pname">' + esc(p.name) + '</div>';
  h += '<div class="period">' + esc(p.period) + '</div>';

  if (p.summary && p.summary.trim()) {
    h += '<div class="summary-desc">' + esc(p.summary.trim()) + '</div>';
  }

  const metaRows = [];
  if (p.venue && p.venue.trim()) metaRows.push('<div class="meta-row"><span class="label">위치</span><span>' + esc(p.venue.trim()) + '</span></div>');
  if (p.price && p.price.trim()) metaRows.push('<div class="meta-row"><span class="label">가격</span><span class="price-val">' + esc(p.price.trim()) + '</span></div>');

  if (metaRows.length) {
    h += '<div class="meta-box">' + metaRows.join("") + '</div>';
  }

  if (p.waiting && p.waiting.trim()) {
    h += '<div class="waiting-box">⏱️ <b>웨이팅/이용:</b> ' + esc(p.waiting.trim()) + '</div>';
  }

  const bits = [];
  if (p.venueAlt && p.venueAlt.trim()) bits.push("장소: " + p.venueAlt.trim());
  if (p.floorAlt && p.floorAlt.trim()) bits.push("층: " + p.floorAlt.trim());
  if (p.endAlt && p.endAlt.trim()) bits.push("종료일: " + p.endAlt.trim());
  if (bits.length > 0) {
    h += '<div class="alt-box">⚠️ <b>출처별 상이:</b> ' + esc(bits.join(" · ")) + '</div>';
  }

  if (p.sources && p.sources.length) {
    h += '<div class="card-foot"><span>출처</span><span>' + p.sources.map(s => s.url
      ? '<a href="' + esc(s.url) + '" target="_blank" rel="noopener">' + esc(s.name) + '</a>'
      : esc(s.name)).join(" · ") + '</span></div>';
  }

  return h + '</div>';
}

function render() {
  const list = document.getElementById("list");
  const shown = DATA.filter(p =>
    (curFloor === "all" || p.floor === curFloor) &&
    (curCat === "all" || p.category === curCat) &&
    (!curQ || (p.name + " " + p.brand + " " + p.venue + " " + p.summary + " " + p.price + " " + p.floor).toLowerCase().includes(curQ))
  );

  const byFloor = {};
  shown.forEach(p => (byFloor[p.floor] ||= []).push(p));

  const order = ["B1", "B2", "1F", "2F", "4F", "5F", "?"];
  let html = "";

  order.forEach(f => {
    if (!byFloor[f] || byFloor[f].length === 0) return;
    const isUnknown = f === "?";
    html += '<div class="floor-section-head">';
    html += '  <div class="head-left">';
    html += '    <span class="floor-code">' + esc(isUnknown ? "ETC" : f) + '</span>';
    html += '    <span class="floor-name">' + esc(floorNames[f] || "층 정보") + '</span>';
    html += '  </div>';
    html += '  <span class="floor-count">' + byFloor[f].length + '개 팝업</span>';
    html += '</div>';
    html += '<div class="grid">' + byFloor[f].map(cardHTML).join("") + '</div>';
  });

  list.innerHTML = html;
  document.getElementById("empty").hidden = shown.length > 0;
  document.getElementById("count").innerHTML = '총 <span class="total-highlight">' + shown.length + '</span>개의 팝업스토어 (전체 ' + DATA.length + '개 중)';
}

document.getElementById("q").addEventListener("input", e => {
  curQ = e.target.value.trim().toLowerCase();
  render();
});

document.querySelectorAll(".ftab").forEach(b => b.addEventListener("click", () => {
  document.querySelectorAll(".ftab").forEach(x => x.classList.remove("active"));
  b.classList.add("active");
  curFloor = b.dataset.floor;
  render();
}));

document.querySelectorAll(".chip").forEach(b => b.addEventListener("click", () => {
  document.querySelectorAll(".chip").forEach(x => x.classList.remove("active"));
  b.classList.add("active");
  curCat = b.dataset.cat;
  render();
}));

// Floating Top Button
const topBtn = document.getElementById("topBtn");
window.addEventListener("scroll", () => {
  if (window.scrollY > 300) {
    topBtn.classList.add("visible");
  } else {
    topBtn.classList.remove("visible");
  }
});
topBtn.addEventListener("click", () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
});

render();
</script>
</body>
</html>
`
