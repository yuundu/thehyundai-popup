// Package server 는 로컬 미리보기용 웹 서버(`-serve` 모드)를 제공한다.
package server

import (
	"io/fs"
	"log"
	"net/http"
)

// Serve 는 fsys(보통 main 패키지가 go:embed 로 붙인 out/ 디렉터리)를 지정 포트로 연다.
func Serve(addr string, fsys fs.FS) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(fsys)))
	log.Printf("웹 서버 기동: http://localhost%s  (Ctrl+C 로 종료)", addr)
	return http.ListenAndServe(addr, mux)
}
