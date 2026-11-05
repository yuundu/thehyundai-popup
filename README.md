# THE HYUNDAI SEOUL POP-UP RADAR 🛰️

**매일 달라지는 더현대 서울 팝업, 층별로 한눈에.**
평일 아침 08:49(KST)에 Go 크롤러가 팝업 정보를 전수 수집해 어제 기록과 비교(Diff)하고,
변동 건만 하이라이트한 층별 레이더 사이트를 GitHub Pages 에 자동 배포합니다.

- 🌐 배포: https://yuundu.github.io/thehyundai-popup/
- 🏗️ 레포: https://github.com/yuundu/thehyundai-popup
- ⏰ 수집 스케줄: 평일(월~금) 08:49 KST = UTC 23:49 (GitHub Actions cron)

## 설치(최초 1회)

1. 레포 **Settings → Pages → Source: GitHub Actions** 선택.
2. `.install/deploy.yml` 을 **`.github/workflows/deploy.yml` 로 추가**.
   (Pickle 의 GitHub 연동 토큰은 워크플로 파일 push 권한이 없어, 파일 생성은
   GitHub 웹에서 한 번 직접 해야 합니다 — Actions 탭에서도 추가 가능)
3. 이후 매일 평일 08:49 KST 에 자동 수집·배포됩니다. `Actions → Run workflow`
   로 즉시 테스트 실행도 가능합니다.

## 동작 방식 (한 장 요약)

```
GitHub Actions cron (UTC 23:49, 평일)
        │
        ▼
 go run main.go
   ├─ 1. crawler: 출처별 동시 수집(goroutine) — 일부 출처가 죽어도 계속
   ├─ 2. merge:   출처 우선순위 병합(공식>플랫폼>언론>SNS) + 불일치 값 병기(Alt)
   ├─ 3. diff:    어제 data/last_snapshot.json 과 비교 → 신규/종료/변경
   ├─ 4. notify:  콘솔(Actions 요약) + data/changes-YYYY-MM-DD.json (+슬랙 웹훅 설정 시)
   └─ 5. site:    out/index.html + out/popups.json 재생성
        │
        ▼
 data/last_snapshot.json 커밋(내일 Diff 기준) → out/ 을 GitHub Pages 배포
```

## 로컬 실행

```bash
go run main.go              # 수집 → 스냅샷 갱신 → out/ 재빌드
go run main.go -dry         # 파일 갱신 없이 수집·비교만 출력
go run main.go -serve       # 빌드된 out/ 을 http://localhost:8080 미리보기
go run main.go -serve -addr :9090
go test ./...               # 날짜 파서·병합·Diff 단위 테스트
```

- `-serve` 모드는 빌드 시점의 `out/` 을 `go:embed` 로 서빙합니다.
- 슬랙 알림: `SLACK_WEBHOOK_URL` 환경변수만 넣으면 변동 요약이 전송됩니다.

## 데이터 규칙 (중요)

1. **연도 필터**: 2025년 등 과년도 행사가 섞이지 않게 모든 기간을 기준 연도로 검증합니다.
   크롤러는 종료일이 기준일보다 과거인 항목을 목록에서 제외합니다.
2. **정보 상이 병기**: 출처마다 층·일정·장소 안내가 다르면 하나로 자르지 않고 둘 다 남깁니다.
   - 카드에 `⚠️ 정보 상이 확인` 배지 + 병기 문구 (예: 올기 — 매장 B1 / 웨이팅 등록 B2)
   - 피크퍼포먼스 종료일 11/25 vs 12/2 사례가 대표적입니다.
3. **부분 실패 허용**: 한 출처가 막혀도 나머지 + 수동 검증 시드로 화면이 유지됩니다.
   사이트가 백지화되는 일은 없습니다.
4. **신뢰도 점수**: 공식+플랫폼 교차 검증 100점 → 단일 SNS 40점까지 카드에 표기합니다.

## 수집 출처

| 등급 | 출처 | 비고 |
|---|---|---|
| official | 더현대 서울 공식몰, 현대백화점 이벤트 페이지 | |
| platform | 팝업나우(기본 활성), 팝플리·팝놀·스팟츠·팝가·디스위켄드 | SPA 사이트는 렌더링 필요 → 기본 비활성(확장 항목) |
| press | 보도자료·언론 기사 | 시드에 반영 |
| sns | 인스타그램 @thehyundai_seoul, 네이버 블로그, X | 로그인·차단 정책상 자동 수집 제외, 수동 확인 채널 |
| seed | 사람이 대조해 확정한 기준 스냅샷(`internal/seed`) | 자동 수집값과 충돌 시 물러나는 최하위 출처 |

## 저장소 구조

```
main.go                     엔트리: -serve / -dry 플래그
internal/crawler/           수집(출처 파서·병합·날짜 파서)
internal/diff/              어제↔오늘 비교 + 상태 배지 계산
internal/model/             데이터 모델(Popup, Snapshot, ChangeSet)
internal/notify/            알림 채널(콘솔·파일·슬랙)
internal/seed/              수동 검증 기준 스냅샷(2026-11-06)
internal/server/            -serve 로컬 미리보기 서버
internal/site/              index.html 템플릿 렌더러
internal/store/             스냅샷/히스토리 저장(원자적 쓰기)
data/last_snapshot.json     어제 기준 스냅샷(Diff 기준 — CI 가 커밋해 보존)
data/history/               날짜별 보관
data/changes-*.json         일일 변동 기록(알림 로그)
out/                        배포 산출물(index.html, popups.json, changes.json)
.github/workflows/deploy.yml  평일 08:49 KST 수집 & Pages 배포
```

## 주의

- 수집 결과는 공식 안내를 대체하지 않습니다. 방문 전 더현대 서울 공식 채널을 확인하세요.
- 인스타그램·네이버 예약 등 로그인이 필요한 채널은 자동 수집 대상에서 제외했습니다.
