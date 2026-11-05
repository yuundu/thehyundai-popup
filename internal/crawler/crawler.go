// Package crawler 는 더현대 서울 팝업 정보를 여러 출처에서 수집·병합한다.
package crawler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yuundu/thehyundai-popup/internal/model"
)

// ErrNoSource 는 모든 출처가 실패했을 때의 명시적 실패 신호다.
var ErrNoSource = errors.New("crawler: 수집 성공한 출처가 없음")

const userAgent = "Mozilla/5.0 (compatible; TheHyundaiPopupRadar/1.0; +https://yuundu.github.io/thehyundai-popup/)"

// Result 는 출처 하나의 수집 결과(성공/실패 모두 기록).
type Result struct {
	Source string
	URL    string
	Popups []model.Popup
	Err    error
}

// Client 는 HTML 가져오기와 파싱을 담당한다.
type Client struct {
	HTTP    *http.Client
	Timeout time.Duration
}

// NewClient 는 기본 타임아웃(15초)의 클라이언트를 만든다.
func NewClient() *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: 15 * time.Second},
		Timeout: 15 * time.Second,
	}
}

// FetchOne 은 출처 하나를 가져와 파싱한다. 실패는 에러로만 남기고 전체를 죽이지 않는다.
func (c *Client) FetchOne(ctx context.Context, src Source, asOf time.Time) Result {
	res := Result{Source: src.Name, URL: src.URL}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src.URL, nil)
	if err != nil {
		res.Err = fmt.Errorf("요청 생성 실패: %w", err)
		return res
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		res.Err = fmt.Errorf("요청 실패: %w", err)
		return res
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		res.Err = fmt.Errorf("HTTP %d", resp.StatusCode)
		return res
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		res.Err = fmt.Errorf("HTML 파싱 실패: %w", err)
		return res
	}
	res.Popups = src.Parse(doc, asOf, model.Source{Name: src.Name, URL: src.URL, Kind: src.Kind})
	if res.Popups == nil {
		res.Popups = []model.Popup{}
	}
	return res
}

// FetchAll 은 모든 출처를 goroutine 으로 동시 수집하고, 각 출처의 성패를 로그로 남긴다.
// 컨텍스트 취소·타임아웃·일부 출처 실패와 무관하게 성공한 출처의 결과는 모두 돌려준다.
func (c *Client) FetchAll(ctx context.Context, sources []Source, asOf time.Time) []Result {
	results := make(chan Result, len(sources))
	sem := make(chan struct{}, 4) // 동시 요청 상한: 상대 서버 배려

	for _, src := range sources {
		if !src.Enabled {
			continue
		}
		go func(s Source) {
			defer func() {
				if r := recover(); r != nil {
					results <- Result{Source: s.Name, URL: s.URL, Err: fmt.Errorf("파서 패닉: %v", r)}
				}
			}()
			sem <- struct{}{}
			defer func() { <-sem }()
			results <- c.FetchOne(ctx, s, asOf)
		}(src)
	}

	out := []Result{}
	nEnabled := 0
	for _, s := range sources {
		if s.Enabled {
			nEnabled++
		}
	}
	for i := 0; i < nEnabled; i++ {
		r := <-results
		if r.Err != nil {
			log.Printf("[수집 실패] %s (%s): %v", r.Source, r.URL, r.Err)
		} else {
			log.Printf("[수집 성공] %s: %d건", r.Source, len(r.Popups))
		}
		out = append(out, r)
	}
	return out
}

// Tagged 는 수집 결과에 출처 등급을 붙인 형태(병합 입력).
func TaggedFromResults(rs []Result) []Tagged {
	var items []Tagged
	for _, r := range rs {
		if r.Err != nil {
			continue
		}
		kind := model.KindPlatform
		for _, s := range DefaultSources() {
			if s.Name == r.Source {
				kind = s.Kind
				break
			}
		}
		for _, p := range r.Popups {
			items = append(items, Tagged{Popup: p, Kind: kind})
		}
	}
	return items
}
