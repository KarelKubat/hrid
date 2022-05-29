package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/KarelKubat/flagnames"
	"github.com/KarelKubat/hrid/id"
)

const (
	usage = `
This is hrid, the Human Readable ID converter.
Usage:
  hrid [FLAGS] NUMBER - generates a human readable ID and prints it on stdout
  hrid [FLAGS] -id ID - re-interprets the ID as a number and prints it on stdout

The flags can be abbreviated: -a for -alphabet, -l for -length etc.
Supported flags:
`
)

var (
	alphabetFlag   = flag.String("alphabet", id.Alphabet, "conversion alphabet: first rune represents 0, second 1, etc.")
	lenFlag        = flag.Int("length", id.StringLen, "minimum length of generated IDs, set to 0 for no padding")
	ignoreCaseFlag = flag.Bool("ignorecase", id.IgnoreCase, "when true, casing is ignored when converting IDs to numbers")
	groupsizeFlag  = flag.Int("groupsize", id.GroupSize, "size of space-delimited groups in generated IDs, for better readability")
	checksumFlag   = flag.Int("checksum", id.ChecksumLen, "number of checksum runes to append")

	idFlag      = flag.Bool("id", false, "when true, arguments are taken as IDs, default: numbers")
	verboseFlag = flag.Bool("verbose", false, "show options with which the converter is instantiated")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flagnames.Patch()
	flag.Parse()
	hrid(flag.Args())
}

// hrid is a helper function that can be called from the unit test.
func hrid(args []string) {
	if len(args) == 0 {
		flag.Usage()
	}
	opts := &id.Opts{
		Alphabet:    *alphabetFlag,
		StringLen:   *lenFlag,
		IgnoreCase:  *ignoreCaseFlag,
		GroupSize:   *groupsizeFlag,
		ChecksumLen: *checksumFlag,
	}
	idConverter, err := id.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	if *verboseFlag {
		log.Printf("Converter options: %+v", *opts)
	}
	for _, a := range args {
		if *idFlag {
			n, err := idConverter.ToNr(a)
			if err != nil {
				log.Printf("%v: not a valid ID: %v", a, err)
			} else {
				fmt.Println(n)
			}
		} else {
			u, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				log.Printf("%v: not a valid number: %v", a, err)
			} else {
				fmt.Println(idConverter.ToString(u))
			}
		}
	}
}
