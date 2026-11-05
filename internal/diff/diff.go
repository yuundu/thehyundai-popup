package diff

import (
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// Compare 는 어제 스냅샷과 오늘 목록을 비교해 변동만 뽑는다.
//
//	added   : 오늘 처음 보인다(신규 오픈)
//	removed : 오늘 사라졌다(종료)
//	updated : 기간·층·장소·가격 등이 바뀌었다(변경)
//
// 스냅샷이 없으면(첫 실행) 전부 added 로 취급하지 않고 빈 결과를 돌려
// 첫날 전체가 "신규"로 알림 도배되는 일을 막는다.
func Compare(prev *model.Snapshot, cur model.Snapshot, asOf time.Time) model.ChangeSet {
	set := model.ChangeSet{
		AsOf:    cur.AsOf,
		Counts:  map[string]int{},
		Changes: []model.Change{},
		Notes:   cur.Notes,
	}
	if prev != nil {
		set.PreviousOf = prev.AsOf
	}

	curIdx := map[string]model.Popup{}
	for _, p := range cur.Popups {
		curIdx[p.ID] = p
	}

	if prev != nil {
		prevIdx := map[string]model.Popup{}
		for _, p := range prev.Popups {
			prevIdx[p.ID] = p
		}

		for id, p := range curIdx {
			q, ok := prevIdx[id]
			if !ok {
				c := model.Change{
					Type: model.ChangeAdded, ID: id, Name: p.Name, Floor: p.FloorLabel(),
					Message: addedMessage(p),
				}
				set.Changes = append(set.Changes, c)
				set.Counts[string(model.ChangeAdded)]++
				continue
			}
			if c, changed := fieldChanges(q, p); changed {
				set.Changes = append(set.Changes, c...)
				set.Counts[string(model.ChangeUpdated)]++
			}
		}
		for id, q := range prevIdx {
			if _, ok := curIdx[id]; !ok {
				c := model.Change{
					Type: model.ChangeRemoved, ID: id, Name: q.Name, Floor: q.FloorLabel(),
					Message: removedMessage(q, asOf),
				}
				set.Changes = append(set.Changes, c)
				set.Counts[string(model.ChangeRemoved)]++
			}
		}
	}

	return set
}

// addedMessage 는 신규 팝업 한 줄 소개를 만든다.
func addedMessage(p model.Popup) string {
	s := "🆕 신규 오픈: " + p.Name
	if p.StartDate != "" && p.EndDate != "" {
		s += " (" + short(p.StartDate) + "~" + short(p.EndDate) + ")"
	}
	if p.Venue != "" {
		s += " · " + p.Venue
	}
	if p.Price != "" {
		s += " · " + p.Price
	}
	return s
}

// removedMessage 는 종료 안내를 만든다. 오늘 종료면 "마지막 날"로 표현한다.
func removedMessage(p model.Popup, asOf time.Time) string {
	if p.EndDate == short(asOf.Format("2006-01-02")) {
		return "⏰ 오늘 마지막: " + p.Name + " — 오늘 안 가면 못 봄"
	}
	return "🏁 종료: " + p.Name
}

// fieldChanges 는 같은 팝업의 바뀐 필드를 찾는다. 층·기간·장소·가격만 감시한다.
func fieldChanges(prev, cur model.Popup) ([]model.Change, bool) {
	var out []model.Change
	add := func(field, before, after string) {
		out = append(out, model.Change{
			Type: model.ChangeUpdated, ID: cur.ID, Name: cur.Name, Floor: cur.FloorLabel(),
			Field: field, Before: before, After: after,
			Message: "🔄 변경(" + field + "): " + cur.Name + " — " + before + " → " + after,
		})
	}
	if v, ok := cmp(cur.Floor, prev.Floor, cur.FloorAlt); ok {
		add("층", v[0], v[1])
	}
	if cur.StartDate != prev.StartDate {
		add("시작일", short(prev.StartDate), short(cur.StartDate))
	}
	if cur.EndDate != prev.EndDate {
		add("종료일", short(prev.EndDate), short(cur.EndDate))
	}
	if v, ok := cmp(cur.Venue, prev.Venue, cur.VenueAlt); ok {
		add("장소", v[0], v[1])
	}
	if cur.Price != prev.Price && cur.Price != "" {
		add("가격", prev.Price, cur.Price)
	}
	return out, len(out) > 0
}

// cmp 는 현재값·병기값 중 하나가 과거값과 다르면 (이전, 현재) 쌍을 돌린다.
func cmp(cur, prev, alt string) ([2]string, bool) {
	if cur != prev && cur != "" {
		return [2]string{fallback(prev), cur}, true
	}
	if alt != "" && alt != prev {
		return [2]string{fallback(prev), alt}, true
	}
	return [2]string{}, false
}

func fallback(s string) string {
	if s == "" {
		return "미상"
	}
	return s
}

func short(d string) string {
	if len(d) == 10 && d[4] == '-' {
		return d[5:] // MM-DD
	}
	return d
}

// Statuses 는 기준일과 기간으로 상태 배지를 계산한다(사이트 카드 배지).
func Statuses(p model.Popup, asOf time.Time) []string {
	var out []string
	if p.FirstSeen == asOf.Format("2006-01-02") {
		out = append(out, model.StatusNew)
	}
	if p.EndDate != "" {
		if e, err := time.ParseInLocation("2006-01-02", p.EndDate, asOf.Location()); err == nil {
			d := int(e.Sub(asOf).Hours() / 24)
			switch {
			case d < 0:
				// 기간 종료: 목록에서 자연스러워지므로 배지 없음
			case d == 0:
				out = append(out, model.StatusClosed)
			case d <= 7:
				out = append(out, model.StatusEnding)
			}
		}
	}
	for _, t := range p.Tags {
		if t == model.StatusDeal {
			out = append(out, model.StatusDeal)
			break
		}
	}
	if p.EndDateAlt != "" || p.FloorAlt != "" || p.VenueAlt != "" {
		out = append(out, model.StatusConflict)
	}
	return out
}
