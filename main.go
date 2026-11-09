// 더현대 서울 팝업 레이더 — 하루 한 번 수집해서 층별 레이더 사이트를 다시 빌드한다.
//
// 동작 모드
//
//	go run main.go                 # 수집 → Diff → 스냅샷 갱신 → out/ 재빌드
//	go run main.go -serve          # 빌드된 out/ 을 http://localhost:8080 으로 미리보기
//	go run main.go -serve -addr :9090
//	go run main.go -dry            # 파일 갱신 없이 수집·비교만 출력
//
// GitHub Actions 에서는 "go run main.go" 가 평일 KST 07:30(일~목 UTC 22:30)에 실행되고,
// out/ 이 GitHub Pages 로 배포된다.
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/crawler"
	"github.com/yuundu/thehyundai-popup/internal/diff"
	"github.com/yuundu/thehyundai-popup/internal/model"
	"github.com/yuundu/thehyundai-popup/internal/notify"
	"github.com/yuundu/thehyundai-popup/internal/seed"
	"github.com/yuundu/thehyundai-popup/internal/server"
	"github.com/yuundu/thehyundai-popup/internal/site"
	"github.com/yuundu/thehyundai-popup/internal/store"
)

// distFS 는 빌드 산출물(out/)을 실행 파일에 붙여 넣는다.
// -serve 모드는 이 임베드 사본을 그대로 서빙한다(빌드 시점 스냅샷).
//
//go:embed all:out
var distFS embed.FS

const (
	storePath  = "data/last_snapshot.json" // 어제 기준 스냅샷(다음 Diff 의 기준)
	historyDir = "data/history"            // 날짜별 보관
	distDir    = "out"                     // Pages 배포 대상
)

// RunOptions 는 한 번의 실행 설정.
type RunOptions struct {
	Serve   bool
	Addr    string
	DryRun  bool
	Timeout time.Duration
}

// todayKST 는 KST 기준 오늘 0시. 밤에 돌려도 "그날의 레이더"가 되도록 한다.
func todayKST() time.Time {
	loc := time.FixedZone("KST", 9*60*60)
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}

