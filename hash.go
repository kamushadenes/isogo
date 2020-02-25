package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

var possibleHashFiles = []string{
	"SHA256SUMS", "sha256sums.txt",
	"SHA1SUMS", "sha1sums.txt",
	"MD5SUMS", "md5sums.txt"}

func getHash(url *url.URL, fname string) (hasHash bool, hashName string, h hash.Hash, sum []byte) {
	var hsum string
	var err error

	fname = filepath.Base(fname)

	for _, v := range possibleHashFiles {
		nurl, _ := url.Parse(url.String())

		pathParts := strings.Split(nurl.Path, "/")
		pathParts = pathParts[:len(pathParts)-1]
		pathParts = append(pathParts, v)

		nurl.Path = strings.Join(pathParts, "/")

		resp, err := http.Get(nurl.String())
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(bytes.NewReader(body))

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, fname) {
				if strings.Contains(strings.ToLower(nurl.Path), "sha256sums") {
					h = sha256.New()
					hashName = "SHA256"
				} else if strings.Contains(strings.ToLower(nurl.Path), "sha1sums") {
					h = sha1.New()
					hashName = "SHA1"
				} else if strings.Contains(strings.ToLower(nurl.Path), "md5sums") {
					h = md5.New()
					hashName = "MD5"
				}
				hsum = strings.Fields(line)[0]
				break
			}
		}
		if err := scanner.Err(); err != nil {
			continue
		} else {
			break
		}

	}

	if hsum != "" {
		sum, err = hex.DecodeString(hsum)
		if err == nil {
			hasHash = true
		}
	}

	return
}
