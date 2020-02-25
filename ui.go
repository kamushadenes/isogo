package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/gosuri/uilive"
)

var writer = uilive.New()

func displayProgress() {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			downloadMutex.Lock()
			wg.Add(1)
			msg := ""

			keys := make([]string, 0)

			for k := range downloading {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				d := downloading[k]

				if d.Status != "Completed" && d.Status != "Failed" {
					msg = fmt.Sprintf(`%s
%s - %s
    Start: %s
    ETA: %s
    Source: %s://%s/.../%s
    Destination: %s/%s
    Expected Hash: %s %s
    Speed %.2f kb/s
    Progress: %v / %v (%.2f%%)

`,
						msg,
						d.Name, d.Status,
						d.StartTime.Format(time.RFC3339),
						d.ETA.Format(time.RFC3339),
						d.URL.Scheme, d.URL.Hostname(), filepath.Base(d.URL.Path),
						d.Destination, filepath.Base(d.URL.Path),
						d.Hash, d.Sum,
						d.BytesPerSecond/1024,
						d.CurrentBytes, d.TotalBytes, 100*d.Progress)
				} else {
					msg = fmt.Sprintf(`%s
%s - %s
    Start: %s
    End: %s
    Source: %s://%s/.../%s
    Destination: %s/%s
    Expected Hash: %s %s

`,
						msg,
						d.Name, d.Status,
						d.StartTime.Format(time.RFC3339),
						d.EndTime.Format(time.RFC3339),
						d.URL.Scheme, d.URL.Hostname(), filepath.Base(d.URL.Path),
						d.Destination, filepath.Base(d.URL.Path),
						d.Hash, d.Sum)

				}
			}
			fmt.Fprint(writer, msg)
			wg.Done()
			downloadMutex.Unlock()
		}
	}
}
