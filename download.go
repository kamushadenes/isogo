package main

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cavaliercoder/grab"
)

var downloadMutex sync.Mutex

type Download struct {
	Name           string
	URL            *url.URL
	Destination    string
	TotalBytes     int64
	CurrentBytes   int64
	Progress       float64
	Hash           string
	Sum            string
	StartTime      *time.Time
	EndTime        *time.Time
	Status         string
	FailReason     string
	BytesPerSecond float64
	ETA            time.Time
	Filename       string
}

var client = grab.NewClient()
var downloading = make(map[string]*Download)

func newDownload(name string, url *url.URL, dst string) *Download {
	t := time.Now()
	return &Download{Name: name, StartTime: &t, URL: url, Destination: dst, Status: "Created"}
}

func addDownload(name string, url string, dst string) error {
	dst, _ = filepath.Abs(dst)

	fi, err := os.Stat(dst)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}
	} else if !fi.Mode().IsDir() {
		err = fmt.Errorf("%s is not an directory", dst)
		return err
	}

	req, _ := grab.NewRequest(dst, url)

	hasHash, hashName, hash, sum := getHash(req.URL(), filepath.Base(req.URL().Path))

	d := newDownload(name, req.URL(), dst)

	if hasHash {
		d.Hash = hashName
		d.Sum = hex.EncodeToString(sum)

		req.SetChecksum(hash, sum, true)
	}

	downloadMutex.Lock()
	downloading[name] = d
	downloadMutex.Unlock()

	wg.Add(1)
	go d.Do(req)

	return nil
}

func (d *Download) Do(req *grab.Request) {

	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			downloadMutex.Lock()
			d.Status = "Downloading"
			d.CurrentBytes = resp.BytesComplete()
			d.TotalBytes = resp.Size
			d.Progress = resp.Progress()
			d.BytesPerSecond = resp.BytesPerSecond()
			d.ETA = resp.ETA()
			d.Filename = resp.Filename
			downloadMutex.Unlock()
		case <-resp.Done:
			downloadMutex.Lock()
			// download is complete
			endt := time.Now()
			d.EndTime = &endt
			downloadMutex.Unlock()
			break Loop
		}
	}

	// check for errors
	downloadMutex.Lock()
	if err := resp.Err(); err != nil {
		d.Status = "Failed"
	} else {
		d.Status = "Completed"
	}
	downloadMutex.Unlock()

	time.Sleep(1 * time.Second)

	wg.Done()
}
