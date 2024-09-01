package main

import (
	"fmt"

	"diikstra.fr/letterboxd-statistics/src/letterboxd"
)

func main() {
	rawFetchers := letterboxd.LoadFetcherFromProxies()

	letterboxdFetchers := letterboxd.Fetchers{}
	letterboxdFetchers.AddFetchers(rawFetchers)

	letterboxdGetter := letterboxd.LetterboxdGetter{Fetchers: letterboxdFetchers}
	popularMovies, _ := letterboxdGetter.GetPopularMovies(10, 0, letterboxd.Week)

	for _, movie := range popularMovies.Movies {
		movieStat, _ := letterboxdGetter.GetMovieStats(movie.Link)
		fmt.Printf("%s : %+v\n", movie.Slug, *movieStat)
	}
}
