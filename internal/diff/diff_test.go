package diff

import (
	"testing"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

func asOf() time.Time { return time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC) }

func TestCompareAddedRemovedUpdated(t *testing.T) {
	prev := &model.Snapshot{
		AsOf: "2026-11-05",
		Popups: []model.Popup{
			{ID: "oki", Name: "오키뷰티", Floor: "B1", EndDate: "2026-11-11"},
			{ID: "old", Name: "어텐스", Floor: "B1", EndDate: "2026-10-26"},
		},
	}
	cur := model.Snapshot{
		AsOf: "2026-11-06",
		Popups: []model.Popup{
			{ID: "oki", Name: "오키뷰티", Floor: "B1", EndDate: "2026-11-11"},
			{ID: "gangneung", Name: "강릉감자닭강정", Floor: "B1", StartDate: "2026-11-06", EndDate: "2026-11-12", FirstSeen: "2026-11-06"},
		},
	}
	set := Compare(prev, cur, asOf())
	if set.Counts["added"] != 1 || set.Counts["removed"] != 1 {
		t.Fatalf("added/removed 계수 오류: %v", set.Counts)
	}
	if set.Changes[0].Name != "강릉감자닭강정" {
		t.Fatalf("신규 항목 오류: %+v", set.Changes[0])
	}
}

func TestFieldChangesConflict(t *testing.T) {
	prev := &model.Snapshot{AsOf: "2026-11-05", Popups: []model.Popup{
		{ID: "peak", Name: "피크퍼포먼스", EndDate: "2026-11-25"},
	}}
	cur := model.Snapshot{AsOf: "2026-11-06", Popups: []model.Popup{
		{ID: "peak", Name: "피크퍼포먼스", EndDate: "2026-12-02", EndDateAlt: "2026-11-25"},
	}}
	set := Compare(prev, cur, asOf())
	if set.Counts["updated"] != 1 {
		t.Fatalf("updated 미감지: %+v", set)
	}
	if set.Changes[0].Field != "종료일" || set.Changes[0].Before != "11-25" || set.Changes[0].After != "12-02" {
		t.Fatalf("변경 내용 오류: %+v", set.Changes[0])
	}
}

func TestFirstRunNoSpam(t *testing.T) {
	cur := model.Snapshot{AsOf: "2026-11-06", Popups: []model.Popup{{ID: "a", Name: "A"}}}
	set := Compare(nil, cur, asOf())
	if len(set.Changes) != 0 {
		t.Fatalf("첫 실행에 변동 알림이 나오면 안 됨: %+v", set.Changes)
	}
}

func TestStatuses(t *testing.T) {
	p := model.Popup{FirstSeen: "2026-11-06", EndDate: "2026-11-08", EndDateAlt: "2026-11-12"}
	st := Statuses(p, asOf())
	want := map[string]bool{"new": false, "ending": false, "conflict": false}
	for _, s := range st {
		want[s] = true
	}
	for k, v := range want {
		if !v {
			t.Fatalf("상태 %s 누락: %v", k, st)
		}
	}
}
