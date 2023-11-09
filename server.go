package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//go:embed data
var content embed.FS

func main() {
	http.HandleFunc("/listseasons.php", func(w http.ResponseWriter, r *http.Request) {
		b, err := content.ReadFile("data/listseasons.php")
		if err != nil {
			log.Printf("ERROR: %v", err)
		}

		c, err := w.Write(b)
		if err != nil {
			log.Printf("ERROR: %v", err)
		} else {
			log.Printf("%d bytes written for '%s'", c, r.Method+" "+r.Host+r.URL.String())
		}
	})

	http.HandleFunc("/showseason.php", func(w http.ResponseWriter, r *http.Request) {
		sb := strings.Builder{}

		if len(r.URL.Query()) > 0 {
			for k, v := range r.URL.Query() {
				if strings.EqualFold(k, "season") {
					sb.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ",")))
				}
			}
		}

		if len(sb.String()) == 0 {
			sb.WriteString("No query parameters.")
		}

		c, err := w.Write([]byte(sb.String()))
		if err != nil {
			log.Printf("ERROR: %v", err)
		} else {
			log.Printf("%d bytes response to '%s'", c, r.Method+" "+r.Host+r.URL.String())
		}
	})

	port := "8080"
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}
