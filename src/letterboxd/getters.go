package letterboxd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	N      int
	Offset int
}

func (lg LetterboxdGetter) GetPopularMovies(n int, offset int, dateRange Range) (*PopularMovies, error) {
	// TODO implement offset
	filename := "films/ajax/popular"
	page := 1
	nFetched := 0

	staticFetcher := lg.Fetchers.GetRandomFetcher()

	popMovies := PopularMovies{
		N:      n,
		Offset: offset,
	}
	popMovies.Movies = make([]MovieMeta, n)

	for {
		resBytes, err := staticFetcher.FetchData(webfetch.FetcherParams{
			Method:   "GET",
			Url:      fmt.Sprintf("%s://%s/%s/%s/page/%d", Scheme, Host, filename, dateRange, page),
			UseProxy: true,
			Params: webfetch.Param{
				"esiAllowFilters": "true",
			},
		})

		if err != nil {
			return nil, err
		}

		htmlBody, err := resBytes.ToHTMLNode()
		if err != nil {
			return nil, err
		}

		posterNodes := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
			ClassNames: "poster film-poster",
			Tag:        "div",
			Multiple:   true,
		})

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

			nFetched += 1
			if nFetched >= n {
				break
			}
		}

		if nFetched >= n {
			break
		}

		page += 1
	}

	return &popMovies, nil
}

type MovieStat struct {
	WatchCount int
	ListCount  int
	LikeCount  int
}

func extractNumber(text string) (int, error) {
	r, err := regexp.Compile(`[\d,]+`)
	if err != nil {
		return -1, fmt.Errorf("erreur de compilation de la regex: %w", err)
	}

	matches := r.FindAllString(text, -1)
	if len(matches) > 0 {
		numberString := matches[0]
		numberString = strings.ReplaceAll(numberString, ",", "")
		number, err := strconv.Atoi(numberString)
		return number, err
	}

	return -1, fmt.Errorf("aucun nombre trouvé dans: %s", text)
}

func (lg LetterboxdGetter) GetMovieStats(movieUrl string) (*MovieStat, error) {
	staticFetcher := lg.Fetchers.GetRandomFetcher()

	resBytes, err := staticFetcher.FetchData(webfetch.FetcherParams{
		Method:   "GET",
		Url:      fmt.Sprintf("%s://%s/csi%sstats", Scheme, Host, movieUrl),
		UseProxy: true,
	})
	if err != nil {
		return nil, err
	}

	htmlBody, err := resBytes.ToHTMLNode()
	if err != nil {
		return nil, err
	}

	watchNode := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
		ClassNames: "icon-watched",
		Tag:        "a",
		Multiple:   false,
	})
	watchTitle := gosoup.GetAttribute(watchNode[0], "title")
	watchStat, err := extractNumber(watchTitle)
	if err != nil {
		return nil, err
	}

	listNode := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
		ClassNames: "icon-list",
		Tag:        "a",
		Multiple:   false,
	})
	listTitle := gosoup.GetAttribute(listNode[0], "title")
	listStat, err := extractNumber(listTitle)
	if err != nil {
		return nil, err
	}

	likeNode := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
		ClassNames: "icon-liked",
		Tag:        "a",
		Multiple:   false,
	})
	likeTitle := gosoup.GetAttribute(likeNode[0], "title")
	likeStat, err := extractNumber(likeTitle)
	if err != nil {
		return nil, err
	}

	return &MovieStat{
		WatchCount: watchStat,
		ListCount:  listStat,
		LikeCount:  likeStat,
	}, nil
}
