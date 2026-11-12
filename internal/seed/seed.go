// Package seed 는 수동 검증된 기준 스냅샷이다.
//
// 자동 크롤러가 처음 도는 날(2026-11-06) 기준으로, 공식 안내·팝업 플랫폼·언론 보도를
// 사람이 대조해 확정한 목록이다. 출처마다 안내가 다른 항목은 Alt 필드에 병기했다.
//
// 유지 규칙
//   - 자동 크롤링이 정상화되면 이 시드는 점차 자체 데이터로 대체된다.
//   - 신뢰도(confidence)는 자동 수집 결과와 충돌할 때 크롤러 쪽 우선순위 규칙을 따른다.
//   - 연도 필터: 모든 기간은 2026년으로 검증된 것만 둔다(2025년 이월 정보 배제).
package seed

import (
	"strings"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// seedSrc 는 수동 검증 출처의 공통 표기.
var seedSrc = model.Source{
	Name: "수동 검증(공식·플랫폼·언론 대조)",
	Kind: model.KindSeed,
	URL:  "https://thehyundaiseoul.ehyundai.com/",
}

// base 는 시드 공통값.
func base(name, floor, cat, start, end, venue, price, summary string) model.Popup {
	return model.Popup{
		Name: name, Floor: floor, Category: cat,
		StartDate: start, EndDate: end, Venue: venue, Price: price, Summary: summary,
		Tags: []string{}, Sources: []model.Source{seedSrc},
	}
}

// Popups 는 2026-11-06 기준 검증 목록이다.
func Popups() []model.Popup {
	ps := []model.Popup{
		// ── B1 식품관 주간 팝업(오늘 신규 3종) ──────────────────────────
		func() model.Popup {
			p := base("강릉감자닭강정", model.FloorB1, model.CategoryFood,
				"2026-11-06", "2026-11-12", "지하 1층 식품행사장", "15,000원",
				"바삭 감자채 + 닭강정. 강릉 로컬 감성 간식의 백화점 입점판.")
			p.FirstSeen = "2026-11-06"
			return p
		}(),
		func() model.Popup {
			p := base("테쿠루 타코야키", model.FloorB1, model.CategoryFood,
				"2026-11-06", "2026-11-12", "지하 1층 식품행사장", "9,000~10,000원",
				"12시간 숙성 가쓰오 반죽 문어 타코야키.")
			p.FirstSeen = "2026-11-06"
			return p
		}(),
		func() model.Popup {
			p := base("에그릿 에그타르트", model.FloorB1, model.CategoryFood,
				"2026-11-06", "2026-11-12", "지하 1층 식품행사장", "4,800원",
				"강릉 중앙시장 에그타르트.")
			p.FirstSeen = "2026-11-06"
			return p
		}(),

		// ── B1 마켓 특가/웨이팅 이슈 ────────────────────────────────────
		func() model.Popup {
			p := base("Tasty Seoul 마켓 알뜰장보기", model.FloorB1, model.CategoryFood,
				"2026-11-06", "2026-11-08", "지하 1층 테이스티 서울 마켓", "한우 채끝 12,900원 · 사과 19,900원",
				"주말 한정 수량 특가(한우 채끝 스틱·사과 등). 조기 소진 가능.")
			p.Tags = append(p.Tags, model.StatusDeal)
			p.Waiting = "주말 오전 오픈 직후 줄 가장 길어요. 오전 방문 권장."
			p.FirstSeen = "2026-11-06"
			return p
		}(),
		func() model.Popup {
			p := base("올기 찹쌀 초콜릿", model.FloorB1, model.CategoryFood,
				"2026-10-20", "2026-11-16", "지하 1층 핑크스핫도그 옆", "5,500~16,900원",
				"찹쌀 초콜릿 전문 브랜드 팝업(연합뉴스 11월 쇼핑대전 보도 확인).")
			p.FloorAlt = "지하 1층(매장) / 웨이팅 등록: B2 지하철 연결통로 09:00~"
			p.Waiting = "현장 웨이팅은 B2 지하철 연결통로에서 09:00부터 등록."
			return p
		}(),

		// ── B1 팝업존/대행사장 ─────────────────────────────────────────
		base("오키뷰티", model.FloorB1, model.CategoryBeauty,
			"2026-11-04", "2026-11-11", "지하 1층 팝업 존", "", "스킨케어 브랜드 팝업."),
		base("주다안", model.FloorB1, model.CategoryFashion,
			"2026-11-04", "2026-11-11", "지하 1층 대행사장", "", "주다안 팝업."),

		// ── B2 ────────────────────────────────────────────────────────
		base("SM X 아모레퍼시픽 팬시그널 립밤", model.FloorB2, model.CategoryBeauty,
			"2026-11-05", "2026-11-11", "지하 2층 POP-UP@ICONIC", "", "팬시그널 립밤 콜라보 팝업."),
		base("파우 비클린", model.FloorB2, model.CategoryFashion,
			"2026-11-05", "2026-11-11", "지하 2층 비클린 매장", "", "파우 팝업스토어."),
		base("롱샴 마이 플리아쥬", model.FloorB2, model.CategoryFashion,
			"2026-11-05", "2026-11-18", "지하 2층 센트럴 팝업존", "", "롱샴 향수·라이프 팝업."),
		func() model.Popup {
			p := base("이외들 IOEDLE", model.FloorB2, model.CategoryFashion,
				"2026-11-12", "2026-11-18", "지하 2층 POP-UP@WEST", "",
				"1990년대 버스정류장 레트로 콘셉트 팝업.")
			p.ID = "이외들-ioedle"
			p.FirstSeen = "2026-11-12"
			p.Sources = []model.Source{
				{
					Name: "더현대 서울 공식 캘린더",
					URL:  "https://thehyundaiseoul.ehyundai.com/event-info/E4002611514269",
					Kind: model.KindOfficial,
				},
			}
			return p
		}(),
		func() model.Popup {
			p := base("알키알카 ALKI ALKA", model.FloorB2, model.CategoryFashion,
				"2026-11-12", "2026-11-18", "지하 2층 POP-UP@EAST", "최대 45% 할인",
				"여성 패션 첫 팝업스토어. 최대 45% 할인, 구매금액별 사은품 및 리퍼브 1~5만원 균일가.")
			p.ID = "알키알카-alki-alka"
			p.FirstSeen = "2026-11-12"
			p.Sources = []model.Source{
				{
					Name: "더현대 서울 공식 캘린더",
					URL:  "https://thehyundaiseoul.ehyundai.com/event-info/E4002611514786",
					Kind: model.KindOfficial,
				},
			}
			return p
		}(),

		// ── 1F ────────────────────────────────────────────────────────
		func() model.Popup {
			p := base("프라다 뷰티 신제품", model.Floor1F, model.CategoryBeauty,
				"2026-11-05", "", "1층 본매장", "", "신제품 출시 기념 매장 내 팝업(소진 시 종료).")
			p.EndDateAlt = ""
			return p
		}(),

		// ── 4F: 출처 간 종료일 상이 사례(11/25 vs 12/2) 병기 표기 대상 ──
		func() model.Popup {
			p := base("피크퍼포먼스", model.Floor4F, model.CategoryFashion,
				"2026-11-05", "2026-12-02", "4층", "",
				"스포츠·아웃도어 브랜드 팝업. 출처별 종료일 안내가 달라 병기함.")
			p.EndDateAlt = "2026-11-25(일부 출처) / 2026-12-02(기준 스냅샷)"
			return p
		}(),

		// ── 5F ────────────────────────────────────────────────────────
		func() model.Popup {
			p := base("크리스마스 인 뉴욕시티 빌리지", model.Floor5F, model.CategorySeason,
				"2026-11-01", "2026-12-31", "5층 사운즈 포레스트", "무료(사전예약제)",
				"'Christmas in New York City' 테마. 1차 사전예약 진행 중.")
			p.StartDateAlt = "캠페인 전체 11/1~, 더현대 서울 빌리지 입장 가능일 11/10(화)~"
			p.Waiting = "사전예약: 계정당 최대 4매 · 1시간 단위(10:30~18:30) · 예약 15분 초과 시 입장 불가. 예약 실패 시 2~4층 QR 현장 웨이팅 가능."
			return p
		}(),
		base("힐링스쿱 스튜디오", model.Floor5F, model.CategoryFashion,
			"2026-11-04", "2026-11-16", "5층 팝업 스튜디오", "", "체험형 스튜디오 팝업."),
	}
	for i := range ps {
		if ps[i].ID == "" {
			ps[i].ID = slugSeed(ps[i].Name)
		}
		if ps[i].FirstSeen == "" {
			// 레이더 최초 등록일: 시작일이 확인되면 그날, 아니면 기준일
			if ps[i].StartDate != "" {
				ps[i].FirstSeen = ps[i].StartDate
			} else {
				ps[i].FirstSeen = "2026-11-06"
			}
		}
		ps[i].UpdatedAt = "2026-11-06"
		ps[i].Confidence = 80
	}
	return ps
}

// slugSeed 는 사이트 패키지를 끌어오지 않고 ID 를 만들기 위한 최소 정규화.
func slugSeed(name string) string {
	name = strings.ToLower(name)
	var b []rune
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r >= 0xAC00 && r <= 0xD7A3:
			b = append(b, r)
		case r == ' ', r == '-':
			b = append(b, '-')
		}
	}
	return string(b)
}

// Notes 는 사이트 하단·기획서에 함께 노출할 검증 메모.
func Notes() []string {
	return []string{
		"크리스마스 빌리지: 캠페인 전체 기간은 11/1~12/31이지만 더현대 서울 입장 가능일은 11/10(화)부터 — 공식 이벤트 페이지 기준.",
		"11/5(목) 1차 사전예약은 서버 오류로 중단 후 당일 18:00 네이버 예약으로 재진행(현대백화점 공지). 2차 예약 일정은 미공개.",
		"피크퍼포먼스 종료일: 출처별 11/25 vs 12/2 상이 — 공식 확인 전까지 병기 표기 유지.",
	}
}
