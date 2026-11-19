package diff

import "github.com/yuundu/thehyundai-popup/internal/model"

// MergeDay 는 같은 날 두 번째 이후 실행에서 당일 변동 내역을 잃지 않도록
// 기존 파일의 변동과 이번 실행의 변동을 합친다.
//
// 원칙
//  1. 이번 실행 결과(fresh)를 기준으로 삼는다. 기준일(어제) 대비 전수 재계산이므로
//     가장 정확하고, 같은 입력이면 같은 결과가 나온다(멱등).
//  2. 기존 파일(existing)의 항목은 fresh 에 없는 것만 덧붙인다. 같은 팝업·같은 필드의
//     변경은 fresh 쪽(최신 수집값)을 남긴다.
//  3. 오늘 목록(today)과 모순되는 기존 항목은 버린다.
//     - 종료(removed)라고 적혔지만 오늘 목록에 아직 있는 팝업 → 되살아난 것
//     - 신규(added)라고 적혔지만 오늘 목록에 없는 팝업 → 사라진 것
//  4. 날짜가 다른 파일(자정을 넘긴 재실행)은 합치지 않고 fresh 를 그대로 쓴다.
func MergeDay(existing, fresh model.ChangeSet, today []model.Popup) model.ChangeSet {
	if existing.AsOf != "" && existing.AsOf != fresh.AsOf {
		return fresh
	}
	if len(existing.Changes) == 0 {
		return fresh
	}

	merged := fresh
	if merged.Changes == nil {
		merged.Changes = []model.Change{}
	}

	seen := make(map[string]bool, len(merged.Changes)+len(existing.Changes))
	for _, c := range merged.Changes {
		seen[changeKey(c)] = true
	}
	present := make(map[string]bool, len(today))
	for _, p := range today {
		present[p.ID] = true
	}

	for _, c := range existing.Changes {
		if seen[changeKey(c)] {
			continue
		}
		switch c.Type {
		case model.ChangeRemoved:
			if present[c.ID] {
				continue // 오늘 다시 확인됨 → 기존 '종료' 표기는 폐기
			}
		case model.ChangeAdded:
			if !present[c.ID] {
				continue
			}
		}
		merged.Changes = append(merged.Changes, c)
		seen[changeKey(c)] = true
	}

	// 기준일 표기: 이번 실행이 오늘자 스냅샷과 비교해 previous_of == as_of 가 됐다면
	// 기존 파일이 갖고 있던 진짜 기준일(어제)을 살린다.
	if (merged.PreviousOf == "" || merged.PreviousOf == merged.AsOf) &&
		existing.PreviousOf != "" && existing.PreviousOf != merged.AsOf {
		merged.PreviousOf = existing.PreviousOf
	}

	merged.Counts = CountChanges(merged.Changes)
	return merged
}

// CountChanges 는 변동 목록에서 counts 를 다시 계산한다.
// Compare 와 같은 규칙으로, updated 는 변동 건수가 아니라 팝업(id) 단위로 센다.
func CountChanges(changes []model.Change) map[string]int {
	counts := make(map[string]int)
	updatedIDs := make(map[string]bool)
	for _, c := range changes {
		if c.Type == model.ChangeUpdated {
			updatedIDs[c.ID] = true
			continue
		}
		counts[string(c.Type)]++
	}
	if n := len(updatedIDs); n > 0 {
		counts[string(model.ChangeUpdated)] = n
	}
	return counts
}

// changeKey 는 같은 변동을 가리키는 키.
// updated 는 팝업·필드당 하나만 남기고(최신값 우선), added/removed 는 팝업당 하나.
func changeKey(c model.Change) string {
	if c.Type == model.ChangeUpdated {
		return string(c.Type) + "|" + c.ID + "|" + c.Field
	}
	return string(c.Type) + "|" + c.ID
}
