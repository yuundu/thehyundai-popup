package crawler

import (
	"sort"
	"strings"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// Tagged 는 "어느 출처에서 왔는지"를 붙인 수집 결과 한 건.
type Tagged struct {
	Popup model.Popup
	Kind  model.SourceKind
}

// priority 는 충돌 시 어느 쪽을 대표값으로 쓸지 정한다.
// 공식 > 플랫폼 > 보도 > 예약/예매 > SNS > 수동 시드 순.
func priority(k model.SourceKind) int {
	switch k {
	case model.KindOfficial:
		return 5
	case model.KindPlatform:
		return 4
	case model.KindPress:
		return 3
	case model.KindReservation:
		return 2
	case model.KindSeed:
		return 1
	case model.KindSNS:
		return 0
	}
	return 0
}

// confidenceOf 는 출처 등급 조합으로 0~100 신뢰도를 매긴다.
// 수동 검증 시드가 섞인 조합은 최소 80점을 보장한다(사람이 확인한 값).
func confidenceOf(kinds map[model.SourceKind]bool) int {
	var score int
	switch {
	case kinds[model.KindOfficial] && kinds[model.KindPlatform]:
		score = 100
	case kinds[model.KindOfficial]:
		score = 92
	case kinds[model.KindPlatform] && kinds[model.KindPress]:
		score = 72
	case kinds[model.KindPlatform]:
		score = 62
	case kinds[model.KindPress]:
		score = 55
	case kinds[model.KindSNS]:
		score = 40
	default:
		score = 80 // 시드만 있는 경우
	}
	if kinds[model.KindSeed] && score < 80 {
		score = 80
	}
	return score
}

// setField 는 우선순위 규칙에 따라 대표값/병기값을 갱신한다.
// 같은 필드에 서로 다른 주장이 오면 지우지 않고 Alt 에 남긴다(누락 금지 원칙).
func setField(dst, alt *string, val string, valPrio, dstPrio int) int {
	val = strings.TrimSpace(val)
	if val == "" {
		return dstPrio
	}
	if *dst == "" {
		*dst = val
		return valPrio
	}
	if *dst == val {
		if valPrio > dstPrio {
			return valPrio
		}
		return dstPrio
	}
	if valPrio > dstPrio {
		if *alt == "" {
			*alt = *dst
		}
		*dst = val
		return valPrio
	}
	if *alt == "" && val != *dst {
		*alt = val
	}
	return dstPrio
}

func addSource(p *model.Popup, s model.Source) {
	for _, ex := range p.Sources {
		if ex.Name == s.Name && ex.URL == s.URL {
			if ex.Claim == "" && s.Claim != "" {
				ex.Claim = s.Claim
			}
			return
		}
	}
	p.Sources = append(p.Sources, s)
}

// agg 는 병합 중인 그룹 하나.
type agg struct {
	p     model.Popup
	prio  map[string]int // 필드별 현재 대표값의 우선순위
	kinds map[model.SourceKind]bool
}

// findFuzzy 는 키가 서로를 포함하는 관계(예: "올기" ⊂ "올기찹쌀초콜릿")나
// 같은 ID 를 가진 기존 그룹을 찾는다. 브랜드명 표기 차이로 같은 팝업이 갈라지는 것을 막는다.
func findFuzzy(groups map[string]*agg, k string, id string) string {
	for gk, g := range groups {
		if gk == k {
			continue
		}
		if id != "" && g.p.ID == id {
			return gk
		}
		if strings.Contains(gk, k) || strings.Contains(k, gk) {
			return gk
		}
	}
	return ""
}

// Merge 는 출처별 수집 결과를 하나의 목록으로 합친다.
//
//  1. 정규화 키(브랜드명, 부분 일치 허용)로 그룹화
//  2. 필드별 우선순위 병합, 불일치는 Alt 필드에 병기
//  3. 출처 등급 조합으로 신뢰도 점수 산정
func Merge(items []Tagged) []model.Popup {
	groups := map[string]*agg{}
	order := []string{}

	for _, it := range items {
		in := it.Popup
		name := strings.TrimSpace(in.Name)
		if name == "" {
			continue
		}
		k := Key(name)
		if fk := findFuzzy(groups, k, in.ID); fk != "" {
			k = fk
		}
		g, ok := groups[k]
		if !ok {
			g = &agg{
				p: model.Popup{
					ID:        in.ID,
					Name:      name,
					Brand:     in.Brand,
					Category:  in.Category,
					FirstSeen: in.FirstSeen,
					Tags:      append([]string{}, in.Tags...),
				},
				prio:  map[string]int{},
				kinds: map[model.SourceKind]bool{},
			}
			if g.p.ID == "" {
				g.p.ID = Slug(name)
			}
			groups[k] = g
			order = append(order, k)
		}
		pr := priority(it.Kind)
		g.kinds[it.Kind] = true

		// 대표명: 더 길고 정보적인 이름을 채택 (예: "올기 찹쌀 초콜릿" > "올기")
		if len(name) > len(g.p.Name) {
			g.p.Name = name
		}

		g.prio["floor"] = setField(&g.p.Floor, &g.p.FloorAlt, in.Floor, pr, g.prio["floor"])
		g.prio["start"] = setField(&g.p.StartDate, &g.p.StartDateAlt, in.StartDate, pr, g.prio["start"])
		g.prio["end"] = setField(&g.p.EndDate, &g.p.EndDateAlt, in.EndDate, pr, g.prio["end"])
		g.prio["venue"] = setField(&g.p.Venue, &g.p.VenueAlt, in.Venue, pr, g.prio["venue"])
		g.prio["price"] = setField(&g.p.Price, new(string), in.Price, pr, g.prio["price"])
		g.prio["summary"] = setField(&g.p.Summary, new(string), in.Summary, pr, g.prio["summary"])
		g.prio["waiting"] = setField(&g.p.Waiting, new(string), in.Waiting, pr, g.prio["waiting"])
		if in.Category != "" && (g.p.Category == "" || pr > g.prio["category"]) {
			g.p.Category = in.Category
			g.prio["category"] = pr
		}
		if g.p.Brand == "" {
			g.p.Brand = in.Brand
		}
		if in.FirstSeen != "" && (g.p.FirstSeen == "" || in.FirstSeen < g.p.FirstSeen) {
			g.p.FirstSeen = in.FirstSeen
		}
		for _, t := range in.Tags {
			if !contains(g.p.Tags, t) {
				g.p.Tags = append(g.p.Tags, t)
			}
		}
		for _, s := range in.Sources {
			addSource(&g.p, s)
		}
	}

	out := make([]model.Popup, 0, len(order))
	// ID 도 최종 대표명 기준으로 재정렬(그룹 합쳐진 경우 대비)
	for _, k := range order {
		g := groups[k]
		if Key(g.p.ID) != k && Key(g.p.Name) != k {
			g.p.ID = Slug(g.p.Name)
		}
		g.p.Confidence = confidenceOf(g.kinds)
		if g.p.Floor == "" {
			g.p.Floor = model.FloorUnknown
		}
		out = append(out, g.p)
	}
	// 층 → 이름 순으로 안정 정렬(수집 순서가 흔들려도 결과가 같게).
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Floor != out[j].Floor {
			return floorRank(out[i].Floor) < floorRank(out[j].Floor)
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func floorRank(f string) int {
	switch f {
	case "B2":
		return 2
	case "B1":
		return 1
	case "1F":
		return 10
	case "2F":
		return 20
	case "5F":
		return 50
	}
	return 99
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
