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
	popularMovies := letterboxdGetter.GetPopularMovies(10, 0, letterboxd.Week)

	fmt.Println(popularMovies)
}
