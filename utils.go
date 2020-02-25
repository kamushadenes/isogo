package main

import "net/url"

func uniqueURL(urlSlice []*url.URL) []*url.URL {
	keys := make(map[string]bool)
	list := make([]*url.URL, 0)

	for _, entry := range urlSlice {
		if _, value := keys[entry.String()]; !value {
			keys[entry.String()] = true
			list = append(list, entry)
		}
	}

	return list
}
