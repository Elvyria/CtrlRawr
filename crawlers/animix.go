package crawlers

import (
	"regexp"
	"errors"
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Animix struct {}

func init() {
	Crawlers["animixplay.to"] = new(Animix)
}

type animix_episodes struct {
	Total    uint                   `json:"eptotal"`
	Episodes map[string]interface{} `json:"-"`
}

func (crawler *Animix) Crawl(s *Spider) (info *Info, err error) {
	c := colly.NewCollector()

	c.OnHTML("#lowerplayerpage", func(e *colly.HTMLElement) {
		info, err = crawler.parseStream(e)
	})

	c.OnError(func(r *colly.Response, e error) {
		err = e
	})

	c.Visit(s.Url)

	return info, err
}

func (_ *Animix) ClearUrl(ref string) string {
	reStream := regexp.MustCompile(`(https:\/\/animixplay\.to\/v\d+\/[^\/]+)`)

	matches := reStream.FindSubmatch([]byte(ref))
	if len(matches) > 1 {
		return string(matches[1])
	}

	return ref
}

func (crawler *Animix) parseStream(e *colly.HTMLElement) (info *Info, err error) {
	info = &Info{}

	info.Title = e.ChildText(".animetitle")

	epslist := e.ChildText("#epslistplace")

	episodes := animix_episodes{}
	if err := json.Unmarshal([]byte(epslist), &episodes.Episodes); err != nil {
		return nil, err
	}

	episodes.Total = uint(episodes.Episodes["eptotal"].(float64))
	info.Message = crawler.message(episodes.Total)

	if (episodes.Total == 0) {
		return nil, errors.New("0 episodes")
	}

	info.Url = "https:" + episodes.Episodes[fmt.Sprint(episodes.Total - 1)].(string)

	return info, err
}

func (_ *Animix) message(episode uint) string {
	return fmt.Sprintf("Episode %d just released!", episode)
}