func run(o RunOptions) error {
	if o.Serve {
		sub, err := fs.Sub(distFS, "out")
		if err != nil {
			return fmt.Errorf("out/ 임베드 실패: %w", err)
		}
		return server.Serve(o.Addr, sub)
	}

	asOf := todayKST()
	asOfStr := asOf.Format("2006-01-02")
	log.Printf("기준일: %s (%s)", asOfStr, weekdayKO(asOf))

	// 1) 어제 기준 스냅샷 로드(첫 실행이면 nil → Diff 스킵)
	prev, err := store.LoadSnapshot(storePath)
	if err != nil {
		log.Printf("어제 스냅샷 로드 실패(무시하고 진행): %v", err)
		prev = nil
	}

	// 2) 수집: 출처별 goroutine 동시 수집 + 부분 실패 허용
	ctx, cancel := context.WithTimeout(context.Background(), o.Timeout)
	defer cancel()

	client := crawler.NewClient()
	results := client.FetchAll(ctx, crawler.DefaultSources(), asOf)
	items := crawler.TaggedFromResults(results)
	okSources := 0
	for _, r := range results {
		if r.Err == nil {
			okSources++
		}
	}

	// 2-1) 검증 시드는 항상 병합 입력에 포함한다.
	//	   - 출처 등급이 가장 낮아(KindSeed) 자동 수집값과 충돌하면 물러난다
	//	   - 자동 수집이 비어 있는 정보(가격·웨이팅·요약)를 채우고,
	//	     전 출처 실패 시 화면 백지화 방지 역할도 한다
	usingSeed := false
	if len(items) == 0 {
		usingSeed = true
		log.Printf("경고: 자동 수집 성공 출처 0곳 — 검증 시드(%d건)로 대체합니다", len(seed.Popups()))
	}
	for _, p := range seed.Popups() {
		items = append(items, crawler.Tagged{Popup: p, Kind: model.KindSeed})
	}

	// 3) 병합: 출처 우선순위 반영 + 불일치 항목 Alt 병기
	popups := crawler.Merge(items)

	// 3-1) 오늘 기준 이미 종료된 팝업 제외(연도/기간 필터)
	kept := popups[:0]
	for _, p := range popups {
		if p.EndDate != "" {
			if e, err := time.ParseInLocation("2006-01-02", p.EndDate, asOf.Location()); err == nil && e.Before(asOf) {
				continue
			}
		}
		kept = append(kept, p)
	}
	popups = kept

	// 3-2) 상태 배지·D-day 계산(렌더 직전에 한 번)
	for i := range popups {
		popups[i].Status = diff.Statuses(popups[i], asOf)
		if popups[i].EndDate != "" {
			if e, err := time.ParseInLocation("2006-01-02", popups[i].EndDate, asOf.Location()); err == nil {
				if d := int(e.Sub(asOf).Hours() / 24); d >= 0 && d <= 14 {
					dCopy := d
					popups[i].DDay = &dCopy
				}
			}
		}
	}
	for _, n := range seed.Notes() {
		log.Printf("[검증 메모] %s", n)
	}

	cur := model.Snapshot{
		GeneratedAt: time.Now().In(time.FixedZone("KST", 9*60*60)).Format("2006-01-02 15:04 KST"),
		AsOf:        asOfStr,
		Store:       "더현대 서울",
		Popups:      popups,
		Notes:       seed.Notes(),
	}

	// 4) Diff: 어제 대비 변동
	set := diff.Compare(prev, cur, asOf)
	log.Printf("변동: 신규 %d · 종료 %d · 변경 %d", set.Counts["added"], set.Counts["removed"], set.Counts["updated"])

	// 5) 알림 채널: 콘솔+파일 기본, SLACK_WEBHOOK_URL 있으면 슬랙 추가
	for _, ch := range notify.All(os.Getenv("SLACK_WEBHOOK_URL"), "data") {
		if err := ch.Send(set); err != nil {
			log.Printf("[%s] 알림 실패: %v", ch.Name(), err)
		}
	}

	if o.DryRun {
		fmt.Printf("DRY RUN: 스냅샷/사이트를 갱신하지 않음 (수집 성공 출처 %d곳, 팝업 %d건, 시드 대체=%v)\n",
			okSources, len(popups), usingSeed)
		return nil
	}

	// 6) 저장: 스냅샷 + 히스토리
	if err := store.SaveSnapshot(storePath, cur); err != nil {
		return fmt.Errorf("스냅샷 저장 실패: %w", err)
	}
	if err := store.SaveHistory(historyDir, cur); err != nil {
		return fmt.Errorf("히스토리 저장 실패: %w", err)
	}

	// 7) 사이트 빌드: out/index.html + out/popups.json + out/changes.json
	newToday := 0
	for _, p := range popups {
		for _, s := range p.Status {
			if s == model.StatusNew {
				newToday++
				break
			}
		}
	}
	if err := site.WriteIndex(distDir+"/index.html", cur, set, newToday); err != nil {
		return fmt.Errorf("index.html 생성 실패: %w", err)
	}
	if err := store.WriteJSON(distDir+"/popups.json", cur); err != nil {
		return fmt.Errorf("popups.json 생성 실패: %w", err)
	}
	if err := store.WriteJSON(distDir+"/changes.json", set); err != nil {
		return fmt.Errorf("changes.json 생성 실패: %w", err)
	}
	log.Printf("빌드 완료: %s/index.html (팝업 %d건)", distDir, len(popups))
	return nil
}

func main() {
	serve := flag.Bool("serve", false, "빌드된 out/ 을 로컬 웹 서버로 미리보기")
	addr := flag.String("addr", ":8080", "미리보기 서버 주소")
	dry := flag.Bool("dry", false, "파일 갱신 없이 수집·비교만 수행")
	timeout := flag.Duration("timeout", 40*time.Second, "수집 전체 타임아웃")
	flag.Parse()

	if err := run(RunOptions{Serve: *serve, Addr: *addr, DryRun: *dry, Timeout: *timeout}); err != nil {
		log.Fatal(err)
	}
}

// weekdayKO 는 로그용 요일.
func weekdayKO(t time.Time) string {
	return []string{"일", "월", "화", "수", "목", "금", "토"}[int(t.Weekday())]
}
