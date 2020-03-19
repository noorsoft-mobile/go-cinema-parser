package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strings"
)

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://kinoteatr.ru/raspisanie-kinoteatrov/city/#"},
		ParseFunc: parseMovies,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}

func parseMovies(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.shedule_movie").Each(func(i int, s *goquery.Selection) {
		var sessions = strings.Split(s.Find(".shedule_session_time").Text(), " \n ")
		sessions = sessions[:len(sessions)-1]

		for i := 0; i < len(sessions); i++ {
			sessions[i] = strings.Trim(sessions[i], "\n ")
		}

		var description string

		if href, ok := s.Find("a.gtm-ec-list-item-movie").Attr("href"); ok {
			g.Get(r.JoinURL(href), func(_g *geziyor.Geziyor, _r *client.Response) {
				description = _r.HTMLDoc.Find("span.announce p.movie_card_description_inform").Text()

				description = strings.ReplaceAll(description, "\t", "")
				description = strings.ReplaceAll(description, "\n", "")
				description = strings.TrimSpace(description)

				g.Exports <- map[string]interface{}{
					"title":        strings.TrimSpace(s.Find("span.movie_card_header.title").Text()),
					"subtitle":    strings.TrimSpace(s.Find("span.sub_title.shedule_movie_text").Text()),
					"sessions":    sessions,
					"description": description,
				}
			})
		}
	})
}
