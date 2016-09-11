package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GeoNet/slink"
	"github.com/ozym/raw"
)

func main() {

	// storage options
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

	// seedlink options
	var netdly int
	flag.IntVar(&netdly, "netdly", 0, "provide network delay")
	var netto int
	flag.IntVar(&netto, "netto", 300, "provide network timeout")
	var keepalive int
	flag.IntVar(&keepalive, "keepalive", 0, "provide keep-alive")
	var selectors string
	//flag.StringVar(&selectors, "selectors", "50L?? 51L??", "provide channel selectors")
	flag.StringVar(&selectors, "selectors", "???", "provide channel selectors")
	var streams string
	flag.StringVar(&streams, "streams", "*_*", "provide streams")
	var statefile string
	flag.StringVar(&statefile, "statefile", "", "provide a running state file")
	var state time.Duration
	flag.DurationVar(&state, "state", 30.0*time.Second, "how often to save state")

	// heartbeat flush interval
	var flush time.Duration
	flag.DurationVar(&flush, "flush", 60.0*time.Second, "how often to update files")

	flag.Parse()

	storage, err := raw.NewTemplate(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// who to call ...
	server := "localhost:18000"
	if flag.NArg() > 0 {
		server = flag.Arg(0)
	}

	// initial seedlink handle
	slconn := slink.NewSLCD()
	defer slink.FreeSLCD(slconn)

	// seedlink settings
	slconn.SetNetDly(netdly)
	slconn.SetNetTo(netto)
	slconn.SetKeepAlive(keepalive)

	// conection
	slconn.SetSLAddr(server)
	defer slconn.Disconnect()

	// configure streams selectors to recover
	slconn.ParseStreamList(streams, selectors)

	if statefile != "" {
		if _, err := os.Stat(statefile); err == nil {
			log.Println("read initial state")
			if x := slconn.RecoverState(statefile); x != 0 {
				log.Println("unable to read state: %s", statefile)
			}
		}
	}

	// handle process signals via channels
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, syscall.SIGINT)
	signal.Notify(halt, syscall.SIGTERM)

	// periodicly save state
	tick := time.NewTicker(state)

	// periodicly flush the buffers
	tock := time.NewTicker(flush)

	var readings []raw.Reading

	log.Printf("collecting: %s (%s) :: %s", streams, selectors, server)

loop:
	for {
		select {
		case <-halt:
			break loop
		case <-tick.C:
			if statefile != "" {
				if x := slconn.SaveState(statefile); x != 0 {
					log.Fatalf("unable to write state: %s", statefile)
				}
			}
		case <-tock.C:
			if len(readings) > 0 {
				log.Printf("flush: %d records", len(readings))
				if err := raw.Store(dir, raw.NewCsv(dp), storage.Execute, readings); err != nil {
					log.Fatalf("unable to store readings: %v", err)
				}
				readings = nil
			}
		default:
			// recover packet ...
			switch p, rc := slconn.CollectNB(); rc {
			case slink.SLTERMINATE:
				log.Printf("terminating")
				break loop
			case slink.SLNOPACKET:
				time.Sleep(100 * time.Millisecond)
				continue loop
			case slink.SLPACKET:
				// check just in case we're shutting down
				if p != nil && p.PacketType() == slink.SLDATA {
					r, err := raw.DecodeMSeedBuffer(p.GetMSRecord(), offset, scale)
					if err != nil {
						log.Fatalf("unable to decode mseed buffer: %v", err)
					}
					readings = append(readings, r...)
				}
			default:
				log.Fatal("invalid packet")
			}
		}
	}

	if len(readings) > 0 {
		log.Printf("flush: %d records", len(readings))
		if err := raw.Store(dir, raw.NewCsv(dp), storage.Execute, readings); err != nil {
			log.Fatalf("unable to store readings: %v", err)
		}
	}

	if statefile != "" {
		log.Println("write final state")
		if x := slconn.SaveState(statefile); x != 0 {
			log.Fatalf("unable to write state: %s", statefile)
		}
	}

	log.Println("terminated")
}
