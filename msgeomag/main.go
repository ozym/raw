package main

import (
	"flag"
	"log"
	"os"

	"github.com/ozym/geomag"
)

func main() {

	var dir string
	flag.StringVar(&dir, "dir", ".", "output base directory")

	var tmpl string
	flag.StringVar(&tmpl, "template", "{{Year .Epoch}}/{{Year .Epoch}}.{{Doy .Epoch}}/{{Year .Epoch}}.{{Doy .Epoch}}.{{Hour .Epoch}}.{{.Source}}.csv", "file name template")

	var scale float64
	flag.Float64Var(&scale, "scale", 1.0, "stream scale factor")

	var dp int
	flag.IntVar(&dp, "dp", -1, "decimal places")

	flag.Parse()

	mseed, err := geomag.NewMSeed(scale)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := geomag.NewStorage(tmpl, dp)
	if err != nil {
		log.Fatal(err)
	}

	var readings []geomag.Reading
	for _, infile := range flag.Args() {
		switch infile {
		case "-":
			log.Println("reading: stdin")
			r, err := mseed.ReadStream(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			readings = append(readings, r...)
		default:
			log.Printf("reading: %s", infile)
			r, err := mseed.ReadFile(infile)
			if err != nil {
				log.Fatal(err)
			}
			readings = append(readings, r...)
		}
	}

	log.Printf("storing %d readings: %s", len(readings), dir)
	if err := storage.Store(dir, readings); err != nil {
		log.Fatal(err)
	}
}
