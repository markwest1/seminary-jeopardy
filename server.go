package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//go:embed data
var content embed.FS

const port = "8080"

func main() {
	http.HandleFunc("/listseasons.php", listSeasons)
	http.HandleFunc("/showseason.php", showSeason)
	http.HandleFunc("/showgame.php", showGame)

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func listSeasons(w http.ResponseWriter, r *http.Request) {
	// Only GET allowed
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	b, err := readFileReplacingHost("data/listseasons.php")
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
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	season := -1
	nonNumeric := ""
	sb := strings.Builder{}
	if len(r.URL.Query()) > 0 {
		for k, v := range r.URL.Query() {
			if strings.EqualFold("season", k) && len(v) == 1 && len(v[0]) > 0 {
				var e error
				season, e = strconv.Atoi(v[0])
				if e != nil {
					nonNumeric = v[0]
					log.Printf("GET %s non-numeric season: %q", r.RequestURI, nonNumeric)
					http.Error(w, "season "+nonNumeric+" not found", http.StatusNotFound)
					return
				} else {
					log.Printf("GET %s numeric season: %d", r.RequestURI, season)
					if season > -1 {
						b, err := readFileReplacingHost(fmt.Sprintf("data/seasons/season_%02d.php", season))
						if err != nil {
							log.Printf("ERROR: season %02d not found", season)
							http.Error(w, fmt.Sprintf("season %d not found", season), http.StatusNotFound)
							return
						}
						sb.Write(b)
					}
				}
			} else {
				http.Error(w, "", http.StatusNotFound)
				return
			}
		}
	}

	if sb.Len() > 0 {
		c, err := w.Write([]byte(sb.String()))
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		log.Printf("%d bytes response to '%s'", c, r.Method+" "+r.URL.String())
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func showGame(w http.ResponseWriter, r *http.Request) {
	// Only GET allowed
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	gameId := -1
	nonNumeric := ""
	sb := strings.Builder{}
	if len(r.URL.Query()) > 0 {
		for k, v := range r.URL.Query() {
			if strings.EqualFold("game_id", k) && len(v) == 1 && len(v[0]) > 0 {
				var err error
				gameId, err = strconv.Atoi(v[0])
				if err != nil {
					nonNumeric = v[0]
					log.Printf("GET %s non-numeric game: %q", r.RequestURI, nonNumeric)
					http.Error(w, "game "+nonNumeric+" not found", http.StatusNotFound)
					return
				} else {
					log.Printf("GET %s numeric game: %d", r.RequestURI, gameId)
					if gameId > -1 {
						path, err := pathForGameId(gameId)
						if err != nil {
							log.Printf("ERROR: file not found for game %d", gameId)
							http.Error(w, fmt.Sprintf("game %d file not found", gameId), http.StatusNotFound)
							return
						}
						b, err := readFileReplacingHost(path)
						if err != nil {
							log.Printf("ERROR: opening file %q: %v", path, err)
							http.Error(w, "", http.StatusInternalServerError)
							return
						}
						sb.Write(b)
					}
				}
			}
		}
	}

	if sb.Len() > 0 {
		n, err := w.Write([]byte(sb.String()))
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		log.Printf("%d bytes written in response to %q", n, r.Method+" "+r.URL.String())
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func pathForGameId(gameId int) (string, error) {
	fillGameIdPathMap()
	if path, ok := gameIdPathMap[gameId]; ok {
		return path, nil
	}
	return "", fmt.Errorf("no path for id %d", gameId)
}

var gameIdPathMap map[int]string
var gamePathRegex *regexp.Regexp = regexp.MustCompile(`^.*\/game_([0-9]+)\.php$`)

func fillGameIdPathMap() {
	if len(gameIdPathMap) > 0 {
		return
	}

	gameIdPathMap = map[int]string{}
	err := fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if gamePathRegex.MatchString(path) {
			sm := gamePathRegex.FindStringSubmatch(path)
			if len(sm) > 1 {
				if gid, e := strconv.Atoi(sm[1]); e == nil {
					gameIdPathMap[gid] = path
				} else {
					log.Printf("ERROR: %v", e)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	log.Printf("len(gameIdPathMap) = %d", len(gameIdPathMap))
}

var replaceHostRegex *regexp.Regexp = regexp.MustCompile(`href=\"https:\/\/j-archive.com\/`)

func readFileReplacingHost(filename string) ([]byte, error) {
	b, err := content.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return replaceHostRegex.ReplaceAll(b, []byte(fmt.Sprintf(`href="http://localhost:%s/`, port))), nil
}
