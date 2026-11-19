package diff

import (
	"testing"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// 같은 날 재실행에서 이번 수집 변동이 0건이어도, 아침에 기록한 당일 변동
// (11-19의 8건 종료)이 빈 목록으로 덮어써지지 않아야 한다.
func TestMergeDayKeepsMorningChangesWhenRerunIsEmpty(t *testing.T) {
	morning := model.ChangeSet{
		AsOf:       "2026-11-19",
		PreviousOf: "2026-11-18",
		Counts:     map[string]int{"removed": 2},
		Changes: []model.Change{
			{Type: model.ChangeRemoved, ID: "999휴머니티", Name: "999휴머니티", Floor: "B2", Message: "🏁 종료: 999휴머니티"},
			{Type: model.ChangeRemoved, ID: "롱샴-마이-플리지", Name: "롱샴", Floor: "B2", Message: "🏁 종료: 롱샴"},
		},
	}
	// 같은 날 재실행: 기준 스냅샷이 이미 오늘자라 변동 0건, previous_of 도 오늘.
	fresh := model.ChangeSet{
		AsOf:       "2026-11-19",
		PreviousOf: "2026-11-19",
		Counts:     map[string]int{},
		Changes:    []model.Change{},
	}
	today := []model.Popup{{ID: "코나카드", Name: "코나카드 아틀리에"}}

	merged := MergeDay(morning, fresh, today)

	if len(merged.Changes) != 2 {
		t.Fatalf("당일 변동이 보존되지 않음: %d건 %+v", len(merged.Changes), merged.Changes)
	}
	if merged.Counts["removed"] != 2 {
		t.Fatalf("counts 재계산 오류: %v", merged.Counts)
	}
	if merged.PreviousOf != "2026-11-18" {
		t.Fatalf("기준일이 오늘로 덮어써짐: %q", merged.PreviousOf)
	}
}

// 오후 재실행에서 새 변동이 잡히면 기존 변동에 더해져야 한다.
func TestMergeDayAppendsNewChanges(t *testing.T) {
	morning := model.ChangeSet{
		AsOf:       "2026-11-19",
		PreviousOf: "2026-11-18",
		Counts:     map[string]int{"removed": 1},
		Changes: []model.Change{
			{Type: model.ChangeRemoved, ID: "old", Name: "종료된 팝업", Floor: "B1", Message: "🏁 종료: 종료된 팝업"},
		},
	}
	fresh := model.ChangeSet{
		AsOf:       "2026-11-19",
		PreviousOf: "2026-11-18",
		Counts:     map[string]int{"added": 1},
		Changes: []model.Change{
			{Type: model.ChangeAdded, ID: "snows", Name: "SNOWS", Floor: "1F", Message: "🆕 신규 오픈: SNOWS"},
		},
	}
	today := []model.Popup{{ID: "snows"}, {ID: "코나카드"}}

	merged := MergeDay(morning, fresh, today)

	if len(merged.Changes) != 2 {
		t.Fatalf("병합 건수 오류: %+v", merged.Changes)
	}
	if merged.Changes[0].ID != "snows" || merged.Changes[1].ID != "old" {
		t.Fatalf("이번 실행 변동이 먼저 오지 않음: %+v", merged.Changes)
	}
	if merged.Counts["added"] != 1 || merged.Counts["removed"] != 1 {
		t.Fatalf("counts 오류: %v", merged.Counts)
	}
}

// 오늘 목록에 아직 있는 팝업의 기존 '종료' 표기는 버린다(오전 수집 오류로 생긴 오탐 정리).
func TestMergeDayDropsContradictedRemoval(t *testing.T) {
	morning := model.ChangeSet{
		AsOf: "2026-11-19",
		Changes: []model.Change{
			{Type: model.ChangeRemoved, ID: "alive", Name: "살아있는 팝업", Message: "🏁 종료: 살아있는 팝업"},
			{Type: model.ChangeRemoved, ID: "gone", Name: "정말 종료", Message: "🏁 종료: 정말 종료"},
		},
	}
	fresh := model.ChangeSet{AsOf: "2026-11-19", Changes: []model.Change{}}
	today := []model.Popup{{ID: "alive"}}

	merged := MergeDay(morning, fresh, today)
	if len(merged.Changes) != 1 || merged.Changes[0].ID != "gone" {
		t.Fatalf("모순 항목 정리 실패: %+v", merged.Changes)
	}
}

// 같은 팝업·같은 필드의 변경은 이번 실행(최신 수집값)이 이긴다.
func TestMergeDayFreshUpdateWins(t *testing.T) {
	morning := model.ChangeSet{
		AsOf: "2026-11-19",
		Changes: []model.Change{
			{Type: model.ChangeUpdated, ID: "p", Name: "P", Field: "층", Before: "B1", After: "B2", Message: "🔄 변경(층): P — B1 → B2"},
		},
		Counts: map[string]int{"updated": 1},
	}
	fresh := model.ChangeSet{
		AsOf: "2026-11-19",
		Changes: []model.Change{
			{Type: model.ChangeUpdated, ID: "p", Name: "P", Field: "층", Before: "B1", After: "1F", Message: "🔄 변경(층): P — B1 → 1F"},
		},
		Counts: map[string]int{"updated": 1},
	}
	merged := MergeDay(morning, fresh, []model.Popup{{ID: "p"}})
	if len(merged.Changes) != 1 || merged.Changes[0].After != "1F" {
		t.Fatalf("최신값 우선 실패: %+v", merged.Changes)
	}
	if merged.Counts["updated"] != 1 {
		t.Fatalf("updated 는 팝업 단위로 세야 함: %v", merged.Counts)
	}
}

// 날짜가 다른 파일은 병합 대상이 아니다.
func TestMergeDayOtherDayPassthrough(t *testing.T) {
	old := model.ChangeSet{AsOf: "2026-11-18", Changes: []model.Change{
		{Type: model.ChangeRemoved, ID: "x", Name: "X", Message: "🏁 종료: X"},
	}}
	fresh := model.ChangeSet{AsOf: "2026-11-19", Changes: []model.Change{}}
	merged := MergeDay(old, fresh, nil)
	if len(merged.Changes) != 0 {
		t.Fatalf("다른 날짜 변동이 섞임: %+v", merged.Changes)
	}
}

// 같은 입력이면 변동 순서도 항상 같아야 한다(재실행 시 파일이 바뀐 것처럼 보이지 않도록).
func TestCompareIsDeterministic(t *testing.T) {
	prev := &model.Snapshot{AsOf: "2026-11-19", Popups: []model.Popup{}}
	cur := model.Snapshot{AsOf: "2026-11-20"}
	for _, id := range []string{"z-9", "a-1", "m-5", "k-3", "b-2", "y-8", "c-4", "n-6"} {
		cur.Popups = append(cur.Popups, model.Popup{ID: id, Name: "팝업 " + id, Floor: "B2"})
	}
	as := time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)

	first := Compare(prev, cur, as)
	for i := 0; i < 20; i++ {
		again := Compare(prev, cur, as)
		if len(again.Changes) != len(first.Changes) {
			t.Fatalf("건수 변화: %d vs %d", len(again.Changes), len(first.Changes))
		}
		for j := range first.Changes {
			if again.Changes[j].ID != first.Changes[j].ID {
				t.Fatalf("%d번째 실행에서 순서가 달라짐: %s vs %s", i, again.Changes[j].ID, first.Changes[j].ID)
			}
		}
	}
	// 이름 오름차순(신규→종료→변경, 같은 종류면 층·이름)으로 정렬된다.
	if first.Changes[0].ID != "a-1" || first.Changes[len(first.Changes)-1].ID != "z-9" {
		t.Fatalf("정렬 규칙 위반: %+v", first.Changes)
	}
}
