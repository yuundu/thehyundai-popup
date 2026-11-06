package site

// indexTmpl 은 단일 파일로 배포 가능한 index.html 템플릿이다.
// 브라우저에서 이 파일에 박힌 목록으로 검색과 필터를 수행한다.
const indexTmpl = `<!DOCTYPE html>
<html lang="ko">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>더현대 서울 팝업 레이더</title>
<meta name="description" content="더현대 서울 팝업을 층별로 모은 목록. 평일 아침 08:49에 다시 수집합니다.">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Hanken+Grotesk:wght@500;600&family=IBM+Plex+Mono:wght@400;500&family=Newsreader:ital,opsz,wght@0,6..72,400;0,6..72,500;1,6..72,400&family=Noto+Sans+KR:wght@400;500;600&family=Noto+Serif+KR:wght@400;500&display=swap" rel="stylesheet">
<style>
{{cssStamp}}
:root {
  color-scheme: light;
  --space-3xs: 2px; --space-2xs: 4px; --space-xs: 8px; --space-sm: 12px; --space-md: 16px; --space-lg: 24px; --space-xl: 40px; --space-2xl: 64px; --space-3xl: 96px; --space-4xl: 144px;
  --text-xs: 0.8125rem; --text-sm: 0.875rem; --text-md: 1rem; --text-lg: 1.25rem; --text-xl: 1.5625rem; --text-2xl: 1.953rem; --text-3xl: 2.441rem;
  --text-display-s: clamp(2.25rem, 4vw + 0.5rem, 3.75rem); --text-display: clamp(2.75rem, 5vw + 1rem, 5.25rem);
  --text-section: clamp(2rem, 2.5vw + 1rem, 3.5rem); --text-title: clamp(1.25rem, 0.75vw + 0.875rem, 1.625rem); --text-lead: clamp(1.125rem, 0.35vw + 1rem, 1.375rem);
  --ease-out: cubic-bezier(0.16, 1, 0.3, 1); --ease-in: cubic-bezier(0.7, 0, 0.84, 0); --ease-in-out: cubic-bezier(0.65, 0, 0.35, 1);
  --dur-micro: 120ms; --dur-short: 220ms; --dur-long: 420ms;
  --color-error: oklch(52% 0.19 27); --color-success: oklch(50% 0.12 150);
  --container: 1440px; --gutter: 64px; --gap-columns: 32px; --section-gap-s: 120px; --section-gap: 144px; --section-gap-l: 176px;
  --z-base: 1; --z-raised: 10; --z-dropdown: 100; --z-sticky: 200; --z-sticky-nav: 300; --z-modal: 400; --z-toast: 500; --z-tooltip: 600;
  --control: 48px;
  --tracking-figure: -0.04em;
  --color-paper: oklch(84% 0.012 245); --color-paper-2: oklch(81% 0.012 245); --color-paper-3: oklch(77% 0.012 245);
  --color-ink: oklch(16% 0.020 245); --color-ink-2: oklch(24% 0.018 245); --color-muted: oklch(38% 0.016 245);
  --color-rule: oklch(73% 0.013 240); --color-rule-2: oklch(32% 0.018 240);
  --color-accent: oklch(45% 0.180 350); --color-accent-2: oklch(37% 0.160 350); --color-accent-ink: oklch(97% 0.006 245); --color-focus: oklch(46% 0.190 350);
  --color-sheet: oklch(93.5% 0.008 245);
  --font-display: "Hanken Grotesk", "Noto Sans KR", sans-serif;
  --font-body: "Newsreader", "Noto Serif KR", serif;
  --font-outlier: "IBM Plex Mono", "Noto Sans KR", monospace;
  --radius-card: 0; --radius-pill: 2px; --rule-hair: 1px; --rule-bar: 2px; --shadow-card: none;
  --tracking-display: -0.02em; --display-weight: 600;
}
*, *::before, *::after { box-sizing: border-box; }
html, body { overflow-x: clip; }
html { scrollbar-color: var(--color-muted) var(--color-paper); scroll-padding-top: var(--space-3xl); }
body { margin: 0; background: var(--color-paper); color: var(--color-ink); font-family: var(--font-body); font-size: var(--text-md); line-height: 1.55; -webkit-font-smoothing: antialiased; caret-color: var(--color-accent); }
::selection { background: color-mix(in oklab, var(--color-accent) 22%, transparent); color: var(--color-ink); }
:focus-visible { outline: 2px solid var(--color-focus); outline-offset: 2px; }
a { color: inherit; text-decoration-color: var(--color-accent); text-decoration-thickness: 1px; text-underline-offset: 0.15em; }
a:hover { text-decoration-thickness: 2px; }
h1, h2, h3 { margin: 0; font-family: var(--font-display); font-weight: var(--display-weight); font-synthesis: none; letter-spacing: var(--tracking-display); line-height: 1.05; overflow-wrap: break-word; min-width: 0; text-wrap: balance; }
h1 { max-width: 20ch; font-size: var(--text-display); } h2 { max-width: 30ch; font-size: var(--text-section); } h3 { max-width: 40ch; font-size: var(--text-title); line-height: 1.2; }
.lead { font-size: var(--text-lead); line-height: 1.45; max-width: 52ch; }
p { margin: 0; text-wrap: pretty; }
p:not([class]) { max-width: 65ch; }
img, video { max-width: 100%; height: auto; display: block; }
[hidden] { display: none !important; }
button, input, select, textarea { font: inherit; color: inherit; }
.container { width: 100%; max-width: var(--container); margin-inline: auto; padding-inline: var(--gutter); }
.num { font-family: var(--font-outlier); font-variant-numeric: tabular-nums; font-weight: 500; }
.visually-hidden { position: absolute; width: 1px; height: 1px; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
.slab { background: var(--color-ink); color: var(--color-sheet); }
.slab .container { display: flex; flex-wrap: wrap; justify-content: space-between; align-items: baseline; gap: var(--space-lg); min-height: var(--space-2xl); }
.brand { display: flex; gap: var(--space-md); align-items: baseline; margin: 0; font-family: var(--font-display); font-weight: 600; font-size: var(--text-lg); letter-spacing: var(--tracking-display); }
.brand a { color: var(--color-sheet); text-decoration: underline; text-decoration-color: var(--color-accent); white-space: nowrap; }
.brand span { font-family: var(--font-body); font-weight: 400; font-size: var(--text-md); letter-spacing: 0; }
.slab time { font-family: var(--font-outlier); font-size: var(--text-sm); font-weight: 500; font-variant-numeric: tabular-nums; }
.slab :focus-visible { outline-color: var(--color-sheet); }
.hero-copy { padding-top: var(--space-xl); padding-bottom: var(--space-xl); }
.lead { margin-top: var(--space-md); color: var(--color-ink-2); }
.find { position: sticky; top: 0; z-index: var(--z-sticky); display: flex; flex-direction: column; gap: var(--space-sm); margin-top: var(--space-lg); padding-block: var(--space-sm); background: var(--color-paper); }
.search { width: 100%; min-height: var(--control); padding: var(--space-sm) var(--space-md); background: var(--color-sheet); border: var(--rule-hair) solid var(--color-rule-2); border-radius: var(--radius-pill); font-family: var(--font-body); font-size: var(--text-md); }
.search::placeholder { color: var(--color-muted); }
.tabs { display: flex; flex-wrap: wrap; gap: var(--space-xs); }
.floor-tab, .kind { min-height: var(--control); padding: var(--space-sm) var(--space-md); border: 0; border-bottom: var(--rule-bar) solid transparent; border-radius: 0; background: transparent; color: var(--color-ink-2); font-family: var(--font-display); font-weight: 500; font-size: var(--text-sm); white-space: nowrap; cursor: pointer; }
.floor-tab:hover, .kind:hover, .search:hover { background: var(--color-paper-3); }
.floor-tab:active, .kind:active { background: var(--color-paper-2); }
.floor-tab.active, .kind.active { color: var(--color-ink); border-bottom-color: var(--color-rule-2); background: transparent; }
button:disabled, input:disabled { color: var(--color-muted); background: var(--color-paper-3); cursor: not-allowed; }
.sum { margin-top: var(--space-sm); color: var(--color-ink-2); }
.sheet { background: var(--color-sheet); border-top: var(--rule-bar) solid var(--color-ink); padding-block: var(--space-2xl); }
#count { color: var(--color-muted); font-size: var(--text-sm); margin-bottom: var(--space-lg); }
.floor { margin-top: var(--space-2xl); }
.floor:first-child { margin-top: 0; }
.floor-head { display: grid; grid-template-columns: auto minmax(0, 1fr); column-gap: var(--space-lg); align-items: center; border-bottom: var(--rule-bar) solid var(--color-rule-2); padding-bottom: var(--space-sm); }
.floor-head.no-code { grid-template-columns: minmax(0, 1fr); }
.code { margin: 0; font-family: var(--font-outlier); font-weight: 500; font-size: var(--text-section); letter-spacing: var(--tracking-figure); line-height: 1; color: var(--color-ink); font-variant-numeric: tabular-nums; }
.zone { margin-top: var(--space-2xs); color: var(--color-muted); font-size: var(--text-sm); }
.entry { display: grid; grid-template-columns: minmax(0, 5fr) minmax(0, 3fr) minmax(0, 3fr) minmax(0, 2fr); column-gap: var(--gap-columns); row-gap: var(--space-2xs); align-items: baseline; padding-block: var(--space-md); border-top: var(--rule-hair) solid var(--color-rule); }
.entry h3 { grid-column: 1; }
.when { grid-column: 2; font-size: var(--text-sm); }
.where { grid-column: 3; color: var(--color-ink-2); }
.price { grid-column: 4; text-align: right; font-size: var(--text-sm); }
.states, .blurb, .wait, .differ, .sources { grid-column: 1 / -1; }
.state { font-size: var(--text-sm); color: var(--color-muted); margin-right: var(--space-sm); }
.state.is-open { color: var(--color-accent-2); }
.blurb { color: var(--color-ink-2); }
.wait, .differ, .sources { font-size: var(--text-sm); color: var(--color-ink-2); }
.sources a { margin-right: var(--space-md); }
.empty { padding-block: var(--space-xl); color: var(--color-ink-2); }
.mast { font-family: var(--font-display); font-weight: 600; font-size: var(--text-section); letter-spacing: var(--tracking-display); line-height: 1.05; margin-bottom: var(--space-md); }
footer { padding-top: var(--space-2xl); padding-bottom: var(--space-xl); }
.fine { margin-top: var(--space-sm); max-width: 65ch; color: var(--color-ink-2); font-size: var(--text-sm); }
@media (max-width: 60rem) {
  :root { --gutter: 24px; }
  .entry { grid-template-columns: minmax(0, 1fr); }
  .entry h3, .when, .where, .price { grid-column: 1; }
  .price { text-align: start; }
}
@media (max-width: 40rem) { :root { --gutter: 16px; --gap-columns: 16px; } }
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { animation-duration: 0.01ms !important; animation-iteration-count: 1 !important; transition-duration: 0.01ms !important; scroll-behavior: auto !important; }
}
</style>
</head>
<body>
<header class="hero">
  <nav class="slab">
    <div class="container">
      <p class="brand"><a href="#list">팝업 레이더</a><span>더현대 서울</span></p>
      <time>{{.AsOfLabel}}</time>
    </div>
  </nav>
  <div class="container hero-copy">
    <h1>오늘의 팝업</h1>
    <p class="lead">매일 달라지는 더현대 팝업, 층별로 한눈에.</p>
    <div class="find">
      <label class="visually-hidden" for="q">팝업 검색</label>
      <input class="search" id="q" type="search" placeholder="이름, 메뉴, 장소" autocomplete="off">
      <div class="tabs" id="floortabs">
        <button class="floor-tab active" type="button" data-floor="all">전체</button>
        <button class="floor-tab" type="button" data-floor="B1">B1</button>
        <button class="floor-tab" type="button" data-floor="B2">B2</button>
        <button class="floor-tab" type="button" data-floor="1F">1F</button>
        <button class="floor-tab" type="button" data-floor="5F">5F</button>
      </div>
      <div class="tabs" id="chips">
        <button class="kind active" type="button" data-cat="all">전체</button>
        <button class="kind" type="button" data-cat="food">식품</button>
        <button class="kind" type="button" data-cat="beauty">뷰티</button>
        <button class="kind" type="button" data-cat="fashion">패션</button>
        <button class="kind" type="button" data-cat="season">시즌</button>
      </div>
    </div>
    <p class="sum">{{.Summary}} 평일 아침 08:49에 다시 수집합니다.</p>
  </div>
</header>

<section class="sheet">
  <div class="container">
    <p id="count"></p>
    <div id="list"></div>
    <p class="empty" id="empty" hidden>조건에 맞는 팝업이 없어요. 검색어나 필터를 바꿔 보세요.</p>
  </div>
</section>

<footer>
  <div class="container">
    <p class="mast">팝업 레이더</p>
    <p>데이터 기준은 더현대 서울 공식, 팝업 플랫폼, 언론 보도입니다.</p>
    <p class="fine">출처마다 안내가 다르면 두 정보를 모두 적습니다. 2026년 일정만 남겼고, 현장과 다를 수 있으니 방문 전 공식 안내를 확인하세요.</p>
  </div>
</footer>

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

const badgeText = {new:"오늘 오픈", ending:"종료 임박", closed:"오늘 마지막", deal:"한정 특가", conflict:"정보 상이"};
const floorNames = {B1:["테이스티&이슈","식품"], B2:["영&트렌디","컬처"], "1F":["럭셔리&뷰티",""], "4F":["컬처&레저",""], "5F":["사운즈포레스트","스튜디오"]};

let curFloor = "all", curCat = "all", curQ = "";

const esc = s => String(s??"").replace(/[&<>"']/g, c => ({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]));

function sourceHTML(s){
  if(s.url) return '<a href="'+esc(s.url)+'" target="_blank" rel="noopener">'+esc(s.name)+'</a>';
  return '<span>'+esc(s.name)+'</span>';
}

function el(name, inner){ return "<" + name + ">" + inner + "</" + name + ">"; }

function cardHTML(p){
  const states = p.statuses.map(s => '<span class="state'+(s==="new"?" is-open":"")+'" role="status">'+esc(badgeText[s]||s)+'</span>').join("");
  let h = '<article class="entry" data-floor="'+esc(p.floor)+'" data-cat="'+esc(p.category)+'">';
  h += el("h3", esc(p.name));
  h += '<p class="when num">'+esc(p.period)+'</p>';
  if(p.venue) h += '<p class="where">'+esc(p.venue)+'</p>';
  if(p.price) h += '<p class="price num">'+esc(p.price)+'</p>';
  if(states) h += '<p class="states">'+states+'</p>';
  if(p.summary) h += '<p class="blurb">'+esc(p.summary)+'</p>';
  if(p.waiting) h += '<p class="wait">'+esc(p.waiting)+'</p>';
  if(p.venueAlt || p.floorAlt || p.endAlt){
    const bits = [];
    if(p.venueAlt) bits.push("장소 "+p.venueAlt);
    if(p.floorAlt) bits.push("층 "+p.floorAlt);
    if(p.endAlt) bits.push("종료일 "+p.endAlt);
    h += '<p class="differ">'+esc(bits.join(". "))+'. 출처마다 안내가 다릅니다.</p>';
  }
  if(p.sources.length){
    h += '<p class="sources">출처 '+p.sources.map(sourceHTML).join(" ")+'</p>';
  }
  return h + '</article>';
}

function render(){
  const list = document.getElementById("list");
  const shown = DATA.filter(p =>
    (curFloor==="all" || p.floor===curFloor) &&
    (curCat==="all" || p.category===curCat) &&
    (!curQ || (p.name+" "+p.brand+" "+p.venue+" "+p.summary+" "+p.price+" "+p.floor).toLowerCase().includes(curQ))
  );
  const byFloor = {};
  shown.forEach(p => (byFloor[p.floor] ||= []).push(p));
  const order = ["B2","B1","1F","2F","4F","5F","?"];
  let html = "";
  order.forEach(f => {
    if(!byFloor[f]) return;
    const pair = floorNames[f] || ["층 미확인", ""];
    const zone = (pair[1] ? pair[1] + ", " : "") + byFloor[f].length + "건";
    const code = f === "?" ? "" : '<p class="code">'+esc(f)+'</p>';
    html += '<div class="floor" id="floor-'+esc(f)+'">';
    html += '<div class="floor-head'+(f === "?" ? " no-code" : "")+'">'+code+'<div>'+el("h2", esc(pair[0]))+'<p class="zone">'+esc(zone)+'</p></div></div>';
    html += byFloor[f].map(cardHTML).join("");
    html += '</div>';
  });
  list.innerHTML = html;
  document.getElementById("empty").hidden = shown.length > 0;
  document.getElementById("count").textContent = "표시 "+shown.length+"건, 전체 "+DATA.length+"건";
}

document.getElementById("q").addEventListener("input", e => { curQ = e.target.value.trim().toLowerCase(); render(); });
document.querySelectorAll(".floor-tab").forEach(b => b.addEventListener("click", () => {
  document.querySelectorAll(".floor-tab").forEach(x=>x.classList.remove("active"));
  b.classList.add("active"); curFloor = b.dataset.floor; render();
}));
document.querySelectorAll(".kind").forEach(b => b.addEventListener("click", () => {
  document.querySelectorAll(".kind").forEach(x=>x.classList.remove("active"));
  b.classList.add("active"); curCat = b.dataset.cat; render();
}));
render();
</script>
</body>
</html>
`
