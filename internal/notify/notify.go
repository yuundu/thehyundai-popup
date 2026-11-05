// Package notify 는 알림 채널이다. 1차는 콘솔/파일, 확장 채널로 슬랙 웹훅을 제공한다.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuundu/thehyundai-popup/internal/model"
)

// Notifier 는 알림 채널의 공통 인터페이스. 채널 추가는 이 인터페이스만 구현하면 된다.
type Notifier interface {
	Name() string
	Send(set model.ChangeSet) error
}

// Console 채널: 빌드 로그로 알림을 남긴다(GitHub Actions 요약에도 그대로 보임).
type Console struct{}

func (Console) Name() string { return "console" }

func (Console) Send(set model.ChangeSet) error {
	if len(set.Changes) == 0 {
		fmt.Println("[알림] 오늘 변동 없음")
		return nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "📣 더현대 서울 팝업 레이더 · %s 기준 변동 %d건\n", set.AsOf, len(set.Changes))
	for _, c := range set.Changes {
		b.WriteString("  - " + c.Message + "\n")
	}
	fmt.Print(b.String())
	return nil
}

// File 채널: 변동 내역을 data/changes-YYYY-MM-DD.json 으로 남긴다.
type File struct{ Dir string }

func (f File) Name() string { return "file" }

func (f File) Send(set model.ChangeSet) error {
	if f.Dir == "" {
		return nil
	}
	if err := os.MkdirAll(f.Dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(f.Dir, "changes-"+set.AsOf+".json")
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// Slack 채널: SLACK_WEBHOOK_URL 이 설정돼 있으면 변동 요약을 보낸다.
// 웹훅이 없거나 전송이 실패해도 수집/배포 파이프라인은 계속 진행된다.
type Slack struct {
	WebhookURL string
	HTTP       *http.Client
}

func (s Slack) Name() string { return "slack" }

// Send 는 슬랙 incoming webhook 으로 텍스트 하나를 보낸다.
func (s Slack) Send(set model.ChangeSet) error {
	if s.WebhookURL == "" || len(set.Changes) == 0 {
		return nil
	}
	if s.HTTP == nil {
		s.HTTP = &http.Client{Timeout: 10 * time.Second}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "*더현대 팝업 레이더* · %s 기준 변동 %d건\n", set.AsOf, len(set.Changes))
	for _, c := range set.Changes {
		b.WriteString("• " + c.Message + "\n")
	}
	payload, _ := json.Marshal(map[string]string{"text": b.String()})
	resp, err := s.HTTP.Post(s.WebhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack webhook HTTP %d", resp.StatusCode)
	}
	return nil
}

// All 은 설정된 채널을 한데 모은다.
func All(slackWebhook, fileDir string) []Notifier {
	channels := []Notifier{Console{}, File{Dir: fileDir}}
	if slackWebhook != "" {
		channels = append(channels, Slack{WebhookURL: slackWebhook})
	}
	return channels
}
