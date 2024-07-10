package letterboxd

import (
	"math/rand/v2"
	"path"
	"path/filepath"
	"runtime"

	"github.com/MathisVerstrepen/go-module/webfetch"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func LoadFetcherFromProxies() *[]webfetch.Fetcher {
	return webfetch.InitFetchers(path.Join(basepath, "../.."))
}

type Fetcher struct {
	*webfetch.Fetcher
}

type Fetchers struct {
	fetchers []Fetcher
}

func (lf *Fetchers) AddFetchers(fetchers *[]webfetch.Fetcher) {
	for i := range *fetchers {
		fetcher := (*fetchers)[i]
		lf.fetchers = append(lf.fetchers, Fetcher{
			Fetcher: &fetcher,
		})
	}
}

func (lf Fetchers) GetFetcher(idx int) Fetcher {
	return lf.fetchers[idx]
}

func (lf Fetchers) GetRandomFetcher() Fetcher {
	idx := rand.IntN(len(lf.fetchers))
	return lf.GetFetcher(idx)
}
