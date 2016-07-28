package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
)

var myname = "jhevolver"
var myversion = "0.1"

var usage = fmt.Sprintf(`%s: 

Usage:
  %s evolve [options]
  %s --help

Options:
  -i, --in=<filename>	  Input file to evolve [default: true].
  -o, --out=<dirname>	  Output directory [default: evolved].
  -l, --limit=<num>	      Limit evolution to first <num> bytes.
  -v, --verbose		      Verbose output.
  -h, --help              Show this screen.
  --version               Show version.
`, myname, myname, myname)

func main() {
	args, err := docopt.Parse(usage, nil, true, myversion, false)
	if err != nil {
		log.Fatal(err)
	}
	verbose := args["--verbose"].(bool)

	if args["evolve"].(bool) {
		binary := args["--in"].(string)
		if _, err := os.Stat(binary); err != nil {
			log.Fatalf("Binary file %s not accessible: %v", binary, err)
		}
		outdir := args["--out"].(string)
		if _, err := os.Stat(outdir); !os.IsNotExist(err) {
			log.Fatalf("Output directory %s already exists!", outdir)
		}
		if err := os.Mkdir(outdir, 0755); err != nil {
			log.Fatalf("Error creating output directory %s: %v", outdir, err)
		}
		limitStr, ok := args["--limit"].(string)
		if !ok {
			limitStr = "0"
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Fatal(err)
		}
		if err := evolve(binary, outdir, limit, verbose); err != nil {
			log.Fatal(err)
		}
	}

}

func evolve(binary, outdir string, limit int, verbose bool) error {
	log.Printf("Started evolution on binary %s.", binary)

	// read entire binary into memory - this won't scale to larger binaries...
	raw, err := ioutil.ReadFile(binary)
	if err != nil {
		return fmt.Errorf("Error reading binary file %s: %v", binary, err)
	}
	log.Printf("Binary contains %d bytes, %d bits", len(raw), len(raw)*8)
	for byteNum, B := range raw {
		log.Printf("Evolving byte #%d/%d: %.8b", byteNum, len(raw), raw[byteNum])
		for bitNum := 7; bitNum >= 0; bitNum-- {
			globalBit := (byteNum * 8) + (8 - bitNum)
			curVal := getBit(int(B), uint(bitNum))
			newVal := 0
			if curVal == 0 {
				newVal = 1
			}
			filename := fmt.Sprintf("%s/evolved-%d", outdir, globalBit)
			err := ioutil.WriteFile(filename, raw, 0644)
			if err != nil {
				return fmt.Errorf("Error writing evolved file %s: %v", filename, err)
			}
			if verbose {
				log.Printf("Evolved bit %d: %b->%b [filename: %s]", globalBit, curVal, newVal, filename)
			}
		}
		if limit > 0 && byteNum >= limit {
			log.Printf("Stoping evolution at byte %d as requested.", limit)
			break
		}
	}
	log.Printf("Finished %s with binary %s.", myname, binary)
	return nil

}

func getBit(n int, pos uint) int {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}
