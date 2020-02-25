package main

import (
	"flag"
	"net/url"
	"os"
	"regexp"
	"sync"

	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	writer.Start()
	defer writer.Stop()

	configPtr := flag.String("config", "isogo.yml", "the YAML config file")

	downloadPtr := flag.Bool("download", false, "download ISOs")
	keepPtr := flag.Bool("keep", false, "run keep jobs")

	flag.Parse()

	if !*downloadPtr && !*keepPtr {
		log.WithFields(log.Fields{
			"configFile": *configPtr,
		}).Error("You need to specify at least one action")

		flag.Usage()

		os.Exit(1)
	}

	config, err := readConfig(*configPtr)
	if err != nil {
		log.WithFields(log.Fields{
			"configFile": *configPtr,
			"error":      err.Error(),
		}).Fatal("Error reading config file")
	}

	log.WithFields(log.Fields{
		"configFile":   *configPtr,
		"downloadJobs": len(config.ISOs),
		"keepJobs":     len(config.Keep),
	}).Info("Successfully read configuration")

	if *downloadPtr {
		log.Info("Running download jobs")

		for k := range config.ISOs {
			iso := config.ISOs[k]

			log.WithFields(log.Fields{
				"name":        iso.Name,
				"url":         iso.URL,
				"discover":    iso.Discover,
				"regex":       iso.Regex,
				"destination": iso.Destination,
			}).Info("Starting download job")

			ur, err := url.Parse(iso.URL)

			if err != nil {
				log.WithFields(log.Fields{
					"name":  iso.Name,
					"error": err.Error(),
				}).Error("Error running download job, skipping")
				continue
			}

			if iso.Discover {
				urls, err := discover(ur, regexp.MustCompile(iso.Regex))

				if err != nil {
					log.WithFields(log.Fields{
						"name":  iso.Name,
						"error": err.Error(),
					}).Error("Error running discovery for download job, skipping")
					continue
				}

				for k := range urls {
					u := urls[k]

					addDownload(iso.Name, u.String(), iso.Destination)
				}
			} else {
				addDownload(iso.Name, ur.String(), iso.Destination)
			}
		}

		go displayProgress()

		wg.Wait()
	}

	if *keepPtr {
		log.Info("Running keep jobs")
		for k := range config.Keep {
			keep := config.Keep[k]

			log.WithFields(log.Fields{
				"directory": keep.Directory,
				"regex":     keep.Regex,
				"last":      keep.Last,
			}).Info("Starting keep job")

			err = keepDir(keep.Directory, regexp.MustCompile(keep.Regex), keep.Last)

			if err != nil {
				log.WithFields(log.Fields{
					"directory": keep.Directory,
					"regex":     keep.Regex,
					"last":      keep.Last,
					"error":     err.Error(),
				}).Error("Error running keep job, skipping")
				continue
			}
		}
	}

}
