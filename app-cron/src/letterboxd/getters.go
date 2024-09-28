package letterboxd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/MathisVerstrepen/go-module/gosoup"
	"github.com/MathisVerstrepen/go-module/webfetch"
	"golang.org/x/net/html"
)

const (
	Scheme string = "https"
	Host   string = "letterboxd.com"
)

type LetterboxdGetter struct {
	Fetchers Fetchers
}

type MovieMeta struct {
	Id       string
	Slug     string
	Link     string
	Title    string
	Rating   float32
	Poster   string
	Backdrop string
}

type PopularMovies struct {
	Movies []MovieMeta
	N      int
	Offset int
}

type MovieIds []string

func (moviesIds *MovieIds) Include(search string) bool {
	for _, moviesId := range *moviesIds {
		if moviesId == search {
			return true
		}
	}
	return false
}

func getPopularNode(staticFetcher Fetcher, dateRange string, page int) (*html.Node, error) {
	urlBasePopular := "films/ajax/popular"

	resBytes, err := staticFetcher.FetchData(webfetch.FetcherParams{
		Method:   "GET",
		Url:      fmt.Sprintf("%s://%s/%s/%s/page/%d", Scheme, Host, urlBasePopular, dateRange, page),
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

	return htmlBody, nil
}

func (lg LetterboxdGetter) GetPopularMovies(n int, offset int, dateRange string) (*PopularMovies, error) {
	// TODO implement offset

	page := 1
	nFetched := 0

	staticFetcher := lg.Fetchers.GetRandomFetcher()

	popMovies := PopularMovies{
		N:      n,
		Offset: offset,
	}
	popMovies.Movies = make([]MovieMeta, n)

	for {
		htmlBody, err := getPopularNode(staticFetcher, dateRange, page)
		if err != nil {
			return nil, err
		}

		filmListNodes := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
			ClassNames: "listitem poster-container",
			Tag:        "li",
			Multiple:   true,
		})

		for _, filmListNode := range filmListNodes {
			movieRating := gosoup.GetAttribute(filmListNode, "data-average-rating")
			movieRatingFloat, _ := strconv.ParseFloat(movieRating, 32)

			filmNode := gosoup.GetNodeBySelector(filmListNode, &gosoup.HtmlSelector{
				ClassNames: "poster film-poster",
				Tag:        "div",
				Multiple:   false,
			})

			movieId := gosoup.GetAttribute(filmNode[0], "data-film-id")
			movieSlug := gosoup.GetAttribute(filmNode[0], "data-film-slug")
			movieLink := gosoup.GetAttribute(filmNode[0], "data-target-link")

			imgNode := gosoup.GetNodeBySelector(filmNode[0], &gosoup.HtmlSelector{
				Tag:      "img",
				OnlyTag:  true,
				Multiple: false,
			})
			movieTitle := gosoup.GetAttribute(imgNode[0], "alt")

			popMovies.Movies[nFetched] = MovieMeta{
				Id:     movieId,
				Slug:   movieSlug,
				Link:   movieLink,
				Title:  movieTitle,
				Rating: float32(movieRatingFloat),
			}

			fmt.Printf("Found id [%10s], slug [%50s], link [%50s], name [%50s], rating [%2.2f] \n", movieId, movieSlug, movieLink, movieTitle, movieRatingFloat)

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
	LikeCount  int
	Rating     float32
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

	fmt.Printf("Found stats WatchCount [%15d], LikeCount [%15d] for movie [%50s]\n", watchStat, likeStat, movieUrl)

	return &MovieStat{
		WatchCount: watchStat,
		LikeCount:  likeStat,
	}, nil
}

func (lg LetterboxdGetter) GetMovieStatsThreaded(movies *PopularMovies) (map[string]*MovieStat, error) {
	var moviesStat map[string]*MovieStat = make(map[string]*MovieStat, len(movies.Movies))

	sem := make(chan struct{}, 10) // au max, 10 coroutines en même temps
	var mu sync.Mutex              // pour protéger l'accès à la map moviesStat
	var wg sync.WaitGroup          // pour attendre toutes les coroutines même après sortie de boucle

	for _, movie := range movies.Movies {
		wg.Add(1)
		sem <- struct{}{}

		go func(movieId string, movieUrl string, movieRating float32) {
			defer wg.Done()

			movieStat, err := lg.GetMovieStats(movieUrl)
			if err != nil {
				log.Default().Println(err)
			}
			movieStat.Rating = movieRating

			mu.Lock()
			moviesStat[movieId] = movieStat
			mu.Unlock()

			<-sem
		}(movie.Id, movie.Link, movie.Rating)
	}

	wg.Wait()

	return moviesStat, nil
}

func (lg LetterboxdGetter) GetMoviePoster(movieSlug string) (string, error) {
	staticFetcher := lg.Fetchers.GetRandomFetcher()

	urlBasePoster := "ajax/poster"

	resBytes, err := staticFetcher.FetchData(webfetch.FetcherParams{
		Method:   "GET",
		Url:      fmt.Sprintf("%s://%s/%s%sstd/1000x1500", Scheme, Host, urlBasePoster, movieSlug),
		UseProxy: true,
	})

	if err != nil {
		return "", err
	}

	htmlBody, err := resBytes.ToHTMLNode()
	if err != nil {
		return "", err
	}

	posterImgNode := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
		ClassNames: "image",
		Tag:        "img",
		Multiple:   false,
	})
	posterUrl := gosoup.GetAttribute(posterImgNode[0], "src")

	fmt.Printf("Found poster [%150s] for movie [%50s]\n", posterUrl, movieSlug)

	return posterUrl, nil
}

