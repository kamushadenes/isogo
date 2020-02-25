package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func discover(durl *url.URL, regex *regexp.Regexp) ([]*url.URL, error) {
	var urls = make([]*url.URL, 0)

	resp, err := http.Get(durl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if regex.MatchString(href) {
				if strings.HasPrefix(href, "/") {
					nurl, _ := url.Parse(durl.String())
					nurl.Path = href

					urls = append(urls, nurl)

				} else if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
					nurl, _ := url.Parse(href)

					if err == nil {
						urls = append(urls, nurl)
					}

				} else {
					nurl, _ := url.Parse(durl.String())
					nurl.Path = path.Join(nurl.Path, href)

					urls = append(urls, nurl)
				}
			}
		}
	})

	urls = uniqueURL(urls)

	return urls, nil
}
