package site

// indexTmpl 은 단일 파일로 배포 가능한 index.html 템플릿이다.
// 브라우저에서 popups.json 을 추가로 받아 검색/필터를 수행한다.
const indexTmpl = `<!DOCTYPE html>
<html lang="ko">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>THE HYUNDAI SEOUL POP-UP RADAR · 매일 달라지는 더현대 팝업 층별 레이더</title>
<meta name="description" content="더현대 서울 팝업을 층별로 한눈에. 매일 아침 8시 49분 자동 업데이트.">
<style>
:root{
  --bg:#0f1115; --card:#171a21; --card2:#1d212b; --line:#2a2f3c;
  --text:#f2f4f8; --sub:#9aa3b2; --gold:#d8b46a; --gold2:#b8944e;
  --new:#ff5d5d; --end:#ffb020; --deal:#3ddc97; --warn:#ff9f43; --info:#5aa2ff;
}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--text);font-family:'Pretendard','Apple SD Gothic Neo','Noto Sans KR',system-ui,-apple-system,sans-serif;line-height:1.55}
a{color:var(--gold);text-decoration:none}
.wrap{max-width:1080px;margin:0 auto;padding:0 16px 60px}
header{padding:34px 0 10px;text-align:center}
.eyebrow{color:var(--gold);font-size:12px;letter-spacing:.35em;font-weight:700}
h1{font-size:clamp(26px,5vw,44px);font-weight:800;letter-spacing:.02em;margin:8px 0 4px}
h1 .thin{font-weight:300;color:var(--sub)}
.tagline{color:var(--sub);font-size:14px}
.asof{display:inline-block;margin-top:12px;font-size:13px;color:var(--gold2);border:1px solid var(--gold2);border-radius:999px;padding:4px 14px}
.summary{margin-top:10px;font-size:13px;color:var(--sub)}
.summary b{color:var(--new)}
.controls{position:sticky;top:0;z-index:10;background:linear-gradient(var(--bg) 82%,transparent);padding:14px 0 10px}
.search{width:100%;background:var(--card);border:1px solid var(--line);border-radius:12px;color:var(--text);padding:12px 16px;font-size:15px;outline:none}
.search:focus{border-color:var(--gold2)}
.floortabs{display:flex;gap:8px;margin-top:12px;overflow-x:auto;scrollbar-width:none;padding-bottom:2px}
.floortabs::-webkit-scrollbar{display:none}
.ftab{white-space:nowrap;border:1px solid var(--line);background:var(--card);color:var(--sub);border-radius:999px;padding:8px 16px;font-size:14px;font-weight:600;cursor:pointer;transition:.15s}
.ftab:hover{color:var(--text)}
.ftab.active{background:var(--gold);border-color:var(--gold);color:#16130b}
.chips{display:flex;gap:6px;margin-top:10px;flex-wrap:wrap}
.chip{border:1px solid var(--line);background:transparent;color:var(--sub);border-radius:999px;padding:5px 12px;font-size:12.5px;cursor:pointer}
.chip.active{border-color:var(--gold2);color:var(--gold);background:rgba(216,180,106,.08)}
.count{margin:16px 2px 10px;font-size:13px;color:var(--sub)}
.floorhead{display:flex;align-items:center;gap:10px;margin:26px 2px 12px}
.floorhead .fb{background:var(--gold);color:#16130b;font-weight:800;border-radius:8px;padding:3px 12px;font-size:15px}
.floorhead .fname{color:var(--sub);font-size:13.5px}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:14px}
.card{background:var(--card);border:1px solid var(--line);border-radius:16px;padding:16px 16px 14px;display:flex;flex-direction:column;gap:8px;transition:.15s;position:relative;overflow:hidden}
.card:hover{border-color:var(--gold2);transform:translateY(-2px)}
.card.is-new{border-color:rgba(255,93,93,.55)}
.card.is-new::before{content:"";position:absolute;inset:0 auto 0 0;width:3px;background:var(--new)}
.badges{display:flex;gap:6px;flex-wrap:wrap}
.bdg{font-size:11px;font-weight:700;border-radius:6px;padding:2px 8px;background:var(--card2);color:var(--sub);border:1px solid var(--line)}
.bdg.floor{background:rgba(216,180,106,.14);color:var(--gold);border-color:rgba(216,180,106,.4);font-size:12px}
.bdg.new{background:rgba(255,93,93,.15);color:var(--new);border-color:rgba(255,93,93,.45)}
.bdg.ending{background:rgba(255,176,32,.12);color:var(--end);border-color:rgba(255,176,32,.4)}
.bdg.deal{background:rgba(61,220,151,.12);color:var(--deal);border-color:rgba(61,220,151,.4)}
.bdg.warn{background:rgba(255,159,67,.12);color:var(--warn);border-color:rgba(255,159,67,.45)}
.bdg.closed{background:rgba(154,163,178,.12);color:var(--sub)}
.pname{font-size:17px;font-weight:750;letter-spacing:-.01em}
.period{font-size:13px;color:var(--gold2);font-weight:600}
.meta{font-size:13px;color:var(--sub)}
.meta .price{color:var(--text);font-weight:650}
.waiting{font-size:12.5px;color:var(--info);background:rgba(90,162,255,.08);border:1px dashed rgba(90,162,255,.4);border-radius:8px;padding:6px 9px}
.alt{font-size:12.5px;color:var(--warn);background:rgba(255,159,67,.07);border-radius:8px;padding:6px 9px}
.src{margin-top:auto;padding-top:8px;border-top:1px solid var(--line);display:flex;gap:8px;flex-wrap:wrap;font-size:12px;color:var(--sub)}
.src a{color:var(--sub);border-bottom:1px dotted var(--sub)}
.src a:hover{color:var(--gold)}
.empty{padding:50px 0;text-align:center;color:var(--sub)}
footer{margin-top:44px;padding-top:18px;border-top:1px solid var(--line);color:var(--sub);font-size:12.5px;text-align:center}
.disclaimer{margin-top:6px;font-size:11.5px;color:#6b7382}
@media(max-width:640px){ .grid{grid-template-columns:1fr} header{padding-top:26px} }
</style>
</head>
<body>
<div class="wrap">
<header>
  <div class="eyebrow">THE HYUNDAI SEOUL</div>
  <h1>POP-UP <span class="thin">RADAR</span></h1>
  <div class="tagline">매일 달라지는 더현대 팝업, 층별로 한눈에</div>
  <div class="asof">📅 {{.AsOfLabel}}</div>
  <div class="summary">🔎 <b>{{.Summary}}</b> · 매일 평일 아침 08:49 자동 수집(KST)</div>
</header>

<div class="controls">
  <input class="search" id="q" type="search" placeholder="팝업명 · 브랜드 · 메뉴 · 위치 검색 (예: 닭강정, B1, 프라다)">
  <div class="floortabs" id="floortabs">
    <button class="ftab active" data-floor="all">전체보기</button>
    <button class="ftab" data-floor="B1">B1 테이스티&이슈</button>
    <button class="ftab" data-floor="B2">B2 영&트렌디</button>
    <button class="ftab" data-floor="1F">1F 럭셔리&뷰티</button>
    <button class="ftab" data-floor="5F">5F 사운즈포레스트</button>
  </div>
  <div class="chips" id="chips">
    <button class="chip active" data-cat="all">전체</button>
    <button class="chip" data-cat="food">🍽️ 식품/디저트</button>
    <button class="chip" data-cat="beauty">💄 뷰티</button>
    <button class="chip" data-cat="fashion">🛍️ 패션/라이프</button>
    <button class="chip" data-cat="season">🎄 시즌/크리스마스</button>
  </div>
</div>

<div class="count" id="count"></div>
<div id="list"></div>
<div class="empty" id="empty" hidden>조건에 맞는 팝업이 없어요. 검색어나 필터를 바꿔 보세요.</div>

<footer>
  데이터 기준: 더현대 서울 공식·팝업 플랫폼·언론 보도 수집 결과 · 자동 업데이트 매일 평일 08:49 KST<br>
  <span class="disclaimer">⚠️ 출처마다 안내가 다른 경우 두 정보를 모두 표기합니다. 연도(2026) 검증을 마친 정보만 수집하며, 현장 상황과 다를 수 있으니 방문 전 공식 안내를 확인하세요.</span>
</footer>
</div>

<script>
const DATA = [
{{- range .Floors}}
  {
    floor: "{{floorLabel .Floor}}",
    name: {{.Name | printf "%q"}},
    brand: {{.Brand | printf "%q"}},
    period: "{{periodKO .}}",
    venue: {{.Venue | printf "%q"}},
    venueAlt: {{.VenueAlt | printf "%q"}},
    floorAlt: {{.FloorAlt | printf "%q"}},
    startAlt: {{.StartDateAlt | printf "%q"}},
    endAlt: {{.EndDateAlt | printf "%q"}},
    price: {{.Price | printf "%q"}},
    summary: {{.Summary | printf "%q"}},
    waiting: {{.Waiting | printf "%q"}},
    category: {{.Category | printf "%q"}},
    chip: "{{chip .Category}}",
    statuses: [{{range .Status}}"{{.}}",{{end}}],
    confidence: {{.Confidence}},
    sources: [
      {{- range .Sources}}
      {name: {{.Name | printf "%q"}}, url: {{.URL | printf "%q"}}, claim: {{.Claim | printf "%q"}}},
      {{- end}}
    ]
  },
{{- end}}
];

const badgeText = {new:"🔥 오늘 오픈!", ending:"⏳ 종료 임박", closed:"⏰ 오늘 마지막", deal:"🛒 한정 특가", conflict:"⚠️ 정보 상이 확인"};
const badgeClass = {new:"new", ending:"ending", deal:"deal", conflict:"warn", closed:"closed"};
const floorNames = {B1:"테이스티&이슈 · F&B", B2:"영&트렌디 · 컬처", "1F":"럭셔리&뷰티", "4F":"컬처&레저", "5F":"사운즈포레스트&스튜디오"};

let curFloor = "all", curCat = "all", curQ = "";

const esc = s => String(s??"").replace(/[&<>"']/g, c => ({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]));

function cardHTML(p){
  const badges = ['<span class="bdg floor">'+esc(p.floor)+'</span>']
    .concat(p.statuses.map(s => '<span class="bdg '+(badgeClass[s]||"")+'">'+(badgeText[s]||esc(s))+'</span>'));
  let h = '<div class="card '+(p.statuses.includes("new")?"is-new":"")+'" data-floor="'+esc(p.floor)+'" data-cat="'+esc(p.category)+'" data-text="'+esc([p.name,p.brand,p.venue,p.summary,p.price,p.floor].join(" "))+'">';
  h += '<div class="badges">'+badges.join("")+'</div>';
  h += '<div class="pname">'+esc(p.name)+'</div>';
  h += '<div class="period">'+esc(p.period)+'</div>';
  if(p.venue) h += '<div class="meta">📍 '+esc(p.venue)+'</div>';
  if(p.price) h += '<div class="meta">💰 <span class="price">'+esc(p.price)+'</span></div>';
  if(p.summary) h += '<div class="meta">'+esc(p.summary)+'</div>';
  if(p.waiting) h += '<div class="waiting">⏱️ '+esc(p.waiting)+'</div>';
  if(p.venueAlt || p.floorAlt || p.endAlt){
    const bits = [];
    if(p.venueAlt) bits.push("장소 병기: "+p.venueAlt);
    if(p.floorAlt) bits.push("층 병기: "+p.floorAlt);
    if(p.endAlt) bits.push("종료일 병기: "+p.endAlt);
    h += '<div class="alt">⚠️ '+esc(bits.join(" · "))+' — 출처별 안내가 다릅니다.</div>';
  }
  if(p.sources.length){
    h += '<div class="src">출처: ' + p.sources.map(s => s.url
      ? '<a href="'+esc(s.url)+'" target="_blank" rel="noopener">'+esc(s.name)+'</a>'
      : esc(s.name)).join(" · ") + '</div>';
  }
  return h + '</div>';
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
    html += '<div class="floorhead"><span class="fb">'+esc(f)+'</span><span class="fname">'+esc(floorNames[f]||"층 미확인")+' · '+byFloor[f].length+'건</span></div>';
    html += '<div class="grid">' + byFloor[f].map(cardHTML).join("") + '</div>';
  });
  list.innerHTML = html;
  document.getElementById("empty").hidden = shown.length > 0;
  document.getElementById("count").textContent = "총 "+shown.length+"개 팝업 표시 중 (전체 "+DATA.length+"개)";
}

document.getElementById("q").addEventListener("input", e => { curQ = e.target.value.trim().toLowerCase(); render(); });
document.querySelectorAll(".ftab").forEach(b => b.addEventListener("click", () => {
  document.querySelectorAll(".ftab").forEach(x=>x.classList.remove("active"));
  b.classList.add("active"); curFloor = b.dataset.floor; render();
}));
document.querySelectorAll(".chip").forEach(b => b.addEventListener("click", () => {
  document.querySelectorAll(".chip").forEach(x=>x.classList.remove("active"));
  b.classList.add("active"); curCat = b.dataset.cat; render();
}));
render();
</script>
</body>
</html>
`