func (lg LetterboxdGetter) SetMoviePosterThreaded(movies *PopularMovies, dbMovieIds *MovieIds) error {
	sem := make(chan struct{}, 10) // au max, 10 coroutines en même temps
	var mu sync.Mutex              // pour protéger l'accès à movies
	var wg sync.WaitGroup          // pour attendre toutes les coroutines même après sortie de boucle

	for idx, movie := range movies.Movies {
		if dbMovieIds.Include(movie.Id) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(movieidx int, movieUrl string) {
			defer wg.Done()

			posterUrl, err := lg.GetMoviePoster(movieUrl)
			if err != nil {
				log.Default().Println(err)
			}

			mu.Lock()
			movies.Movies[movieidx].Poster = posterUrl
			mu.Unlock()

			<-sem
		}(idx, movie.Link)
	}

	wg.Wait()

	return nil
}

func (lg LetterboxdGetter) GetMovieBackdrop(movieUrl string) (string, error) {
	staticFetcher := lg.Fetchers.GetRandomFetcher()

	resBytes, err := staticFetcher.FetchData(webfetch.FetcherParams{
		Method:   "GET",
		Url:      fmt.Sprintf("%s://%s%s", Scheme, Host, movieUrl),
		UseProxy: true,
	})

	if err != nil {
		return "", err
	}

	htmlBody, err := resBytes.ToHTMLNode()
	if err != nil {
		return "", err
	}

	backdropNode := gosoup.GetNodeBySelector(htmlBody, &gosoup.HtmlSelector{
		Id:       "backdrop",
		Multiple: false,
		Tag:      "div",
	})
	if len(backdropNode) == 0 {
		// when movie as no backdrop on letterboxd
		return "", nil
	}
	backdropUrl := gosoup.GetAttribute(backdropNode[0], "data-backdrop2x")

	fmt.Printf("Found backdrop [%150s] for movie [%50s]\n", backdropUrl, movieUrl)

	return backdropUrl, nil
}

func (lg LetterboxdGetter) SetMovieBackdropThreaded(movies *PopularMovies, dbMovieIds *MovieIds) error {
	sem := make(chan struct{}, 10) // au max, 10 coroutines en même temps
	var mu sync.Mutex              // pour protéger l'accès à movies
	var wg sync.WaitGroup          // pour attendre toutes les coroutines même après sortie de boucle

	for idx, movie := range movies.Movies {
		if dbMovieIds.Include(movie.Id) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(movieidx int, movieUrl string) {
			defer wg.Done()

			backdropUrl, err := lg.GetMovieBackdrop(movieUrl)
			if err != nil {
				log.Default().Println(err)
			}

			mu.Lock()
			movies.Movies[movieidx].Backdrop = backdropUrl
			mu.Unlock()

			<-sem
		}(idx, movie.Link)
	}

	wg.Wait()

	return nil
}
