package main

import (
	"flag"
	"log"
	"os"

	"github.com/ozym/raw"
)

func main() {

	var dir string
	flag.StringVar(&dir, "dir", ".", "output base directory")

	var tmpl string
	flag.StringVar(&tmpl, "template", "{{Year .Epoch}}/{{Year .Epoch}}.{{Doy .Epoch}}/{{Year .Epoch}}.{{Doy .Epoch}}.{{Hour .Epoch}}.{{.Source}}.csv", "file name template")

	var scale float64
	flag.Float64Var(&scale, "scale", 1.0, "stream scale factor")

	var offset float64
	flag.Float64Var(&offset, "offset", 0.0, "stream offset factor")

	var dp int
	flag.IntVar(&dp, "dp", -1, "decimal places")

	flag.Parse()

	storage, err := raw.NewTemplate(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	var readings []raw.Reading
	for _, infile := range flag.Args() {
		switch infile {
		case "-":
			log.Println("reading: stdin")
			r, err := raw.ReadMSeedStream(os.Stdin, offset, scale)
			if err != nil {
				log.Fatal(err)
			}
			readings = append(readings, r...)
		default:
			log.Printf("reading: %s", infile)
			r, err := raw.ReadMSeedFile(infile, offset, scale)
			if err != nil {
				log.Fatal(err)
			}
			readings = append(readings, r...)
		}
	}

	log.Printf("storing %d readings: %s", len(readings), dir)
	if err := raw.Store(dir, raw.NewCsv(dp), storage.Execute, readings); err != nil {
		log.Fatal(err)
	}
}
