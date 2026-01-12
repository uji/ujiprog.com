package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/syumai/workers"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline'; img-src 'self' blob:;")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		res := `<!doctype html>
  <html lang="en">
    <body class="flex flex-col items-center gap-y-5">
      <main class="flex flex-col items-center gap-y-5">
      </main>
      <footer>
        <a
          href="https://github.com/uji/ujiprog.com"
          class="text-blue-600 dark:text-blue-500 hover:underline"
        >
          github.com/uji/ujiprog.com
        </a>
      </footer>
    </body>
  </html>`
		io.Copy(w, bytes.NewReader([]byte(res)))
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
