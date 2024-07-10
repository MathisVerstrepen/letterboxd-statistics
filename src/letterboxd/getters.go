package letterboxd

import (
	"fmt"
	"log"

	"github.com/MathisVerstrepen/go-module/gosoup"
	"github.com/MathisVerstrepen/go-module/webfetch"
)

const (
	Scheme string = "https"
	Host   string = "letterboxd.com"
)

type LetterboxdGetter struct {
	Fetchers Fetchers
}

type Range string

const (
	Week  Range = "this/week"
	Month Range = "this/month"
	Year  Range = "this/year"
	All   Range = ""
)

type MovieMeta struct {
	Id   string
	Slug string
	Link string
}

type PopularMovies struct {
	Movies []MovieMeta
	n      int
	offset int
}

func (lg LetterboxdGetter) GetPopularMovies(n int, offset int, dateRange Range) PopularMovies {
	// TODO implement n and offset
	filename := "films/ajax/popular"
	staticFetcher := lg.Fetchers.GetRandomFetcher()

	popMovies := PopularMovies{
		n:      n,
		offset: offset,
	}

	resBytes, err := staticFetcher.FetchData(webfetch.FetcherParams{
		Method:   "GET",
		Url:      fmt.Sprintf("%s://%s/%s/%s", Scheme, Host, filename, dateRange),
		UseProxy: true,
		Params: webfetch.Param{
			"esiAllowFilters": "true",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	htmlBody, err := resBytes.ToHTMLNode()
	if err != nil {
		log.Fatal(err)
	}

	posterNodes := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
		ClassNames: "poster film-poster",
		Tag:        "div",
		Multiple:   true,
	})

	popMovies.Movies = make([]MovieMeta, len(posterNodes))
	for i, posterNode := range posterNodes {
		movieId := gosoup.GetAttribute(posterNode, "data-film-id")
		movieSlug := gosoup.GetAttribute(posterNode, "data-film-slug")
		movieLink := gosoup.GetAttribute(posterNode, "data-target-link")

		fmt.Printf("Found id [%10s], name [%50s], link [%50s] \n", movieId, movieSlug, movieLink)

		popMovies.Movies[i] = MovieMeta{
			Id:   movieId,
			Slug: movieSlug,
			Link: movieLink,
		}
	}

	return popMovies
}
