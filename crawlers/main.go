package crawlers

import (
	"time"
	"net/url"
	"errors"
)

type Crawler interface {
	Crawl(s *Spider) (*Info, error)
}

type Spider struct {
	Hash string `json:"hash"`
	Host string `json:"host"`
	Url  string `json:"url"`
	LastUpdated time.Time `json:"lastUpdate"`
}

type Info struct {
	Title   string
	Message string
	Picture string
	Url     string
}

func Create(ref string) (*Spider, error) {
	u, err := url.Parse(ref)
	if (err != nil) {
		return nil, err
	}

	_, exists := Crawlers[u.Host]
	if (!exists) {
		return nil, errors.New("no scrapper for host")
	}

	s := Spider {
		Hash: "",
		Host: u.Host,
		Url: ref,
		LastUpdated: time.Unix(0, 0),
	}

	return &s, nil
}

func (s *Spider) Crawl() (info *Info, err error) {
	crawler, exists := Crawlers[s.Host]
	if (!exists) {
		return nil, err
	}

	info, err = crawler.Crawl(s)
	if (err != nil) {
		return nil, err
	}

	s.LastUpdated = time.Now()
	s.Hash = info.Message

	return info, err
}

var Crawlers = map[string]Crawler{}
