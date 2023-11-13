package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:embed data
var content embed.FS

func main() {
	http.HandleFunc("/listseasons.php", listSeasons)
	http.HandleFunc("/showseason.php", showSeason)

	port := "8080"
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func listSeasons(w http.ResponseWriter, r *http.Request) {
	// Only GET allowed
	if r.Method != "GET" {
		http.Error(w, "Not found", 404)
	}

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
}

func showSeason(w http.ResponseWriter, r *http.Request) {
	// Only GET allowed
	if r.Method != "GET" {
		http.Error(w, "Not found", 404)
	}

	season := -1
	nonNumeric := ""
	if len(r.URL.Query()) > 0 {
		sb := strings.Builder{}
		for k, v := range r.URL.Query() {
			if len(v) == 1 && len(v[0]) > 0 {
				if strings.EqualFold(k, "season") {
					var e error
					season, e = strconv.Atoi(v[0])
					if e != nil {
						nonNumeric = v[0]
						log.Printf("GET %s non-numeric season: %q", r.RequestURI, nonNumeric)
						http.Error(w, "season "+nonNumeric+" not found", 404)
					} else {
						log.Printf("GET %s numeric season: %d", r.RequestURI, season)
						if season > -1 {
							b, err := content.ReadFile(fmt.Sprintf("data/seasons/season_%02d.php", season))
							if err != nil {
								log.Printf("ERROR: season %02d not found", season)
								http.Error(w, fmt.Sprintf("season %d not found", season), 404)
							}
							sb.Write(b)
						}
					}
				}
			} else {
				http.Error(w, "season not found", 404)
			}
		}
		if sb.Len() > 0 {
			c, err := w.Write([]byte(sb.String()))
			if err != nil {
				log.Printf("ERROR: %v", err)
			} else {
				log.Printf("%d bytes response to '%s'", c, r.Method+" "+r.Host+r.URL.String())
			}
		}
	} else {
		http.Error(w, fmt.Sprintf("Bad request: %q query parameter required to show season", "seasons"), 400)
	}
}
