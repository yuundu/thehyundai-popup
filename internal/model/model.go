// Package model 은 레이더가 다루는 모든 데이터 구조를 정의한다.
//
// 설계 원칙
//   - 모든 필드는 JSON 태그를 가지며, 그대로 popups.json 으로 직렬화된다.
//   - "확인되지 않은 값"은 빈 문자열로 두고, 절대 추정값으로 채우지 않는다.
//   - 출처가 서로 다른 값을 주장하면 Alt 필드에 병기한다(누락 금지).
package model

// 층 코드. 화면 표기와 수집 표기를 분리하기 위해 상수로 고정한다.
const (
	FloorB2      = "B2"
	FloorB1      = "B1"
	Floor1F      = "1F"
	Floor4F      = "4F"
	Floor5F      = "5F"
	FloorUnknown = "UNKNOWN"
)

// 카테고리 코드.
const (
	CategoryFood    = "food"
	CategoryBeauty  = "beauty"
	CategoryFashion = "fashion"
	CategorySeason  = "season"
)

// 상태 배지 코드.
const (
	StatusNew      = "new"      // 오늘 신규 오픈
	StatusEnding   = "ending"   // 종료 임박(D-7 이내)
	StatusDeal     = "deal"     // 한정 특가
	StatusConflict = "conflict" // 출처 간 정보 불일치
	StatusClosed   = "closed"   // 오늘 종료
)

// SourceKind 는 출처의 신뢰 등급을 구분한다. 우선순위 official > platform > press > sns.
type SourceKind string

const (
	KindOfficial    SourceKind = "official"    // 더현대 서울/현대백화점 공식
	KindPlatform    SourceKind = "platform"    // 팝업 전문 플랫폼
	KindPress       SourceKind = "press"       // 보도자료/언론
	KindSNS         SourceKind = "sns"         // 인스타그램, X, 블로그
	KindReservation SourceKind = "reservation" // 네이버 예약 등
	KindSeed        SourceKind = "seed"        // 수동 검증된 기준 스냅샷
)

// Source 는 개별 출처와, 그 출처가 주장하는 내용(Claim)을 함께 남긴다.
// Claim 은 불일치 병기 화면에서 "어느 출처가 무엇을 말했는지" 보여주기 위한 필드다.
type Source struct {
	Name  string     `json:"name"`
	URL   string     `json:"url,omitempty"`
	Kind  SourceKind `json:"kind"`
	Claim string     `json:"claim,omitempty"`
}

// Popup 은 팝업 한 건. ID 는 브랜드/이벤트명을 정규화해 만든 안정 키다.
type Popup struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Brand    string `json:"brand,omitempty"`
	Floor    string `json:"floor"`
	Category string `json:"category"`
	// StartDate/EndDate 는 "2026-11-06" 형식. 소진 시 종료처럼 미정이면 빈 값.
	StartDate    string   `json:"start_date,omitempty"`
	EndDate      string   `json:"end_date,omitempty"`
	StartDateAlt string   `json:"start_date_alt,omitempty"`
	EndDateAlt   string   `json:"end_date_alt,omitempty"`
	FloorAlt     string   `json:"floor_alt,omitempty"`
	Venue        string   `json:"venue,omitempty"`
	VenueAlt     string   `json:"venue_alt,omitempty"`
	Price        string   `json:"price,omitempty"`
	Summary      string   `json:"summary,omitempty"`
	Waiting      string   `json:"waiting,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Sources      []Source `json:"sources,omitempty"`
	FirstSeen    string   `json:"first_seen,omitempty"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
	Confidence   int      `json:"confidence"`
	// 이하 필드는 빌드 시 계산된다(수집 데이터에는 없음).
	Status []string `json:"status,omitempty"`
	DDay   *int     `json:"d_day,omitempty"`
}

// FloorLabel 은 화면에 노출할 층 라벨.
func (p Popup) FloorLabel() string {
	switch p.Floor {
	case FloorB2:
		return "B2"
	case FloorB1:
		return "B1"
	case Floor1F:
		return "1F"
	case Floor5F:
		return "5F"
	default:
		return "층 미확인"
	}
}

// Snapshot 은 특정 시각의 전수 목록. 파일로 저장해 다음 날 Diff 기준이 된다.
type Snapshot struct {
	GeneratedAt string   `json:"generated_at"`
	AsOf        string   `json:"as_of"`
	Store       string   `json:"store"`
	Popups      []Popup  `json:"popups"`
	Notes       []string `json:"notes,omitempty"`
}

// ChangeType 은 Diff 결과의 종류.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change 는 하루치 변동 한 건. 알림 메시지의 최소 단위다.
type Change struct {
	Type    ChangeType `json:"type"`
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Floor   string     `json:"floor"`
	Field   string     `json:"field,omitempty"`
	Before  string     `json:"before,omitempty"`
	After   string     `json:"after,omitempty"`
	Message string     `json:"message"`
}

// ChangeSet 은 하루치 변동 묶음. 화면 상단 요약 배지와 알림 본문이 여기서 나온다.
type ChangeSet struct {
	AsOf       string         `json:"as_of"`
	PreviousOf string         `json:"previous_of,omitempty"`
	Counts     map[string]int `json:"counts"`
	Changes    []Change       `json:"changes"`
	Notes      []string       `json:"notes,omitempty"`
}

// Count 는 특정 변동 종류의 건수.
func (c ChangeSet) Count(t ChangeType) int { return c.Counts[string(t)] }
